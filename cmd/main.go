package main

import (
	"fmt"
	"zodo/internal/conf"
	"zodo/internal/orders"
	"zodo/internal/param"
	"zodo/internal/task"
)

func init() {
	param.Parse()
}

func main() {
	if param.Server {
		if conf.Data.DailyReport.Enabled {
			task.StartDailyReport()
		}
		if conf.Data.Reminder.Enabled {
			task.StartReminder()
		}
		select {}
	}

	err := orders.Handle(param.Input)
	if err != nil {
		fmt.Println(err.Error())
	}
}
