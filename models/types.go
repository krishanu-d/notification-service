// models/types.go
package models

// Recipient defines the structure for the notification target user's details.
// This data will be provided by the NestJS backend within the message.
type Recipient struct {
	UserID       string `json:"user_id"`                 // Unique ID of the user (volunteer or NGO)
	PlatformType string `json:"platform_type,omitempty"` // e.g., "mobile", "web" (for push)
	DeviceToken  string `json:"device_token,omitempty"`  // FCM token for push notifications
	EmailAddress string `json:"email_address,omitempty"` // Email address for email notifications
	// PhoneNumber   string `json:"phone_number,omitempty"` // Uncomment if you add SMS later

	// Prefs contains the user's general notification preferences.
	// These are also fetched by NestJS and included here.
	Prefs struct {
		ReceivePush  bool `json:"receive_push"`
		ReceiveEmail bool `json:"receive_email"`
		// You can add more granular preferences here if needed later
	} `json:"prefs"`
}

// Payload defines the actual content of the notification.
// This content is prepared by the NestJS backend.
type Payload struct {
	Title        string `json:"title,omitempty"`
	Body         string `json:"body,omitempty"`
	Subject      string `json:"subject,omitempty"`
	BodyHTML     string `json:"body_html,omitempty"`
	DeepLink     string `json:"deep_link,omitempty"`
	TemplateName string `json:"template_name,omitempty"`

	// Specific data related to the event, included for context/personalization.
	// These are optional and depend on the 'notification_type'.
	ApplicationID    int    `json:"application_id,omitempty"`
	OpportunityID    int    `json:"opportunity_id,omitempty"`
	NGOID            int    `json:"ngo_id,omitempty"`
	VolunteerID      int    `json:"volunteer_id,omitempty"`
	OldStatus        string `json:"old_status,omitempty"`
	NewStatus        string `json:"new_status,omitempty"`
	OpportunityTitle string `json:"opportunity_title,omitempty"`
	VolunteerName    string `json:"volunteer_name,omitempty"`
	NgoName          string `json:"ngo_name,omitempty"`
}

// NotificationMessage is the top-level struct for an incoming message from RabbitMQ.
// This is the full structure of the JSON message expected from the NestJS backend.
type NotificationMessage struct {
	// notification_type specifies what kind of business event this notification is for.
	// This helps the Go service route to the correct internal logic.
	NotificationType string    `json:"notification_type"`        // e.g., "VOLUNTEER_APP_STATUS_UPDATE", "NGO_NEW_APPLICATION"
	Recipient        Recipient `json:"recipient"`                // Details about who receives the notification
	Payload          Payload   `json:"payload"`                  // The actual content to be delivered
	SenderService    string    `json:"sender_service,omitempty"` // Optional: for auditing/debugging, e.g., "VolHub_ApplicationsService"
	Timestamp        int64     `json:"timestamp,omitempty"`      // Optional: when the event occurred (Unix timestamp)
}
