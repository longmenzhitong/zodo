package param

import (
	"flag"
	"strings"
)

var (
	Interactive bool
	Server      bool
)

var (
	Input string
)

func init() {
	flag.BoolVar(&Interactive, "i", false, "enter interactive mode")
	flag.BoolVar(&Server, "s", false, "enter server mode")
	flag.Parse()

	if len(flag.Args()) > 0 {
		Input = strings.Join(flag.Args(), " ")
	}
}
