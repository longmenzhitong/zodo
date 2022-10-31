package main

import (
	"fmt"
	"zodo/internal/conf"
	"zodo/internal/orders"
	"zodo/internal/param"
	"zodo/internal/stdin"
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
