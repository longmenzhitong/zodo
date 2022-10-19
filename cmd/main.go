package main

import (
	"fmt"
	"zodo/internal/orders"
	"zodo/internal/stdin"
	"zodo/internal/todo"
)

func main() {
	defer todo.Save()

	fmt.Println("================")
	fmt.Println("Welcome to ZODO!")
	fmt.Println("================")

	for true {
		fmt.Print("$ ")
		input := stdin.ReadString()

		// todo help

		if orders.IsExit(input) {
			return
		}

		if orders.IsList(input) {
			todo.List()
			continue
		}

		// todo detail

		if orders.IsAdd(input) {
			content, err := orders.ParseAdd(input)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			todo.Add(content)
			continue
		}

		if orders.IsModify(input) {
			id, content, err := orders.ParseModify(input)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			todo.Modify(id, content, "")
			continue
		}

		if orders.IsPending(input) {
			id, err := orders.ParsePending(input)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			todo.Pending(id)
			continue
		}

		if orders.IsDone(input) {
			id, err := orders.ParseDone(input)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			todo.Done(id)
			continue
		}

		if orders.IsAbandon(input) {
			id, err := orders.ParseAbandon(input)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			todo.Abandon(id)
			continue
		}

		if orders.IsDelete(input) {
			id, err := orders.ParseDelete(input)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			todo.Delete(id)
			continue
		}

		// todo hint
	}
}
