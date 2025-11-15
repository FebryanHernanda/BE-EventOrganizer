package provider

import (
	"github.com/FebryanHernanda/BE-EventOrganizer/config"
	"gopkg.in/gomail.v2"
)

type SMTPProvider struct {
	cfg config.SMTPConfig
}

func NewSMTPProvider(cfg config.SMTPConfig) *SMTPProvider {
	return &SMTPProvider{cfg}
}

func (p *SMTPProvider) Send(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", p.cfg.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(p.cfg.Host, p.cfg.Port, p.cfg.Username, p.cfg.Password)

	return d.DialAndSend(m)
}
