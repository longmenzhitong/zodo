package zodo

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
)

func readTodoLines(storageType string) []string {
	switch storageType {
	case StorageTypeFile:
		return ReadLinesFromPath(todoPath)
	case StorageTypeRedis:
		var lines []string
		cmd := Redis().Get(todoRedisKey)
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
		panic(&InvalidConfigError{
			Message: fmt.Sprintf("storage.type: %s", storageType),
		})
	}
}

func writeTodoLines(lines []string, storageType string) {
	switch storageType {
	case StorageTypeFile:
		RewriteLinesToPath(todoPath, lines)
		return
	case StorageTypeRedis:
		linesJson, err := json.Marshal(lines)
		if err != nil {
			panic(err)
		}
		Redis().Set(todoRedisKey, linesJson, 0)
		if Config.Storage.Redis.Localize {
			writeTodoLines(lines, StorageTypeFile)
		}
		return
	default:
		panic(&InvalidConfigError{
			Message: fmt.Sprintf("storage.type: %s", storageType),
		})
	}
}
