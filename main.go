package main

import (
	"github.com/jessevdk/go-flags"
	zodo "zodo/src"
)

func main() {
	zodo.InitConfig()
	zodo.InitCache()

	var opt zodo.Option
	_, _ = flags.Parse(&opt)
}
