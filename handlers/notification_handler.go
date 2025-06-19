// handlers/notification_handler.go
package handlers

import (
	"encoding/json"
	"log"

	"notification-service/models" // Make sure this path is correct for your models
	// You will eventually import your actual email and FCM service packages here, e.g.:
	// "notification-service/services/email"
	// "notification-service/services/fcm"
)

// NotificationHandler handles incoming notification messages
type NotificationHandler struct {
	// You should add fields for your actual email and push notification services here.
	// For example:
	// EmailService *email.Service
	// FCMService *fcm.Service
}

// NewNotificationHandler creates a new notification handler
func NewNotificationHandler() *NotificationHandler {
	return &NotificationHandler{
		// Initialize your services here when they are implemented, e.g.:
		// EmailService: email.NewService(),
		// FCMService: fcm.NewService(),
	}
}

// ProcessMessage processes a notification message by unmarshaling it
// and routing it to the appropriate handler based on its type.
func (h *NotificationHandler) ProcessMessage(body []byte) error {
	var msg models.NotificationMessage

	err := json.Unmarshal(body, &msg)
	if err != nil {
		log.Printf("Failed to unmarshal JSON message: %v. Raw Body: %s", err, body)
		return err
	}

	log.Printf("Successfully unmarshaled message: NotificationType=%s, RecipientID=%s, PushPref=%t, EmailPref=%t",
		msg.NotificationType, msg.Recipient.UserID, msg.Recipient.Prefs.ReceivePush, msg.Recipient.Prefs.ReceiveEmail)

	// Process based on notification type
	switch msg.NotificationType {
	// --- Volunteer-centric Application Status Updates ---
	case "APPLICATION_ACCEPTED":
		return h.handleApplicationStatusUpdate(msg)
	case "APPLICATION_REJECTED":
		return h.handleApplicationStatusUpdate(msg)
	case "APPLICATION_COMPLETED":
		return h.handleApplicationStatusUpdate(msg)
	case "APPLICATION_STATUS_CHANGED": // General fallback for any other status change to volunteer
		return h.handleApplicationStatusUpdate(msg)

	// --- NGO-centric Application Events ---
	case "APPLICATION_WITHDRAWN": // This event is directed to the NGO
		return h.handleNgoApplicationEvent(msg)
	case "NGO_NEW_APPLICATION":
		return h.handleNgoNewApplication(msg) // Existing handler for new applications to NGOs

	// --- Other Specific Notification Types ---
	case "VOLUNTEER_NEW_MATCHING_OPPORTUNITY":
		return h.handleVolunteerNewOpportunity(msg)

		// --- Opportunity Management Notifications ---
	case "OPPORTUNITY_UPDATED":
		log.Printf("Handling Opportunity Updated: %s", msg.Payload.OpportunityTitle)
		return h.handleOpportunityUpdate(msg)
	case "OPPORTUNITY_DELETED":
		log.Printf("Handling Opportunity Deleted: %s", msg.Payload.OpportunityTitle)
		return h.handleOppotunityDeleted(msg)

	default:
		log.Printf("Unknown notification type received: %s. Raw Message: %s", msg.NotificationType, body)
		// Consider returning an error for unhandled types in a production system
		// return fmt.Errorf("unknown notification type: %s", msg.NotificationType)
	}

	return nil
}

// handleApplicationStatusUpdate processes notifications for volunteers about their application status changes.
// This single function handles ACCEPTED, REJECTED, COMPLETED, and general STATUS_CHANGED notifications.
func (h *NotificationHandler) handleApplicationStatusUpdate(msg models.NotificationMessage) error {
	log.Printf("Handling Volunteer Application Status Update: Type=%s, AppID=%d, OldStatus=%s, NewStatus=%s, VolunteerName=%s",
		msg.NotificationType, msg.Payload.ApplicationID, msg.Payload.OldStatus, msg.Payload.NewStatus, msg.Payload.VolunteerName)

	// Extract common notification content from payload
	title := msg.Payload.Title
	body := msg.Payload.Body
	subject := msg.Payload.Subject   // Used for email
	deepLink := msg.Payload.DeepLink // Declared and assigned

	// --- Email Logic for Volunteer ---
	if msg.Recipient.Prefs.ReceiveEmail && msg.Recipient.EmailAddress != "" {
		log.Printf("Attempting to send Email to Volunteer %s (%s) for status '%s'. Subject: '%s'",
			msg.Recipient.UserID, msg.Recipient.EmailAddress, msg.Payload.NewStatus, subject)
		// TODO: Call your actual Email sending function here, e.g.:
		// if h.EmailService != nil {
		//     err := h.EmailService.SendEmail(msg.Recipient.EmailAddress, subject, body /* , optional HTML body */)
		//     if err != nil {
		//         log.Printf("Error sending email to %s: %v", msg.Recipient.EmailAddress, err)
		//     }
		// }
	} else {
		log.Printf("Skipping Email for Volunteer %s (Pref: %t, Email: %t). Status: %s",
			msg.Recipient.UserID, msg.Recipient.Prefs.ReceiveEmail, msg.Recipient.EmailAddress != "", msg.Payload.NewStatus)
	}

	// --- Push Notification Logic for Volunteer ---
	if msg.Recipient.Prefs.ReceivePush && msg.Recipient.DeviceToken != "" {
		// Used deepLink in log.Printf
		log.Printf("Attempting to send Push to Volunteer %s (%s) for status '%s'. Title: '%s', DeepLink: '%s', Body: '%s'",
			msg.Recipient.UserID, msg.Recipient.DeviceToken, msg.Payload.NewStatus, title, deepLink, body)
		// TODO: Call your actual FCM sending function here, e.g.:
		// if h.FCMService != nil {
		//     err := h.FCMService.SendPushNotification(msg.Recipient.DeviceToken, title, body, deepLink)
		//     if err != nil {
		//         log.Printf("Error sending push notification to %s: %v", msg.Recipient.UserID, err)
		//     }
		// }
	} else {
		log.Printf("Skipping Push for Volunteer %s (Pref: %t, Token: %t). Status: %s",
			msg.Recipient.UserID, msg.Recipient.Prefs.ReceivePush, msg.Recipient.DeviceToken != "", msg.Payload.NewStatus)
	}
	return nil
}

// handleNgoApplicationEvent processes notifications for NGOs about application events (e.g., withdrawn).
func (h *NotificationHandler) handleNgoApplicationEvent(msg models.NotificationMessage) error {
	log.Printf("Handling NGO Application Event: Type=%s, AppID=%d, VolunteerName=%s, OpportunityTitle=%s",
		msg.NotificationType, msg.Payload.ApplicationID, msg.Payload.VolunteerName, msg.Payload.OpportunityTitle)

	// Extract common notification content from payload
	title := msg.Payload.Title
	body := msg.Payload.Body
	subject := msg.Payload.Subject   // Used for email
	deepLink := msg.Payload.DeepLink // Declared and assigned

	// --- Email Logic for NGO ---
	if msg.Recipient.Prefs.ReceiveEmail && msg.Recipient.EmailAddress != "" {
		log.Printf("Attempting to send Email to NGO %s (%s) for event '%s'. Subject: '%s'",
			msg.Recipient.UserID, msg.Recipient.EmailAddress, msg.NotificationType, subject)
		// TODO: Call your actual Email sending function here
		// if h.EmailService != nil { ... }
	} else {
		log.Printf("Skipping Email for NGO %s (Pref: %t, Email: %t). Event: %s",
			msg.Recipient.UserID, msg.Recipient.Prefs.ReceiveEmail, msg.Recipient.EmailAddress != "", msg.NotificationType)
	}

	// --- Push Notification Logic for NGO ---
	if msg.Recipient.Prefs.ReceivePush && msg.Recipient.DeviceToken != "" {
		// Used deepLink in log.Printf
		log.Printf("Attempting to send Push to NGO %s (%s) for event '%s'. Title: '%s', DeepLink: '%s', Body: '%s'",
			msg.Recipient.UserID, msg.Recipient.DeviceToken, msg.NotificationType, title, deepLink, body)
		// TODO: Call your actual FCM sending function here
		// if h.FCMService != nil { ... }
	} else {
		log.Printf("Skipping Push for NGO %s (Pref: %t, Token: %t). Event: %s",
			msg.Recipient.UserID, msg.Recipient.Prefs.ReceivePush, msg.Recipient.DeviceToken != "", msg.NotificationType)
	}
	return nil
}

// handleNgoNewApplication handles new applications for NGOs.
func (h *NotificationHandler) handleNgoNewApplication(msg models.NotificationMessage) error {
	log.Printf("Handling NGO New Application: Type=%s, AppID=%d, VolunteerName=%s, OpportunityTitle=%s",
		msg.NotificationType, msg.Payload.ApplicationID, msg.Payload.VolunteerName, msg.Payload.OpportunityTitle)

	title := msg.Payload.Title
	body := msg.Payload.Body
	subject := msg.Payload.Subject
	deepLink := msg.Payload.DeepLink // Declared and assigned

	// --- Email Logic for NGO ---
	if msg.Recipient.Prefs.ReceiveEmail && msg.Recipient.EmailAddress != "" {
		log.Printf("Attempting to send Email for NGO: %s - Subject: '%s'",
			msg.Recipient.UserID, subject)
		// TODO: Call actual Email sending function here
	} else {
		log.Printf("Skipping Email for NGO %s (Pref: %t, Email: %t)",
			msg.Recipient.UserID, msg.Recipient.Prefs.ReceiveEmail, msg.Recipient.EmailAddress != "")
	}

	// --- Push Notification Logic for NGO ---
	if msg.Recipient.Prefs.ReceivePush && msg.Recipient.DeviceToken != "" {
		// Used deepLink in log.Printf
		log.Printf("Attempting to send Push for NGO: %s - Title: '%s', DeepLink: '%s', Body: '%s'",
			msg.Recipient.UserID, title, deepLink, body)
		// TODO: Call actual FCM sending function here for NGO
	} else {
		log.Printf("Skipping Push for NGO %s (Pref: %t, Token: %t)",
			msg.Recipient.UserID, msg.Recipient.Prefs.ReceivePush, msg.Recipient.DeviceToken != "")
	}
	return nil
}

// handleVolunteerNewOpportunity handles notifications for volunteers about new matching opportunities.
func (h *NotificationHandler) handleVolunteerNewOpportunity(msg models.NotificationMessage) error {
	log.Printf("Handling Volunteer New Matching Opportunity: Type=%s, OpportunityTitle=%s",
		msg.NotificationType, msg.Payload.OpportunityTitle)

	title := msg.Payload.Title
	body := msg.Payload.Body
	subject := msg.Payload.Subject
	deepLink := msg.Payload.DeepLink // Declared and assigned

	// For new opportunities, assume only push notifications for volunteers (or add email if desired)
	if msg.Recipient.Prefs.ReceivePush && msg.Recipient.DeviceToken != "" {
		// Used deepLink in log.Printf
		log.Printf("Attempting to send Push for Volunteer (New Matching Opportunity): %s - Title: '%s', DeepLink: '%s', Body: '%s'",
			msg.Recipient.UserID, title, deepLink, body)
		// TODO: Call actual FCM sending function here
	} else {
		log.Printf("Skipping Push for Volunteer %s (Pref: %t, Token: %t). Type: %s",
			msg.Recipient.UserID, msg.Recipient.Prefs.ReceivePush, msg.Recipient.DeviceToken != "", msg.NotificationType)
	}

	// If email should also be sent for new opportunities:
	if msg.Recipient.Prefs.ReceiveEmail && msg.Recipient.EmailAddress != "" {
		log.Printf("Attempting to send Email for Volunteer (New Matching Opportunity): %s - Subject: '%s'",
			msg.Recipient.UserID, subject)
		// TODO: Call actual Email sending function here
	} else {
		log.Printf("Skipping Email for Volunteer %s (Pref: %t, Email: %t). Type: %s",
			msg.Recipient.UserID, msg.Recipient.Prefs.ReceiveEmail, msg.Recipient.EmailAddress != "", msg.NotificationType)
	}
	return nil
}

// handleOpportunityUpdate handles notifications for updates to opportunities.
func (h *NotificationHandler) handleOpportunityUpdate(msg models.NotificationMessage) error {
	log.Printf("Handling Opportunity Update: Type=%s, OpportunityID=%d, Title=%s",
		msg.NotificationType, msg.Payload.OpportunityID, msg.Payload.OpportunityTitle)

	title := msg.Payload.Title
	body := msg.Payload.Body
	subject := msg.Payload.Subject
	deepLink := msg.Payload.DeepLink // Declared and assigned

	// --- Email Logic for NGO ---
	if msg.Recipient.Prefs.ReceiveEmail && msg.Recipient.EmailAddress != "" {
		log.Printf("Attempting to send Email for Opportunity Update to NGO: %s - Subject: '%s'",
			msg.Recipient.UserID, subject)
		// TODO: Call actual Email sending function here
	} else {
		log.Printf("Skipping Email for NGO %s (Pref: %t, Email: %t). Type: %s",
			msg.Recipient.UserID, msg.Recipient.Prefs.ReceiveEmail, msg.Recipient.EmailAddress != "", msg.NotificationType)
	}
	// --- Push Notification Logic for NGO ---
	if msg.Recipient.Prefs.ReceivePush && msg.Recipient.DeviceToken != "" {
		// Used deepLink in log.Printf
		log.Printf("Attempting to send Push for Opportunity Update to NGO: %s - Title: '%s', DeepLink: '%s', Body: '%s'",
			msg.Recipient.UserID, title, deepLink, body)
		// TODO: Call actual FCM sending function here
	} else {
		log.Printf("Skipping Push for NGO %s (Pref: %t, Token: %t). Type: %s",
			msg.Recipient.UserID, msg.Recipient.Prefs.ReceivePush, msg.Recipient.DeviceToken != "", msg.NotificationType)
	}

	return nil
}

// handleOppotunityDeleted handles notifications for deleted opportunities.
func (h *NotificationHandler) handleOppotunityDeleted(msg models.NotificationMessage) error {
	log.Printf("Handling Opportunity Deleted: Type=%s, OpportunityID=%d, Title=%s",
		msg.NotificationType, msg.Payload.OpportunityID, msg.Payload.OpportunityTitle)

	title := msg.Payload.Title
	body := msg.Payload.Body
	subject := msg.Payload.Subject
	deepLink := msg.Payload.DeepLink // Declared and assigned

	// --- Email Logic for NGO ---
	if msg.Recipient.Prefs.ReceiveEmail && msg.Recipient.EmailAddress != "" {
		log.Printf("Attempting to send Email for Opportunity Deletion to Volunteers: %s - Subject: '%s'",
			msg.Recipient.UserID, subject)
		// TODO: Call actual Email sending function here
	} else {
		log.Printf("Skipping Email for NGO %s (Pref: %t, Email: %t). Type: %s",
			msg.Recipient.UserID, msg.Recipient.Prefs.ReceiveEmail, msg.Recipient.EmailAddress != "", msg.NotificationType)
	}

	// --- Push Notification Logic for Volunteers ---
	if msg.Recipient.Prefs.ReceivePush && msg.Recipient.DeviceToken != "" {
		// Used deepLink in log.Printf
		log.Printf("Attempting to send Push for Opportunity Deletion to Volunteers: %s - Title: '%s', DeepLink: '%s', Body: '%s'",
			msg.Recipient.UserID, title, deepLink, body)
		// TODO: Call actual FCM sending function here
	} else {
		log.Printf("Skipping Push for NGO %s (Pref: %t, Token: %t). Type: %s",
			msg.Recipient.UserID, msg.Recipient.Prefs.ReceivePush, msg.Recipient.DeviceToken != "", msg.NotificationType)
	}
	return nil
}
