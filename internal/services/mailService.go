package services

import (
	"mymod/internal/config"

	log "github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

var GlobalEmail config.EmailConfig

func SendEmail(reciever string) {

	msg := gomail.NewMessage()
	msg.SetHeader("From", GlobalEmail.Sender)
	msg.SetHeader("To", reciever)
	msg.SetHeader("Subject", "subject")
	msg.SetBody("text/html", "<b>This is the body of the mail</b>")

	n := gomail.NewDialer("smtp.gmail.com", 587, GlobalEmail.Sender, GlobalEmail.Password)

	if err := n.DialAndSend(msg); err != nil {
		panic(err)
	}

	log.Info("message sended")
}
