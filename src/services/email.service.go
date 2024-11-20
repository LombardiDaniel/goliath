package services

type EmailService interface {
	SendEmailConfirmation(email string, name string, otp string) error
	SendAccountCreated(email string, name string) error
	SendOrganizationInvite(email string, name string, otp string, orgName string) error
	SendPasswordReset(email string, name string, otp string) error
}
