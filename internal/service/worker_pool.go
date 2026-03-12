package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/Shubham23jha/go-gin-postgres-clean/internal/models"
	"github.com/Shubham23jha/go-gin-postgres-clean/internal/repository"
	amqp "github.com/rabbitmq/amqp091-go"
)

type WorkerPool interface {
	Start(ctx context.Context)
}

type emailWorkerPool struct {
	repo       repository.CampaignRepository
	conn       *amqp.Connection
	ch         *amqp.Channel
	numWorkers int
	smtpHost   string
	smtpPort   string
	smtpUser   string
	smtpPass   string
	smtpSender string
}

func NewWorkerPool(repo repository.CampaignRepository, amqpURL string, numWorkers int) WorkerPool {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		log.Fatalf("failed to connect to RabbitMQ (Worker): %s", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("failed to open a channel (Worker): %s", err)
	}

	// 1. Declare DLX (Dead Letter Exchange)
	err = ch.ExchangeDeclare(
		"email_dlx", // name
		"direct",    // type
		true,        // durable
		false,       // auto-deleted
		false,       // internal
		false,       // no-wait
		nil,         // arguments
	)
	if err != nil {
		log.Fatalf("failed to declare DLX: %s", err)
	}

	// 2. Declare DLQ (Dead Letter Queue)
	_, err = ch.QueueDeclare(
		"email_dlq", // name
		true,        // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	if err != nil {
		log.Fatalf("failed to declare DLQ: %s", err)
	}

	// Bind DLQ to DLX
	err = ch.QueueBind("email_dlq", "email_failed", "email_dlx", false, nil)
	if err != nil {
		log.Fatalf("failed to bind DLQ: %s", err)
	}

	// 3. Declare Main Queue with DLX configuration
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
		log.Fatalf("failed to declare main queue: %s", err)
	}

	return &emailWorkerPool{
		repo:       repo,
		conn:       conn,
		ch:         ch,
		numWorkers: numWorkers,
		smtpHost:   os.Getenv("SMTP_HOST"),
		smtpPort:   os.Getenv("SMTP_PORT"),
		smtpUser:   os.Getenv("SMTP_USER"),
		smtpPass:   os.Getenv("SMTP_PASS"),
		smtpSender: os.Getenv("SMTP_SENDER"),
	}
}

func (p *emailWorkerPool) Start(ctx context.Context) {
	// ... (rest of Start remains same)
	msgs, err := p.ch.Consume(
		"email_queue", // queue
		"",            // consumer
		false,         // auto-ack
		false,         // exclusive
		false,         // no-local
		false,         // no-wait
		nil,           // args
	)
	if err != nil {
		log.Fatalf("failed to register a consumer: %s", err)
	}

	var wg sync.WaitGroup
	log.Printf("📥 Email Worker Pool started with %d workers (Connecting to %s)...", p.numWorkers, p.smtpHost)

	for i := 1; i <= p.numWorkers; i++ {
		wg.Add(1)
		go p.worker(i, msgs, &wg)
	}

	wg.Wait()
}

func (p *emailWorkerPool) worker(id int, msgs <-chan amqp.Delivery, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Printf("Worker %d: Waiting for messages...", id)

	for d := range msgs {
		log.Printf("Worker %d: Received a message (ID: %s)", id, d.MessageId)

		var payload map[string]string
		if err := json.Unmarshal(d.Body, &payload); err != nil {
			log.Printf("Worker %d: Failed to unmarshal payload: %s", id, err)
			d.Nack(false, false)
			continue
		}

		// 1. Idempotency Check
		processed, err := p.repo.IsMessageProcessed(d.MessageId)
		if err == nil && processed {
			log.Printf("Worker %d: Message %s already processed. Skipping.", id, d.MessageId)
			d.Ack(false)
			continue
		}

		// Real SMTP Delivery
		err = p.sendEmail(payload["recipient"], payload["subject"], payload["body"])
		
		status := "SUCCESS"
		errMsg := ""
		if err != nil {
			log.Printf("Worker %d: Failed to send email to %s: %s", id, payload["recipient"], err)
			
			// 3. Retry Logic with Exponential Backoff
			retryCount := 0
			if val, ok := d.Headers["x-retry-count"]; ok {
				retryCount = int(val.(int32))
			}

			if retryCount < 3 {
				retryCount++
				backoff := time.Duration(retryCount*retryCount) * time.Second
				log.Printf("Worker %d: Retrying in %v (Attempt %d/3)...", id, backoff, retryCount)
				
				// Sleep and Nack with Requeue
				time.Sleep(backoff)
				
				// Re-publish with incremented retry count
				p.ch.Publish(
					"",              // exchange
					d.RoutingKey,    // routing key
					false,           // mandatory
					false,           // immediate
					amqp.Publishing{
						Headers: amqp.Table{
							"x-retry-count": int32(retryCount),
						},
						ContentType:  d.ContentType,
						Body:         d.Body,
						MessageId:    d.MessageId,
						DeliveryMode: amqp.Persistent,
					},
				)
				d.Ack(false) // Ack the old one
				continue
			} else {
				log.Printf("Worker %d: Max retries reached for %s. Moving to DLQ.", id, d.MessageId)
				status = "FAILED"
				errMsg = err.Error()
				d.Nack(false, false) // Don't requeue -> Goes to DLX -> DLQ
			}
		} else {
			log.Printf("Worker %d: Successfully sent email to %s!", id, payload["recipient"])
			d.Ack(false)
		}

		// 4. Record Log in DB
		campaignID, _ := strconv.Atoi(payload["campaign_id"])
		emailLog := &models.EmailLog{
			CampaignID:   uint(campaignID),
			MessageID:    d.MessageId,
			Recipient:    payload["recipient"],
			Status:       status,
			ErrorMessage: errMsg,
		}
		p.repo.CreateEmailLog(emailLog)
	}
}

func (p *emailWorkerPool) sendEmail(to, subject, body string) error {
	auth := smtp.PlainAuth("", p.smtpUser, p.smtpPass, p.smtpHost)
	
	msg := fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s\r\n", p.smtpSender, to, subject, body)

	addr := fmt.Sprintf("%s:%s", p.smtpHost, p.smtpPort)
	return smtp.SendMail(addr, auth, p.smtpSender, []string{to}, []byte(msg))
}
