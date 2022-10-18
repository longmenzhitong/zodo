package main

import (
	"errors"
	"fmt"
	"zodo/internal/errs"
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

		if err := orders.IsExit(input); !errors.Is(err, &errs.WrongOrderError{}) {
			return
		}

		if err := orders.IsList(input); !errors.Is(err, &errs.WrongOrderError{}) {
			todo.List()
			continue
		}

		if content, err := orders.IsAdd(input); !errors.Is(err, &errs.WrongOrderError{}) {
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			todo.Add(content)
			continue
		}

		if id, content, err := orders.IsModify(input); !errors.Is(err, &errs.WrongOrderError{}) {
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			todo.Modify(id, content, "")
			continue
		}

		if id, err := orders.IsPending(input); !errors.Is(err, &errs.WrongOrderError{}) {
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			todo.Pending(id)
			continue
		}

		if id, err := orders.IsDone(input); !errors.Is(err, &errs.WrongOrderError{}) {
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			todo.Done(id)
			continue
		}

		if id, err := orders.IsAbandon(input); !errors.Is(err, &errs.WrongOrderError{}) {
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			todo.Abandon(id)
			continue
		}

		if id, err := orders.IsDelete(input); !errors.Is(err, &errs.WrongOrderError{}) {
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			todo.Delete(id)
			continue
		}

	}
}
