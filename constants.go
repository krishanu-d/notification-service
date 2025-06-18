// constants.go
package main

const (
	// Exchange and Queue Names
	ExchangeName       = "notification_exchange"
	VolunteerPushQueue = "volunteer_push_queue"
	NgoEmailQueue      = "ngo_email_queue"
	NgoPushQueue       = "ngo_push_queue"

	// RabbitMQ Defaults
	DefaultRabbitMQURL = "amqp://guest:guest@localhost:5672/"

	// Routing Keys
	RoutingKeyAppStatusChanged   = "application.status_changed"
	RoutingKeyAppCancelled       = "application.cancelled"
	RoutingKeyOpportunityCreated = "opportunity.created"
	RoutingKeyApplicationNew     = "application.new"

	// Notification Types
	NotificationTypeVolunteerAppStatusUpdate = "VOLUNTEER_APPLICATION_STATUS_UPDATE"
	NotificationTypeNgoNewApplication        = "NGO_NEW_APPLICATION"
	NotificationTypeNgoAppCancelled          = "NGO_APPLICATION_CANCELLED"
	NotificationTypeVolunteerNewOpportunity  = "VOLUNTEER_NEW_MATCHING_OPPORTUNITY"
)
