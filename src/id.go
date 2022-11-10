package zodo

import (
	"fmt"
	"strconv"
)

const (
	idFileName = "id"
	idRedisKey = "zd:id"
)

var idPath string

func init() {
	idPath = Path(idFileName)
}

func GetId(storageType string) int {
	switch storageType {
	case StorageTypeFile:
		var id int
		lines := ReadLinesFromPath(idPath)
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
	case StorageTypeRedis:
		cmd := Redis().Get(idRedisKey)
		idStr, err := cmd.Result()
		if err != nil {
			panic(err)
		}
		id, err := strconv.Atoi(idStr)
		if err != nil {
			panic(err)
		}
		return id
	default:
		panic(&InvalidConfigError{
			Message: fmt.Sprintf("storage.type: %s", storageType),
		})
	}
}

func SetId(id int, storageType string) {
	switch storageType {
	case StorageTypeFile:
		RewriteLinesToPath(idPath, []string{strconv.Itoa(id)})
		return
	case StorageTypeRedis:
		Redis().Set(idRedisKey, id, 0)

		if Config.Storage.Redis.Localize {
			SetId(id, StorageTypeFile)
		}

		return
	default:
		panic(&InvalidConfigError{
			Message: fmt.Sprintf("storage.type: %s", storageType),
		})
	}
}

func Id(storageType string) int {
	id := GetId(storageType)
	SetId(id+1, storageType)
	return id
}
