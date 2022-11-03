package emails

import (
	"github.com/go-gomail/gomail"
	"zodo/internal/conf"
)

func Send(title, text string) error {
	m := gomail.NewMessage()
	m.SetAddressHeader("From", conf.Data.Email.From, "ZODO")
	m.SetHeader("To", conf.Data.Email.To...)
	m.SetHeader("Subject", title)
	m.SetBody("text/plain", text)

	d := gomail.NewDialer(conf.Data.Email.Server, conf.Data.Email.Port, conf.Data.Email.From, conf.Data.Email.Auth)
	return d.DialAndSend(m)
}
