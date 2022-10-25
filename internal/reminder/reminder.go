package reminder

import (
	"github.com/robfig/cron"
	"zodo/internal/conf"
	"zodo/internal/todo"
)

func StartDailyReport() {
	c := cron.New()
	err := c.AddFunc(conf.All.Reminder.DailyReport.Cron, func() {
		todo.DailyReport()
	})
	if err != nil {
		panic(err)
	}
	c.Start()
}
