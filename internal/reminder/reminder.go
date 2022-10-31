package reminder

import (
	"fmt"
	"github.com/robfig/cron"
	"zodo/internal/conf"
	"zodo/internal/todos"
)

func StartDailyReport() {
	// fixme
	c := cron.New()
	err := c.AddFunc(conf.Data.Reminder.DailyReport.Cron, func() {
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
