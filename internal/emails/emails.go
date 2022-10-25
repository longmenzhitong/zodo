package emails

import (
	"fmt"
	"github.com/jordan-wright/email"
	"net/smtp"
	"strings"
	"zodo/internal/conf"
)

func Send(title, text string) {
	em := email.NewEmail()
	from := conf.All.Reminder.Email.From
	em.From = fmt.Sprintf("ZODO <%s>", from)
	em.To = conf.All.Reminder.Email.To
	em.Subject = title
	em.Text = []byte(text)

	addr := conf.All.Reminder.Email.Server
	auth := conf.All.Reminder.Email.Auth
	err := em.Send(addr, smtp.PlainAuth("", from, auth, strings.Split(addr, ":")[0]))
	if err != nil {
		panic(err)
	}
}
