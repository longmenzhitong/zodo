package main

import (
	"zodo/cmd"
	zodo "zodo/src"
	"zodo/src/todo"
)

func main() {
	zodo.Config.Init()
	todo.Cache.Init()
	zodo.Id.Init()

	cmd.Execute()
}
