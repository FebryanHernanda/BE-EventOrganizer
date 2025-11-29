package services

import (
	"fmt"

	"github.com/FebryanHernanda/BE-EventOrganizer/internal/mailer/provider"
	"github.com/FebryanHernanda/BE-EventOrganizer/internal/mailer/templates"
	"github.com/sirupsen/logrus"
)

type MailerService struct {
	provider *provider.SMTPProvider
}

func NewMailerService(provider *provider.SMTPProvider) *MailerService {
	return &MailerService{provider}
}

func (s *MailerService) SendActivationEmail(name, email, registeredDate, token string) error {
	link := fmt.Sprintf("http://localhost:3000/auth/activation?token=%s", token)

	body := templates.ActivationTemplate(name, email, registeredDate, link)

	if err := s.provider.Send(email, "Activate Your Account", body); err != nil {
		logrus.WithFields(logrus.Fields{
			"email": email,
			"error": err,
		}).Error("Failed sending activation email")
		return err
	}

	logrus.WithField("email", email).Info("Activation email sent successfully")
	return nil
}
