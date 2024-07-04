package mailer

import (
	"backend/app/common/utils"
	"fmt"
	"net/smtp"
)

type Config struct {
	Auth   smtp.Auth
	From   string
	Server string
}

var mailerInstance *Config

func InitializeMailer() {
	fmt.Println("===== Initialize Mailer =====")

	email := utils.GetEnv("MAILER_EMAIL")
	password := utils.GetEnv("MAILER_PASSWORD")

	mailerInstance = &Config{
		Auth:   smtp.PlainAuth("", email, password, "smtp.gmail.com"),
		From:   email,
		Server: "smtp.gmail.com:587",
	}
}

func GetMailerInstance() *Config {
	return mailerInstance
}
