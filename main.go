// main.go
package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"notification-service/handlers"
	"notification-service/rabbitmq"
)

func main() {
	log.Println("Starting Go Notification Microservice...")

	// Create RabbitMQ connection
	conn, err := rabbitmq.NewConnection(DefaultRabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}
	defer conn.Close()

	// Declare exchange
	err = conn.DeclareExchange(ExchangeName, "topic")
	if err != nil {
		log.Fatalf("Failed to declare exchange: %s", err)
	}

	// Declare and bind queues
	setupQueues(conn)

	// Create notification handler
	notificationHandler := handlers.NewNotificationHandler()

	// Create consumer
	consumer := rabbitmq.NewConsumer(conn, notificationHandler.ProcessMessage)

	// Start consuming from all queues
	log.Println("Starting to consume messages...")
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
	_, err := conn.DeclareQueue(VolunteerPushQueue)
	if err != nil {
		log.Fatalf("Failed to declare '%s': %s", VolunteerPushQueue, err)
	}

	err = conn.BindQueue(VolunteerPushQueue, RoutingKeyAppStatusChanged, ExchangeName)
	if err != nil {
		log.Fatalf("Failed to bind '%s' for '%s': %s", VolunteerPushQueue, RoutingKeyAppStatusChanged, err)
	}

	err = conn.BindQueue(VolunteerPushQueue, RoutingKeyAppCancelled, ExchangeName)
	if err != nil {
		log.Fatalf("Failed to bind '%s' for '%s': %s", VolunteerPushQueue, RoutingKeyAppCancelled, err)
	}

	err = conn.BindQueue(VolunteerPushQueue, RoutingKeyOpportunityCreated, ExchangeName)
	if err != nil {
		log.Fatalf("Failed to bind '%s' for '%s': %s", VolunteerPushQueue, RoutingKeyOpportunityCreated, err)
	}

	// NGO Email Queue
	_, err = conn.DeclareQueue(NgoEmailQueue)
	if err != nil {
		log.Fatalf("Failed to declare '%s': %s", NgoEmailQueue, err)
	}

	err = conn.BindQueue(NgoEmailQueue, RoutingKeyApplicationNew, ExchangeName)
	if err != nil {
		log.Fatalf("Failed to bind '%s' for '%s': %s", NgoEmailQueue, RoutingKeyApplicationNew, err)
	}

	err = conn.BindQueue(NgoEmailQueue, RoutingKeyAppCancelled, ExchangeName)
	if err != nil {
		log.Fatalf("Failed to bind '%s' for '%s': %s", NgoEmailQueue, RoutingKeyAppCancelled, err)
	}

	// NGO Push Queue
	_, err = conn.DeclareQueue(NgoPushQueue)
	if err != nil {
		log.Fatalf("Failed to declare '%s': %s", NgoPushQueue, err)
	}

	err = conn.BindQueue(NgoPushQueue, RoutingKeyApplicationNew, ExchangeName)
	if err != nil {
		log.Fatalf("Failed to bind '%s' for '%s': %s", NgoPushQueue, RoutingKeyApplicationNew, err)
	}

	err = conn.BindQueue(NgoPushQueue, RoutingKeyAppCancelled, ExchangeName)
	if err != nil {
		log.Fatalf("Failed to bind '%s' for '%s': %s", NgoPushQueue, RoutingKeyAppCancelled, err)
	}
}

// waitForShutdown waits for a termination signal
func waitForShutdown() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Shutting down gracefully...")
	time.Sleep(2 * time.Second) // Give some time for resources to close
}
