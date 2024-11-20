package services

import (
	"bytes"
	"log/slog"
	"net/url"
	"path/filepath"
	"text/template"

	"github.com/LombardiDaniel/go-gin-template/common"
	"github.com/resendlabs/resend-go"
)

type EmailServiceResendImpl struct {
	resendClient                *resend.Client
	accountConfirmationTemplate *template.Template
	accountCreationTemplate     *template.Template
	organizationInviteTemplate  *template.Template
	passwordResetTemplate       *template.Template
}

func NewEmailServiceResendImpl(resendApiKey string, templatesDir string) EmailService {
	return &EmailServiceResendImpl{
		resendClient:                resend.NewClient(resendApiKey),
		accountConfirmationTemplate: common.LoadHTMLTemplate(filepath.Join(templatesDir, "account-confirmation.html")),
		accountCreationTemplate:     common.LoadHTMLTemplate(filepath.Join(templatesDir, "account-created.html")),
		organizationInviteTemplate:  common.LoadHTMLTemplate(filepath.Join(templatesDir, "organization-invite.html")),
		passwordResetTemplate:       common.LoadHTMLTemplate(filepath.Join(templatesDir, "password-reset.html")),
	}
}

type htmlConfirmationVars struct {
	ProjectName string
	FirstName   string
	OtpUrl      string
}

func (s *EmailServiceResendImpl) SendEmailConfirmation(email string, name string, otp string) error {

	body := new(bytes.Buffer)
	confirmUrl, err := url.JoinPath(common.API_HOST_URL, "/v1/users/confirm")
	if err != nil {
		return err
	}
	err = s.accountConfirmationTemplate.Execute(body, htmlConfirmationVars{
		ProjectName: common.PROJECT_NAME,
		FirstName:   name,
		OtpUrl:      confirmUrl + "?otp=" + otp,
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
	err := s.accountConfirmationTemplate.Execute(body, htmlAccountCreatedVars{
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

	acceptUrl, err := url.JoinPath(common.API_HOST_URL, "/v1/organizations/accept-invite")
	if err != nil {
		return err
	}
	body := new(bytes.Buffer)
	err = s.accountConfirmationTemplate.Execute(body, htmlOrgInviteVars{
		ProjectName:      common.PROJECT_NAME,
		OrganizationName: orgName,
		FirstName:        name,
		OtpUrl:           acceptUrl,
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

	resetUrl, err := url.JoinPath(common.API_HOST_URL, "/v1/users/set-password-reset-cookie")
	if err != nil {
		return err
	}
	body := new(bytes.Buffer)
	err = s.accountConfirmationTemplate.Execute(body, htmlPwResetVars{
		ProjectName: common.PROJECT_NAME,
		FirstName:   name,
		OtpUrl:      resetUrl,
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
