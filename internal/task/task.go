package task

import (
	"fmt"
	"github.com/robfig/cron"
	"zodo/internal/conf"
	"zodo/internal/todos"
)

func StartDailyReport() {
	c := cron.New()
	err := c.AddFunc(conf.Data.DailyReport.Cron, func() {
		err := todos.DailyReport()
		if err != nil {
			fmt.Println(err.Error())
		}
	})
	if err != nil {
		panic(err)
	}
	c.Start()
}

func StartReminder() {
	c := cron.New()
	err := c.AddFunc(conf.Data.Reminder.Cron, func() {
		err := todos.Remind()
		if err != nil {
			fmt.Println(err.Error())
		}
	})
	if err != nil {
		panic(err)
	}
	c.Start()
}
