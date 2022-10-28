package param

import (
	"flag"
	"strings"
)

var (
	Interactive bool
	Server      bool
	All         bool
	Delete      bool
	ParentId    int
)

var Input string

func init() {
	flag.BoolVar(&Interactive, "i", false, "enter interactive mode")
	flag.BoolVar(&Server, "s", false, "enter server mode")
	flag.BoolVar(&All, "a", false, "all")
	flag.BoolVar(&Delete, "d", false, "delete: -d [id1] [id2] [...]")
	flag.IntVar(&ParentId, "p", -1, "add and set parent: -p [parentId] [content]")
}

func Parse() {
	flag.Parse()
	if len(flag.Args()) > 0 {
		Input = strings.Join(flag.Args(), " ")
	}
}
