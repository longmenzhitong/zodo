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
	flag.StringVar(&Deadline, "D", "", "add and set deadline: -D [deadline] [content]")
}

func Parse() {
	flag.Parse()
	if len(flag.Args()) > 0 {
		Input = strings.Join(flag.Args(), " ")
	}
}
