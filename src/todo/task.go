package todo

import (
	"fmt"
	zodo "zodo/src"

	"github.com/robfig/cron"
)

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
