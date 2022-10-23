package ids

import (
	"strconv"
	"zodo/internal/cst"
	"zodo/internal/files"
)

const (
	fileName = "id"
)

var (
	path string
)

func init() {
	path = files.Dir + cst.PathSep + fileName
	files.EnsureExist(path)
}

func Get() int {
	var id int
	lines := files.ReadLinesFromPath(path)
	if len(lines) == 0 {
		id = 1
	} else {
		n, err := strconv.Atoi(lines[0])
		if err != nil {
			panic(err)
		}
		id = n
	}

	files.RewriteLinesToPath(path, []string{strconv.Itoa(id + 1)})

	return id
}
