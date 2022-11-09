package main

import (
	"github.com/jessevdk/go-flags"
	"zodo/internal/command"
)

func main() {
	var opt command.Option
	_, err := flags.Parse(&opt)
	if err != nil {
		panic(err)
	}
}
