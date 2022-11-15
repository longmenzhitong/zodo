package todo

import (
	"fmt"
	"github.com/robfig/cron"
	"zodo/src"
)

func StartDailyReport() {
	c := cron.New()
	err := c.AddFunc(zodo.Config.DailyReport.Cron, func() {
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
	err := c.AddFunc(zodo.Config.Reminder.Cron, func() {
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
