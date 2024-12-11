package services

import (
	"bytes"
	"log/slog"
	"net/url"
	"path/filepath"
	"text/template"

	"github.com/LombardiDaniel/gopherbase/common"
	"github.com/LombardiDaniel/gopherbase/models"
	"github.com/resendlabs/resend-go"
)

type EmailServiceResendImpl struct {
	resendClient *resend.Client

	emailConfirmationTemplate  *template.Template
	accountCreationTemplate    *template.Template
	organizationInviteTemplate *template.Template
	passwordResetTemplate      *template.Template
	paymentAcceptedTemplate    *template.Template

	usersConfirmUrl  string
	acceptInviteUrl  string
	passwordResetUrl string
}

func NewEmailServiceResendImpl(resendApiKey string, templatesDir string) EmailService {
	usersConfirmUrl, err := url.JoinPath(common.API_HOST_URL, "/v1/users/confirm")
	if err != nil {
		panic(err)
	}

	acceptInviteUrl, err := url.JoinPath(common.API_HOST_URL, "/v1/organizations/accept-invite")
	if err != nil {
		panic(err)
	}

	passwordResetUrl, err := url.JoinPath(common.API_HOST_URL, "/v1/users/set-password-reset-cookie")
	if err != nil {
		panic(err)
	}

	return &EmailServiceResendImpl{
		resendClient:               resend.NewClient(resendApiKey),
		emailConfirmationTemplate:  common.LoadHTMLTemplate(filepath.Join(templatesDir, "email-confirmation.html")),
		accountCreationTemplate:    common.LoadHTMLTemplate(filepath.Join(templatesDir, "account-created.html")),
		organizationInviteTemplate: common.LoadHTMLTemplate(filepath.Join(templatesDir, "organization-invite.html")),
		passwordResetTemplate:      common.LoadHTMLTemplate(filepath.Join(templatesDir, "password-reset.html")),
		paymentAcceptedTemplate:    common.LoadHTMLTemplate(filepath.Join(templatesDir, "payment-accepted.html")),
		usersConfirmUrl:            usersConfirmUrl,
		acceptInviteUrl:            acceptInviteUrl,
		passwordResetUrl:           passwordResetUrl,
	}
}

type htmlConfirmationVars struct {
	ProjectName string
	FirstName   string
	OtpUrl      string
}

func (s *EmailServiceResendImpl) SendEmailConfirmation(email string, name string, otp string) error {
	body := new(bytes.Buffer)
	err := s.emailConfirmationTemplate.Execute(body, htmlConfirmationVars{
		ProjectName: common.PROJECT_NAME,
		FirstName:   name,
		OtpUrl:      s.usersConfirmUrl + "?otp=" + otp,
	})
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	params := &resend.SendEmailRequest{
		From:    common.NOREPLY_EMAIL,
		To:      []string{email},
		Subject: "Confirm Your Account!",
		Html:    body.String(),
	}

	_, err = s.resendClient.Emails.Send(params)

	return err
}

type htmlAccountCreatedVars struct {
	FirstName string
}

func (s *EmailServiceResendImpl) SendAccountCreated(email string, name string) error {
	body := new(bytes.Buffer)
	err := s.accountCreationTemplate.Execute(body, htmlAccountCreatedVars{
		FirstName: name,
	})
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	params := &resend.SendEmailRequest{
		From:    common.NOREPLY_EMAIL,
		To:      []string{email},
		Subject: "Account Created!",
		Html:    body.String(),
	}

	_, err = s.resendClient.Emails.Send(params)

	return err
}

type htmlOrgInviteVars struct {
	ProjectName      string
	OrganizationName string
	FirstName        string
	OtpUrl           string
}

func (s *EmailServiceResendImpl) SendOrganizationInvite(email string, name string, otp string, orgName string) error {
	body := new(bytes.Buffer)
	err := s.organizationInviteTemplate.Execute(body, htmlOrgInviteVars{
		ProjectName:      common.PROJECT_NAME,
		OrganizationName: orgName,
		FirstName:        name,
		OtpUrl:           s.acceptInviteUrl + "?otp=" + otp,
	})
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	params := &resend.SendEmailRequest{
		From:    common.NOREPLY_EMAIL,
		To:      []string{email},
		Subject: "Organization Invite",
		Html:    body.String(),
	}

	_, err = s.resendClient.Emails.Send(params)

	return err
}

type htmlPwResetVars struct {
	ProjectName string
	FirstName   string
	OtpUrl      string
}

func (s *EmailServiceResendImpl) SendPasswordReset(email string, name string, otp string) error {
	body := new(bytes.Buffer)
	err := s.passwordResetTemplate.Execute(body, htmlPwResetVars{
		ProjectName: common.PROJECT_NAME,
		FirstName:   name,
		OtpUrl:      s.passwordResetUrl + "?otp=" + otp,
	})
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	params := &resend.SendEmailRequest{
		From:    common.NOREPLY_EMAIL,
		To:      []string{email},
		Subject: "Password Reset",
		Html:    body.String(),
	}

	_, err = s.resendClient.Emails.Send(params)

	return err
}

type htmlPaymentAccepted struct {
	FirstName string
	PaymentId string
}

func (s *EmailServiceResendImpl) SendPaymentAccepted(email string, name string, payment models.Payment) error {
	body := new(bytes.Buffer)
	err := s.paymentAcceptedTemplate.Execute(body, htmlPaymentAccepted{
		FirstName: name,
		PaymentId: payment.PaymentId,
	})
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	params := &resend.SendEmailRequest{
		From:    common.NOREPLY_EMAIL,
		To:      []string{email},
		Subject: "Payment Accepted",
		Html:    body.String(),
	}

	_, err = s.resendClient.Emails.Send(params)

	return err
}
