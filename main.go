package main

import (
	zodo "zodo/src"
	"zodo/src/cmd"
	"zodo/src/todo"

	"github.com/jessevdk/go-flags"
)

func main() {
	zodo.Config.Init()
	todo.Cache.Init()

	var opt cmd.Option
	flags.Parse(&opt)
}
