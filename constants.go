// notification-service/constants/constants.go
package main

const (
	// RabbitMQ Connection Defaults
	DefaultRabbitMQURL = "amqp://guest:guest@localhost:5672/"

	// RabbitMQ Exchange Names
	ExchangeName = "notification_exchange"

	// RabbitMQ Queue Names
	VolunteerPushQueue = "volunteer_push_queue"
	NgoEmailQueue      = "ngo_email_queue"
	NgoPushQueue       = "ngo_push_queue"

	// RabbitMQ Routing Keys (Must match NestJS RabbitMQRoutingKey enum values)
	// These define how messages are routed to specific queues via the exchange.
	RoutingKeyApplicationNew     = "application.new"
	RoutingKeyAppStatusChanged   = "application.status_changed"
	RoutingKeyOpportunityCreated = "opportunity.created"
	RoutingKeyOpportunityUpdated = "opportunity.updated"
	RoutingKeyOpportunityDeleted = "opportunity.deleted"
	// Removed: RoutingKeyAppCancelled ("application.cancelled") as it's not used by NestJS producer.
)

// Notification Types (Must match NestJS RabbitMQEventType enum values)
// These are the values found *inside* the message payload's `notification_type` field.
// They are used in the `models/notification_message.go` for unmarshaling
// and in `handlers/notification_handler.go`'s switch statement for routing to specific handlers.
const (
	NotificationTypeNgoNewApplication        = "NGO_NEW_APPLICATION"
	NotificationTypeApplicationAccepted      = "APPLICATION_ACCEPTED"
	NotificationTypeApplicationRejected      = "APPLICATION_REJECTED"
	NotificationTypeApplicationWithdrawn     = "APPLICATION_WITHDRAWN" // Matches NestJS enum
	NotificationTypeApplicationCompleted     = "APPLICATION_COMPLETED"
	NotificationTypeVolunteerAppStatusUpdate = "VOLUNTEER_APPLICATION_STATUS_UPDATE" // Matches NestJS enum
	NotificationTypeVolunteerNewOpportunity  = "VOLUNTEER_NEW_MATCHING_OPPORTUNITY"  // Matches NestJS enum
	NotificationTypeOpportunityUpdated       = "OPPORTUNITY_UPDATED"
	NotificationTypeOpportunityDeleted       = "OPPORTUNITY_DELETED"
	// Removed: NotificationTypeNgoAppCancelled ("NGO_APPLICATION_CANCELLED") as it doesn't match a NestJS event type.
)
