package mailer

import (
	"backend/app/config/mailer"
	"fmt"
	"log"
	"net/smtp"
)

type MailInfo struct {
	EmailTarget []string
	Subject     string
	Body        string
}

func SendMail(mailInfo MailInfo) {
	go func() {
		mailerInstance := mailer.GetMailerInstance()
		message := []byte(fmt.Sprintf("Subject: %s\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s\r\n", mailInfo.Subject, mailInfo.Body))

		log.Println("Sending email to:", mailInfo.EmailTarget)
		err := smtp.SendMail(mailerInstance.Server, mailerInstance.Auth, mailerInstance.From, mailInfo.EmailTarget, message)
		if err != nil {
			log.Println("Error sending email to:", mailInfo.EmailTarget, err)
			return
		}
	}()
}
