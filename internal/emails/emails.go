package emails

import (
	"fmt"
	"github.com/jordan-wright/email"
	"net/smtp"
	"strings"
	"zodo/internal/conf"
)

func Send(title, text string) error {
	em := email.NewEmail()
	from := conf.Data.Reminder.Email.From
	em.From = fmt.Sprintf("ZODO <%s>", from)
	em.To = conf.Data.Reminder.Email.To
	em.Subject = title
	em.Text = []byte(text)

	addr := conf.Data.Reminder.Email.Server
	auth := conf.Data.Reminder.Email.Auth
	return em.Send(addr, smtp.PlainAuth("", from, auth, strings.Split(addr, ":")[0]))
}
