package main

import (
	"fmt"
	"os"
	"zodo/internal/orders"
	"zodo/internal/param"
	"zodo/internal/stdin"
	"zodo/internal/todo"
)

func main() {
	defer todo.Save()

	if param.Interactive {
		fmt.Println("================")
		fmt.Println("Welcome to ZODO!")
		fmt.Println("================")
		for true {
			fmt.Print("$ ")
			err := handle(stdin.ReadString())
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	} else {
		err := handle(param.Input)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

}

func handle(input string) error {
	if orders.IsExit(input) {
		todo.Save()
		os.Exit(0)
	}

	if orders.IsList(input) {
		todo.List()
		return nil
	}

	if orders.IsAdd(input) {
		content, err := orders.ParseAdd(input)
		if err != nil {
			return err
		}
		todo.Add(content)
		return nil
	}

	if orders.IsModify(input) {
		id, content, err := orders.ParseModify(input)
		if err != nil {
			return err
		}
		todo.Modify(id, content)
		return nil
	}

	if orders.IsDeadline(input) {
		id, deadline, err := orders.ParseDeadline(input)
		if err != nil {
			return err
		}
		todo.Deadline(id, deadline)
		return nil
	}

	if orders.IsPending(input) {
		id, err := orders.ParsePending(input)
		if err != nil {
			return err
		}
		todo.Pending(id)
		return nil
	}

	if orders.IsDone(input) {
		id, err := orders.ParseDone(input)
		if err != nil {
			return err
		}
		todo.Done(id)
		return nil
	}

	if orders.IsAbandon(input) {
		id, err := orders.ParseAbandon(input)
		if err != nil {
			return err
		}
		todo.Abandon(id)
		return nil
	}

	if orders.IsDelete(input) {
		id, err := orders.ParseDelete(input)
		if err != nil {
			return err
		}
		todo.Delete(id)
		return nil
	}

	// todo help
	// todo detail
	// todo hint

	todo.List()
	return nil
}
