package reminder

import (
	"fmt"
	"github.com/jordan-wright/email"
	"github.com/robfig/cron"
	"net/smtp"
	"strings"
	"zodo/internal/conf"
	"zodo/internal/todo"
)

func StartDailyReport() {
	c := cron.New()
	err := c.AddFunc(conf.All.Reminder.DailyReport.Cron, func() {
		sendEmail("Daily Report", todo.DailyReport())
	})
	if err != nil {
		panic(err)
	}
	c.Start()
}

func sendEmail(title, text string) {
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
