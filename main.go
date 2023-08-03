package main

import (
	"fmt"
	"os"
	"zodo/cmd"
	zodo "zodo/src"
	"zodo/src/todo"
)

func main() {
	zodo.Config.Init()
	todo.Cache.Init()
	zodo.Id.Init()

	if err := cmd.RootCmd.Execute(); err != nil {
		if err != cmd.SilentErr {
			fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(1)
	}
}
