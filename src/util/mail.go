package util

import (
	"net/smtp"
	"strings"
)

func SendMail(to string, subject string, body string) (err error) {
	var (
		auth   smtp.Auth
		sendTo []string
		msg    []byte
	)
	auth = smtp.PlainAuth("", "linfangrong@oneniceapp.com", "a417498669", "smtp.exmail.qq.com")
	sendTo = strings.Split(to, ";")
	msg = []byte("To: " + to + "\r\n" +
		"From: auto@oneniceapp.com\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: text/plain; charset=UTF-8\r\n\r\n" +
		body)
	err = smtp.SendMail("smtp.exmail.qq.com:25", auth, "linfangrong@oneniceapp.com", sendTo, msg)
	return
}
