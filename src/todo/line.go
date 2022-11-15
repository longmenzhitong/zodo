package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"zodo/src"
)

func readTodoLines(storageType string) []string {
	switch storageType {
	case zodo.StorageTypeFile:
		return zodo.ReadLinesFromPath(todoPath)
	case zodo.StorageTypeRedis:
		var lines []string
		cmd := zodo.Redis().Get(todoRedisKey)
		linesJson, err := cmd.Result()
		if errors.Is(err, redis.Nil) {
			return lines
		}
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal([]byte(linesJson), &lines)
		if err != nil {
			panic(err)
		}
		return lines
	default:
		panic(&zodo.InvalidConfigError{
			Message: fmt.Sprintf("storage.type: %s", storageType),
		})
	}
}

func writeTodoLines(lines []string, storageType string) {
	switch storageType {
	case zodo.StorageTypeFile:
		zodo.RewriteLinesToPath(todoPath, lines)
		return
	case zodo.StorageTypeRedis:
		linesJson, err := json.Marshal(lines)
		if err != nil {
			panic(err)
		}
		zodo.Redis().Set(todoRedisKey, linesJson, 0)
		if zodo.Config.Storage.Redis.Localize {
			writeTodoLines(lines, zodo.StorageTypeFile)
		}
		return
	default:
		panic(&zodo.InvalidConfigError{
			Message: fmt.Sprintf("storage.type: %s", storageType),
		})
	}
}
