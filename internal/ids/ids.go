package ids

import (
	"strconv"
	"zodo/internal/conf"
	"zodo/internal/errs"
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

func Get(storageType string) int {
	if conf.IsFileStorage(storageType) {
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
		return id
	}
	if conf.IsRedisStorage(storageType) {
		cmd := redish.Client.Get(key)
		idStr, err := cmd.Result()
		if err != nil {
			panic(err)
		}
		id, err := strconv.Atoi(idStr)
		if err != nil {
			panic(err)
		}
		return id
	}
	panic(&errs.InvalidConfigError{
		Config:  "storage.type",
		Message: storageType,
	})
}

func Set(id int, storageType string) {
	if conf.IsFileStorage(storageType) {
		files.RewriteLinesToPath(path, []string{strconv.Itoa(id)})
		return
	}
	if conf.IsRedisStorage(storageType) {
		redish.Client.Set(key, id, 0)
		return
	}
	panic(&errs.InvalidConfigError{
		Config:  "storage.type",
		Message: storageType,
	})
}

func GetAndSet(storageType string) int {
	id := Get(storageType)
	Set(id+1, storageType)
	return id
}
