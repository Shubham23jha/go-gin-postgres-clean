package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Stats struct {
	QueueLength   int `json:"queue_length"`
	ActiveWorkers int `json:"active_workers"`
}

var (
	clientset *kubernetes.Clientset
	amqpConn  *amqp.Connection
	amqpCh    *amqp.Channel
)

func main() {
	// 1. Setup Kubernetes Client
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Printf("Not running in cluster, falling back to local proxy or skipping K8s stats: %v", err)
	} else {
		clientset, err = kubernetes.NewForConfig(config)
		if err != nil {
			log.Fatalf("Error creating k8s client: %v", err)
		}
	}

	// 2. Setup RabbitMQ Connection
	amqpURL := os.Getenv("RABBITMQ_URL")
	if amqpURL == "" {
		amqpURL = "amqp://guest:guest@localhost:5672/"
	}
	
	connectRabbitMQ(amqpURL)
	defer amqpConn.Close()
	defer amqpCh.Close()

	// 3. Serve Frontend and API
	http.Handle("/", http.FileServer(http.Dir("./web")))
	http.HandleFunc("/api/stats", statsHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("📊 Dashboard Service starting on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func connectRabbitMQ(url string) {
	var err error
	for i := 0; i < 5; i++ {
		amqpConn, err = amqp.Dial(url)
		if err == nil {
			break
		}
		log.Printf("Failed to connect to RabbitMQ, retrying in 5s... (%d/5)", i+1)
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		log.Fatalf("Could not connect to RabbitMQ: %v", err)
	}

	amqpCh, err = amqpConn.Channel()
	if err != nil {
		log.Fatalf("Could not open RabbitMQ channel: %v", err)
	}
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	// Set headers for SSE (Server-Sent Events)
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-r.Context().Done():
			return
		case <-ticker.C:
			stats := getStats()
			data, _ := json.Marshal(stats)
			fmt.Fprintf(w, "data: %s\n\n", data)
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
		}
	}
}

func getStats() Stats {
	stats := Stats{}

	// Get Queue Length
	q, err := amqpCh.QueueInspect("email_queue")
	if err == nil {
		stats.QueueLength = q.Messages
	} else {
		log.Printf("Error inspecting queue: %v", err)
	}

	// Get Active Workers from K8s
	if clientset != nil {
		deploy, err := clientset.AppsV1().Deployments("default").Get(context.TODO(), "email-worker", metav1.GetOptions{})
		if err == nil {
			stats.ActiveWorkers = int(deploy.Status.Replicas)
		} else {
			log.Printf("Error getting worker deployment: %v", err)
		}
	} else {
		// Mock data for local testing if needed, or just 0
		stats.ActiveWorkers = 0 
	}

	return stats
}
