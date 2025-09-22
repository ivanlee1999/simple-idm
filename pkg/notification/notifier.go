package notification

type Notifier interface {
	Send(noticeType NoticeType, data NotificationData, template NoticeTemplate) error
}

type NotificationData struct {
	To      string            // Recipient identifier (e.g., email address, phone number, Slack channel)
	Subject string            // Optional: Subject for notifications like email
	Body    string            // The content or message to send
	Data    map[string]string // Additional metadata (e.g., for SMS sender ID, Slack channel ID, email CC/BCC as comma-separated values)
}
