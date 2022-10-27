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
	err := c.AddFunc(conf.Data.Reminder.DailyReport.Cron, func() {
		err := backup.Pull()
		if err != nil {
			fmt.Println(err.Error())
		}
		err = todos.DailyReport()
		if err != nil {
			fmt.Println(err.Error())
		}
	})
	if err != nil {
		panic(err)
	}
	c.Start()
}
