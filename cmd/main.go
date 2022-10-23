package main

import (
	"fmt"
	"zodo/internal/orders"
	"zodo/internal/param"
	"zodo/internal/stdin"
)

func init() {
	param.Parse()
}

func main() {
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
	} else {
		err := orders.Handle(param.Input)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}
