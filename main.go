package main

import (
	"github.com/jessevdk/go-flags"
	zodo "zodo/src"
	"zodo/src/cmd"
	"zodo/src/todo"
)

func main() {
	zodo.InitConfig()
	todo.InitCache()

	var opt cmd.Option
	_, _ = flags.Parse(&opt)
}
