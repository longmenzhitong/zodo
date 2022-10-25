package reminder

import (
	"fmt"
	"github.com/robfig/cron"
	"zodo/internal/backup"
	"zodo/internal/conf"
	"zodo/internal/todos"
)

func StartDailyReport() {
	c := cron.New()
	err := c.AddFunc(conf.All.Reminder.DailyReport.Cron, func() {
		err := backup.Pull()
		if err != nil {
			fmt.Println(err.Error())
		}
		todos.DailyReport()
	})
	if err != nil {
		panic(err)
	}
	c.Start()
}
