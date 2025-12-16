package email_service

import (
	"context"
	"os"
)

type Email struct {
	UserId  string
	Subject string
	Body    string
	EventID string
}

type EmailSender interface {
	Send(ctx context.Context, msg Email) error
}

type EmailService struct{}

func (e *EmailService) Send(_ context.Context, _ Email) error {
	return nil
}

type EmailServiceIntegrationTests struct{}

func (e *EmailServiceIntegrationTests) Send(_ context.Context, msg Email) error {
	path := os.Getenv("INTEGRATION_TESTS_LOG_PATH")
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	_, err = f.WriteString(msg.EventID + "\n")
	if err != nil {
		return err
	}

	return nil
}
