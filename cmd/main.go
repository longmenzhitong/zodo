package main

import (
	"fmt"
	"zodo/internal/backup"
	"zodo/internal/conf"
	"zodo/internal/orders"
	"zodo/internal/param"
	"zodo/internal/reminder"
	"zodo/internal/stdin"
)

func init() {
	param.Parse()

	err := backup.CheckPull()
	if err != nil {
		fmt.Println(err.Error())
	}
}

func main() {
	if param.Server {
		if conf.All.Reminder.DailyReport.Enabled {
			reminder.StartDailyReport()
		}
		select {}
	}

	if param.Interactive {
		fmt.Println("================")
		fmt.Println("Welcome to ZODO!")
		fmt.Println("================")
		for true {
			fmt.Print("$ ")
			err := orders.Handle(stdin.ReadString())
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}

	err := orders.Handle(param.Input)
	if err != nil {
		fmt.Println(err.Error())
	}
}
