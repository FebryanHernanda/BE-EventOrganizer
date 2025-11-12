package config

import (
	"os"
	"strconv"
)

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

func LoadSMTPConfig() SMTPConfig {
	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))

	return SMTPConfig{
		Host:     os.Getenv("SMTP_HOST"),
		Port:     port,
		Username: os.Getenv("SMTP_USER"),
		Password: os.Getenv("SMTP_PASS"),
		From:     os.Getenv("SMTP_USER"),
	}
}
