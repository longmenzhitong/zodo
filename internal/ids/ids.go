package ids

import (
	"strconv"
	"zodo/internal/cst"
	"zodo/internal/file"
)

const (
	fileName = "id"
)

var (
	path string
)

func init() {
	path = file.Dir + cst.PathSep + fileName
	file.EnsureExist(path)
}

func Get() int {
	var id int
	lines := file.ReadLinesFromPath(path)
	if len(lines) == 0 {
		id = 1
	} else {
		n, err := strconv.Atoi(lines[0])
		if err != nil {
			panic(err)
		}
		id = n
	}

	file.RewriteLinesToPath(path, []string{strconv.Itoa(id + 1)})

	return id
}
