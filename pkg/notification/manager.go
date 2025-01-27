package notification

import (
	"fmt"
)

// NotificationSystem represents a type of notification system (e.g., email, SMS, Slack).
type NotificationSystem string

// NoticeType represents a type of notification (e.g., "welcome", "password_reset").
type NoticeType string

const (
	EmailSystem NotificationSystem = "email"
	SMSSystem   NotificationSystem = "sms"
	SlackSystem NotificationSystem = "slack"

	ExampleNotice NoticeType = "example"
)

type NoticeTemplate struct {
	Subject  string
	Body     string
	BodyPath string
}

// NotificationManager manages notifiers and notification templates.
type NotificationManager struct {
	notifiers            map[NotificationSystem]Notifier                      // Map of notification systems to their Notifier implementations
	notificationRegistry map[NoticeType]map[NotificationSystem]NoticeTemplate // Registry for notification templates
}

// NewNotificationManager creates and returns a new NotificationManager.
func NewNotificationManager() *NotificationManager {
	return &NotificationManager{
		notifiers:            make(map[NotificationSystem]Notifier),
		notificationRegistry: make(map[NoticeType]map[NotificationSystem]NoticeTemplate),
	}
}

// RegisterNotifier registers a notifier for a specific system.
func (nm *NotificationManager) RegisterNotifier(system NotificationSystem, notifier Notifier) {
	nm.notifiers[system] = notifier
}

// RegisterNotification dynamically adds a notification template to the registry.
func (nm *NotificationManager) RegisterNotification(noticeType NoticeType, system NotificationSystem, template NoticeTemplate) error {
	// Validate input
	if noticeType == "" || system == "" || template.BodyPath == "" || template.Subject == "" {
		return fmt.Errorf("invalid input: notification type, system, subject, and bodyPath cannot be empty")
	}

	// Check if the notification type exists in the registry
	if _, exists := nm.notificationRegistry[noticeType]; !exists {
		nm.notificationRegistry[noticeType] = make(map[NotificationSystem]NoticeTemplate)
	}

	// Add or update the template for the system under the given notification type
	nm.notificationRegistry[noticeType][system] = template
	return nil
}

// Send sends a notification to all systems registered for the specified notification type.
func (nm *NotificationManager) Send(noticeType NoticeType, notification NotificationData) error {
	// Check if the notification type exists in the registry
	systemTemplates, exists := nm.notificationRegistry[noticeType]
	if !exists {
		return fmt.Errorf("no templates registered for notification type: %s", noticeType)
	}

	var lastError error
	notifierFound := false

	// Iterate through all systems registered for the notification type
	for system, template := range systemTemplates {
		// Get the notifier for the current system
		notifier, notifierExists := nm.notifiers[system]
		if !notifierExists {
			lastError = fmt.Errorf("no notifier registered for system: %s", system)
			continue
		}

		notifierFound = true

		// Render the template (if applicable)
		fmt.Printf("Using template for system %s: %s\n", system, template.Subject)

		// Send the notification using the notifier
		err := notifier.Send(noticeType, notification, template)
		if err != nil {
			// Log the error and store it as the last error (if any)
			fmt.Printf("Error sending notification via %s: %v\n", system, err)
			lastError = err
		}
	}

	if !notifierFound {
		return lastError
	}

	// Return the last error if any occurred during the process
	return lastError
}
