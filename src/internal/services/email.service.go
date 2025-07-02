package services

import "github.com/LombardiDaniel/goliath/src/internal/models"

// EmailService defines the interface for email-related operations.
// It provides methods for sending various types of emails, such as
// email confirmations, account creation notifications, organization
// invites, password reset emails, and payment acceptance notifications.
type EmailService interface {
	// SendEmailConfirmation sends an email confirmation to a user.
	SendEmailConfirmation(email string, name string, otp string) error

	// SendAccountCreated notifies a user that their account has been created.
	SendAccountCreated(email string, name string) error

	// SendOrganizationInvite sends an invitation email to join an organization.
	SendOrganizationInvite(email string, name string, otp string, orgName string) error

	// SendPasswordReset sends a password reset email to a user.
	SendPasswordReset(email string, name string, otp string) error

	// SendPaymentAccepted notifies a user that their payment has been accepted.
	SendPaymentAccepted(email string, name string, payment models.Payment) error
}

type EmailServiceMock struct{}

func (s *EmailServiceMock) SendEmailConfirmation(email string, name string, otp string) error {
	return nil
}
func (s *EmailServiceMock) SendAccountCreated(email string, name string) error {
	return nil
}
func (s *EmailServiceMock) SendOrganizationInvite(email string, name string, otp string, orgName string) error {
	return nil
}
func (s *EmailServiceMock) SendPasswordReset(email string, name string, otp string) error {
	return nil
}
func (s *EmailServiceMock) SendPaymentAccepted(email string, name string, payment models.Payment) error {
	return nil
}
