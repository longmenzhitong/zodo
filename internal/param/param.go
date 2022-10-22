package param

import (
	"flag"
	"strings"
)

var (
	// 模式参数
	Interactive bool
	Server      bool
)

var (
	All   bool
	Input string
)

func init() {
	flag.BoolVar(&Interactive, "i", false, "enter interactive mode")
	flag.BoolVar(&Server, "s", false, "enter server mode")
	flag.BoolVar(&All, "a", false, "all")
	flag.Parse()

	if len(flag.Args()) > 0 {
		Input = strings.Join(flag.Args(), " ")
	}
}
