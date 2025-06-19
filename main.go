package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"notification-service/handlers"
	"notification-service/rabbitmq"
	// "notification-service/services/email"
	// "notification-service/services/fcm"
)

func main() {
	log.Println("Starting Go Notification Microservice...")

	// Create RabbitMQ connection
	// Use constant from the new package
	conn, err := rabbitmq.NewConnection(DefaultRabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}
	defer conn.Close()

	// Declare exchange
	// Use constant from the new package
	err = conn.DeclareExchange(ExchangeName, "topic")
	if err != nil {
		log.Fatalf("Failed to declare exchange: %s", err)
	}

	// Declare and bind queues
	setupQueues(conn) // This function will also use
	// // Initialize Email and FCM Services
	// emailService := email.NewService()
	// fcmService := fcm.NewService()

	// // Create notification handler, passing the initialized services
	notificationHandler := handlers.NewNotificationHandler()

	// Create consumer
	consumer := rabbitmq.NewConsumer(conn, notificationHandler.ProcessMessage)

	// Start consuming from all queues
	log.Println("Starting to consume messages...")
	// Use from the new package
	err = consumer.StartConsuming(VolunteerPushQueue)
	if err != nil {
		log.Fatalf("Failed to register volunteer push consumer: %s", err)
	}

	err = consumer.StartConsuming(NgoEmailQueue)
	if err != nil {
		log.Fatalf("Failed to register NGO email consumer: %s", err)
	}

	err = consumer.StartConsuming(NgoPushQueue)
	if err != nil {
		log.Fatalf("Failed to register NGO push consumer: %s", err)
	}

	log.Println("Go Notification Microservice started. Waiting for messages. To exit, press CTRL+C")

	// Wait for termination signal
	waitForShutdown()
}

// setupQueues declares and binds all necessary queues
func setupQueues(conn *rabbitmq.Connection) {
	// Volunteer Push Queue
	// Use from the new package
	_, err := conn.DeclareQueue(VolunteerPushQueue)
	if err != nil {
		log.Fatalf("Failed to declare '%s': %s", VolunteerPushQueue, err)
	}

	err = conn.BindQueue(VolunteerPushQueue, RoutingKeyAppStatusChanged, ExchangeName)
	if err != nil {
		log.Fatalf("Failed to bind '%s' for '%s': %s", VolunteerPushQueue, RoutingKeyAppStatusChanged, err)
	}

	// Removed binding for RoutingKeyAppCancelled, as it's not used by NestJS.
	// If you later decide to use a separate routing key for application cancellations,
	// you would add it back here and in go.

	err = conn.BindQueue(VolunteerPushQueue, RoutingKeyOpportunityCreated, ExchangeName)
	if err != nil {
		log.Fatalf("Failed to bind '%s' for '%s': %s", VolunteerPushQueue, RoutingKeyOpportunityCreated, err)
	}

	// NEW: Bind for Opportunity Deleted and Updated
	err = conn.BindQueue(VolunteerPushQueue, RoutingKeyOpportunityDeleted, ExchangeName)
	if err != nil {
		log.Fatalf("Failed to bind '%s' for '%s': %s", VolunteerPushQueue, RoutingKeyOpportunityDeleted, err)
	}
	err = conn.BindQueue(VolunteerPushQueue, RoutingKeyOpportunityUpdated, ExchangeName)
	if err != nil {
		log.Fatalf("Failed to bind '%s' for '%s': %s", VolunteerPushQueue, RoutingKeyOpportunityUpdated, err)
	}

	// NGO Email Queue
	// Use from the new package
	_, err = conn.DeclareQueue(NgoEmailQueue)
	if err != nil {
		log.Fatalf("Failed to declare '%s': %s", NgoEmailQueue, err)
	}

	err = conn.BindQueue(NgoEmailQueue, RoutingKeyApplicationNew, ExchangeName)
	if err != nil {
		log.Fatalf("Failed to bind '%s' for '%s': %s", NgoEmailQueue, RoutingKeyApplicationNew, err)
	}

	// Removed binding for RoutingKeyAppCancelled
	// err = conn.BindQueue(NgoEmailQueue, RoutingKeyAppCancelled, ExchangeName)
	// if err != nil {
	// 	log.Fatalf("Failed to bind '%s' for '%s': %s", NgoEmailQueue, RoutingKeyAppCancelled, err)
	// }

	// NGO Push Queue
	// Use from the new package
	_, err = conn.DeclareQueue(NgoPushQueue)
	if err != nil {
		log.Fatalf("Failed to declare '%s': %s", NgoPushQueue, err)
	}

	err = conn.BindQueue(NgoPushQueue, RoutingKeyApplicationNew, ExchangeName)
	if err != nil {
		log.Fatalf("Failed to bind '%s' for '%s': %s", NgoPushQueue, RoutingKeyApplicationNew, err)
	}

	// Removed binding for RoutingKeyAppCancelled
	// err = conn.BindQueue(NgoPushQueue, RoutingKeyAppCancelled, ExchangeName)
	// if err != nil {
	// 	log.Fatalf("Failed to bind '%s' for '%s': %s", NgoPushQueue, RoutingKeyAppCancelled, err)
	// }
}

// waitForShutdown waits for a termination signal
func waitForShutdown() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Shutting down gracefully...")
	time.Sleep(2 * time.Second) // Give some time for resources to close
}
