package ids

import (
	"strconv"
	"zodo/internal/conf"
	"zodo/internal/files"
	"zodo/internal/redish"
)

const fileName = "id"

const key = "zd:id"

var path string

func init() {
	path = files.GetPath(fileName)
	files.EnsureExist(path)
}

func Get() int {
	if conf.IsFileStorage() {
		return getFromFile()
	} else {
		return getFromRedis()
	}
}

func getFromFile() int {
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

func getFromRedis() int {
	cmd := redish.Client.Incr(key)
	id, err := cmd.Result()
	if err != nil {
		panic(err)
	}
	return int(id)
}
