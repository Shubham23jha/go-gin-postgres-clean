package service

import (
	"context"
	"log"
	"time"

	"github.com/Shubham23jha/go-gin-postgres-clean/internal/models"
	"github.com/Shubham23jha/go-gin-postgres-clean/internal/repository"
	amqp "github.com/rabbitmq/amqp091-go"
)

type OutboxPublisher interface {
	Start(ctx context.Context)
}

type outboxPublisher struct {
	repo repository.CampaignRepository
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewOutboxPublisher(repo repository.CampaignRepository, amqpURL string) OutboxPublisher {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		log.Fatalf("failed to connect to RabbitMQ: %s", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("failed to open a channel: %s", err)
	}

	// Declare Queue
	_, err = ch.QueueDeclare(
		"email_queue", // name
		true,          // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		log.Fatalf("failed to declare a queue: %s", err)
	}

	return &outboxPublisher{
		repo: repo,
		conn: conn,
		ch:   ch,
	}
}

func (p *outboxPublisher) Start(ctx context.Context) {
	if ctx == nil {
		ctx = context.Background()
	}
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	log.Println("📤 Outbox Publisher started...")

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			p.processOutbox()
		}
	}
}

func (p *outboxPublisher) processOutbox() {
	// 1. Fetch pending items
	items, err := p.repo.FetchPendingOutbox(10) // fetch 10 at a time
	if err != nil {
		log.Printf("failed to fetch pending outbox: %s", err)
		return
	}

	for _, item := range items {
		log.Printf("📤 Publishing Outbox item: %d (Recipient: %s)", item.ID, item.Recipient)

		// 2. Mark as PICKED_UP to avoid double processing
		p.repo.UpdateOutboxStatus(item.ID, models.OutboxStatusPickedUp)

		// 3. Publish to RabbitMQ
		err := p.ch.PublishWithContext(context.Background(),
			"",            // exchange
			"email_queue", // routing key
			false,         // mandatory
			false,         // immediate
			amqp.Publishing{
				ContentType: "application/json",
				Body:        []byte(item.Payload),
				MessageId:   item.MessageID.String(),
			})

		if err != nil {
			log.Printf("failed to publish to RabbitMQ: %s", err)
			p.repo.UpdateOutboxStatus(item.ID, models.OutboxStatusFailed)
			continue
		}

		// 4. Mark as PUBLISHED
		p.repo.UpdateOutboxStatus(item.ID, models.OutboxStatusPublished)
	}
}
