package zodo

import (
	"fmt"
	"github.com/robfig/cron"
)

func StartDailyReport() {
	c := cron.New()
	err := c.AddFunc(Config.DailyReport.Cron, func() {
		err := Report()
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
	err := c.AddFunc(Config.Reminder.Cron, func() {
		err := Remind()
		if err != nil {
			fmt.Println(err.Error())
		}
	})
	if err != nil {
		panic(err)
	}
	c.Start()
}
