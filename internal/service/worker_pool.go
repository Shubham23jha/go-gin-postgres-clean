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
	amqpURL    string
	conn       *amqp.Connection
	ch         *amqp.Channel
	numWorkers int
	smtpHost   string
	smtpPort   string
	smtpUser   string
	smtpPass   string
	smtpSender string
	mu         sync.Mutex
}

func NewWorkerPool(repo repository.CampaignRepository, amqpURL string, numWorkers int) WorkerPool {
	return &emailWorkerPool{
		repo:       repo,
		amqpURL:    amqpURL,
		numWorkers: numWorkers,
		smtpHost:   os.Getenv("SMTP_HOST"),
		smtpPort:   os.Getenv("SMTP_PORT"),
		smtpUser:   os.Getenv("SMTP_USER"),
		smtpPass:   os.Getenv("SMTP_PASS"),
		smtpSender: os.Getenv("SMTP_SENDER"),
	}
}

func (p *emailWorkerPool) connect() error {
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

	// 1. Declare DLX
	err = ch.ExchangeDeclare("email_dlx", "direct", true, false, false, false, nil)
	if err != nil {
		ch.Close()
		conn.Close()
		return err
	}

	// 2. Declare DLQ
	_, err = ch.QueueDeclare("email_dlq", true, false, false, false, nil)
	if err != nil {
		ch.Close()
		conn.Close()
		return err
	}

	// Bind DLQ
	err = ch.QueueBind("email_dlq", "email_failed", "email_dlx", false, nil)
	if err != nil {
		ch.Close()
		conn.Close()
		return err
	}

	// 3. Declare Main Queue
	args := amqp.Table{
		"x-dead-letter-exchange":    "email_dlx",
		"x-dead-letter-routing-key": "email_failed",
	}
	_, err = ch.QueueDeclare("email_queue", true, false, false, false, args)
	if err != nil {
		ch.Close()
		conn.Close()
		return err
	}

	p.conn = conn
	p.ch = ch
	return nil
}

func (p *emailWorkerPool) Start(ctx context.Context) {
	for {
		if err := p.connect(); err != nil {
			log.Printf("❌ Worker Pool failed to connect to RabbitMQ: %s. Retrying in 5s...", err)
			time.Sleep(5 * time.Second)
			continue
		}

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
			log.Printf("❌ Failed to register a consumer: %s. Retrying...", err)
			time.Sleep(5 * time.Second)
			continue
		}

		var wg sync.WaitGroup
		log.Printf("📥 Email Worker Pool started with %d workers (Connecting to %s)...", p.numWorkers, p.smtpHost)

		workerCtx, cancel := context.WithCancel(ctx)

		for i := 1; i <= p.numWorkers; i++ {
			wg.Add(1)
			go p.worker(workerCtx, i, msgs, &wg)
		}

		// Watch for connection closure or context done
		closeErr := <-p.conn.NotifyClose(make(chan *amqp.Error))
		log.Printf("⚠️ RabbitMQ connection closed: %v", closeErr)
		cancel()
		wg.Wait()

		select {
		case <-ctx.Done():
			return
		default:
			log.Println("♻️ Attempting to restart worker pool after connection loss...")
		}
	}
}

func (p *emailWorkerPool) worker(ctx context.Context, id int, msgs <-chan amqp.Delivery, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Printf("Worker %d: Waiting for messages...", id)

	for {
		select {
		case <-ctx.Done():
			return
		case d, ok := <-msgs:
			if !ok {
				return
			}
			p.handleDelivery(ctx, id, d)
		}
	}
}

func (p *emailWorkerPool) handleDelivery(ctx context.Context, id int, d amqp.Delivery) {
	log.Printf("Worker %d: Received a message (ID: %s)", id, d.MessageId)

	var payload map[string]string
	if err := json.Unmarshal(d.Body, &payload); err != nil {
		log.Printf("Worker %d: Failed to unmarshal payload: %s", id, err)
		d.Nack(false, false)
		return
	}

	// 1. Idempotency Check
	processed, err := p.repo.IsMessageProcessed(d.MessageId)
	if err == nil && processed {
		log.Printf("Worker %d: Message %s already processed. Skipping.", id, d.MessageId)
		d.Ack(false)
		return
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
			switch v := val.(type) {
			case int:
				retryCount = v
			case int32:
				retryCount = int(v)
			case float64:
				retryCount = int(v)
			}
		}

		if retryCount < 3 {
			retryCount++
			backoff := time.Duration(retryCount*retryCount) * time.Second
			log.Printf("Worker %d: Retrying in %v (Attempt %d/3)...", id, backoff, retryCount)

			time.Sleep(backoff)

			// Re-publish with incremented retry count
			pubErr := p.ch.PublishWithContext(ctx,
				"",           // exchange
				d.RoutingKey, // routing key
				false,        // mandatory
				false,        // immediate
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
			if pubErr != nil {
				log.Printf("Worker %d: Failed to re-publish for retry: %s", id, pubErr)
				d.Nack(false, true) // Requeue to original queue if possible
			} else {
				d.Ack(false)
			}
			return
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
