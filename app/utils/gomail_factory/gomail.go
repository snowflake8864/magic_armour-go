package gomail_factory

import (
	"gopkg.in/gomail.v2"
	"goskeleton/app/global/variable"
	"strings"
)

func CreateMailerFactory() *Mailer {
	host := variable.ConfigYml.GetString("EmailConfig.Host")
	port := variable.ConfigYml.GetInt("EmailConfig.Port")
	username := variable.ConfigYml.GetString("EmailConfig.Username")
	password := variable.ConfigYml.GetString("EmailConfig.Password")
	sendName := variable.ConfigYml.GetString("EmailConfig.SendName")
	encryption := variable.ConfigYml.GetString("EmailConfig.Encryption")

	d := gomail.NewDialer(host, port, username, password)
	if strings.ToLower(encryption) == "ssl" {
		d.SSL = true
	}

	return &Mailer{
		Dialer:    d,
		FromEmail: username,
		FromName:  sendName,
	}

}

// 结构体
type Mailer struct {
	Dialer    *gomail.Dialer
	FromEmail string
	FromName  string
}

func (k *Mailer) SendHtmlMsg(to, subject, body string) error {
	message := gomail.NewMessage()
	message.SetAddressHeader("From", k.FromEmail, k.FromName)
	message.SetHeader("To", to)
	message.SetHeader("Subject", subject)
	message.SetBody("text/html", body)

	err := k.Dialer.DialAndSend(message)
	if err != nil {
		variable.ZapLog.Sugar().Debug(err.Error())
		return err
	}
	return nil
}
