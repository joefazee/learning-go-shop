package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-aws/sqs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/joefazee/learning-go-shop/internal/config"
	"github.com/joefazee/learning-go-shop/internal/models"
	"github.com/joefazee/learning-go-shop/internal/notifications"
	"github.com/joefazee/learning-go-shop/internal/providers"
)

func main() {
	log.Println("Starting notification service...")

	ctx := context.Background()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize email notifier
	emailConfig := &notifications.SMTPConfig{
		Host:     cfg.SMTP.Host,
		Port:     cfg.SMTP.Port,
		Username: cfg.SMTP.Username,
		Password: cfg.SMTP.Password,
		From:     cfg.SMTP.From,
	}

	emailNotifier := notifications.NewEmailNotifier(emailConfig)

	// Create AWS config for SQS
	awsConfig, err := providers.CreateAWSConfig(ctx, cfg.AWS.S3Endpoint, cfg.AWS.Region)
	if err != nil {
		log.Fatalf("Failed to create AWS config: %v", err)
	}

	// Create SQS subscriber
	logger := watermill.NewStdLogger(false, false)
	subscriber, err := sqs.NewSubscriber(sqs.SubscriberConfig{
		AWSConfig: awsConfig,
	}, logger)

	if err != nil {
		log.Fatalf("Failed to create subscriber: %v", err)
	}

	// Subscribe to messages
	messages, err := subscriber.Subscribe(ctx, cfg.AWS.EventQueueName)
	if err != nil {
		subscriber.Close()
		log.Fatalf("Failed to subscribe to queue: %v", err)
	}

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	log.Println("Notification service started. Waiting for messages...")

	for {
		select {
		case msg := <-messages:
			if err := processMessage(msg, emailNotifier); err != nil {
				log.Printf("Error processing message: %v", err)
				msg.Nack()
			} else {
				msg.Ack()
			}
		case <-sigChan:
			log.Println("Shutting down notification service...")
			subscriber.Close()
			return
		}
	}
}

func processMessage(msg *message.Message, emailNotifier *notifications.EmailNotifier) error {
	eventType := msg.Metadata.Get("event_type")
	switch eventType {
	case notifications.UserLoggedIn:
		return handleUserLoggedIn(msg, emailNotifier)
	default:
		log.Printf("Unknown event type: %s", eventType)
		return nil
	}
}

func handleUserLoggedIn(msg *message.Message, emailNotifier *notifications.EmailNotifier) error {
	var user models.User
	if err := json.Unmarshal(msg.Payload, &user); err != nil {
		return err
	}

	userName := user.FirstName + " " + user.LastName
	if userName == " " {
		userName = "User"
	}

	log.Printf("Sending login notification to %s", user.Email)

	return emailNotifier.SendLoginNotification(user.Email, userName)
}
