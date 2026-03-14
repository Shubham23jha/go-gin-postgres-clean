package service

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/Shubham23jha/go-gin-postgres-clean/internal/models"
	"github.com/Shubham23jha/go-gin-postgres-clean/internal/repository"
	amqp "github.com/rabbitmq/amqp091-go"
)

type OutboxPublisher interface {
	Start(ctx context.Context)
}

type outboxPublisher struct {
	repo    repository.CampaignRepository
	amqpURL string
	conn    *amqp.Connection
	ch      *amqp.Channel
	mu      sync.Mutex
}

func NewOutboxPublisher(repo repository.CampaignRepository, amqpURL string) OutboxPublisher {
	return &outboxPublisher{
		repo:    repo,
		amqpURL: amqpURL,
	}
}

func (p *outboxPublisher) connect() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.conn != nil && !p.conn.IsClosed() {
		return nil
	}

	conn, err := amqp.Dial(p.amqpURL)
	if err != nil {
		return err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return err
	}

	// Declare Queue with DLX configuration
	args := amqp.Table{
		"x-dead-letter-exchange":    "email_dlx",
		"x-dead-letter-routing-key": "email_failed",
	}
	_, err = ch.QueueDeclare(
		"email_queue", // name
		true,          // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		args,          // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return err
	}

	p.conn = conn
	p.ch = ch
	return nil
}

func (p *outboxPublisher) Start(ctx context.Context) {
	if ctx == nil {
		ctx = context.Background()
	}

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	reclaimTicker := time.NewTicker(1 * time.Minute)
	defer reclaimTicker.Stop()

	log.Println("📤 Outbox Publisher started...")

	for {
		select {
		case <-ctx.Done():
			if p.conn != nil {
				p.conn.Close()
			}
			return
		case <-reclaimTicker.C:
			count, err := p.repo.ReclaimStalledOutbox(5)
			if err == nil && count > 0 {
				log.Printf("♻️ Reclaimed %d stalled outbox items", count)
			}
		case <-ticker.C:
			if err := p.connect(); err != nil {
				log.Printf("❌ Failed to connect to RabbitMQ: %s", err)
				continue
			}
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

		// 3. Publish to RabbitMQ with a simple retry
		var pubErr error
		for attempt := 1; attempt <= 3; attempt++ {
			pubErr = p.ch.PublishWithContext(context.Background(),
				"",            // exchange
				"email_queue", // routing key
				false,         // mandatory
				false,         // immediate
				amqp.Publishing{
					ContentType: "application/json",
					Body:        []byte(item.Payload),
					MessageId:   item.MessageID.String(),
				})

			if pubErr == nil {
				break
			}

			log.Printf("⚠️ Failed to publish item %d (attempt %d/3): %s", item.ID, attempt, pubErr)
			if attempt < 3 {
				time.Sleep(time.Duration(attempt) * time.Second)
				// Reconnect if needed
				if err := p.connect(); err != nil {
					log.Printf("❌ Reconnection failed during publish retry: %s", err)
				}
			}
		}

		if pubErr != nil {
			log.Printf("❌ Final failure publishing item %d to RabbitMQ: %s", item.ID, pubErr)
			// Apply application-level retry/fail logic
			if item.RetryCount < 3 {
				// We don't mark as FAILED here, just log and let it stay PICKED_UP
				// It will be reclaimed by ReclaimStalledOutbox later if it stays PICKED_UP
				// Or we could explicitly move it back to PENDING if it's not a connection error
				p.repo.UpdateOutboxStatus(item.ID, models.OutboxStatusPending)
				// Note: Ideally we increment RetryCount here too.
			} else {
				p.repo.UpdateOutboxFailure(item.ID, pubErr.Error(), item.RetryCount+1)
			}
			continue
		}

		// 4. Mark as PUBLISHED
		p.repo.UpdateOutboxStatus(item.ID, models.OutboxStatusPublished)
	}
}
