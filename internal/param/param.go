package param

import (
	"flag"
	"strings"
)

var (
	Server   bool
	All      bool
	ParentId int
	Deadline string
)

var Input string

func init() {
	flag.BoolVar(&Server, "s", false, "enter server mode")
	flag.BoolVar(&All, "a", false, "all")
	flag.IntVar(&ParentId, "p", 0, "add and set parent: -p [parentId] [content]")
	flag.StringVar(&Deadline, "d", "", "add and set deadline: -d [deadline] [content]")
}

func Parse() {
	flag.Parse()
	if len(flag.Args()) > 0 {
		Input = strings.Join(flag.Args(), " ")
	}
}
