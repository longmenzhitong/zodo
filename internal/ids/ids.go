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
	Path string
)

func init() {
	Path = file.Dir + cst.PathSep + fileName
	file.EnsureExist(Path)
}

func Get() int {
	var id int
	lines := file.ReadLinesFromPath(Path)
	if len(lines) == 0 {
		id = 1
	} else {
		n, err := strconv.Atoi(lines[0])
		if err != nil {
			panic(err)
		}
		id = n
	}

	file.RewriteLinesToPath(Path, []string{strconv.Itoa(id + 1)})

	return id
}
