package services

type EmailService interface {
	SendAccountConfirmation(name string, email string, otp string) error
	SendAccountCreated(name string, email string) error
	SendOrganizationInvite(name string, email string, orgId string, orgName string) error
}
