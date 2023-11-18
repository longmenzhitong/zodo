package todo

import (
	zodo "zodo/src"
)

func readTodoLines() []string {
	// TODO 自动同步
	return zodo.ReadLinesFromPath(path)
	// switch storageType {
	// case zodo.StorageTypeFile:
	// 	return zodo.ReadLinesFromPath(path)
	// case zodo.StorageTypeRedis:
	// 	var lines []string
	// 	cmd := zodo.Redis().Get(redisKey)
	// 	linesJson, err := cmd.Result()
	// 	if errors.Is(err, redis.Nil) {
	// 		return lines
	// 	}
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	err = json.Unmarshal([]byte(linesJson), &lines)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	return lines
	// default:
	// 	panic(&zodo.InvalidConfigError{
	// 		Message: fmt.Sprintf("storage.type: %s", storageType),
	// 	})
	// }
}

func writeTodoLines(lines []string) {
	zodo.RewriteLinesToPath(path, lines)
	// TODO 自动同步
	return
	// switch storageType {
	// case zodo.StorageTypeFile:
	// 	zodo.RewriteLinesToPath(path, lines)
	// 	return
	// case zodo.StorageTypeRedis:
	// 	linesJson, err := json.Marshal(lines)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	zodo.Redis().Set(redisKey, linesJson, 0)
	// 	if zodo.Config.Storage.Redis.Localize {
	// 		writeTodoLines(lines, zodo.StorageTypeFile)
	// 	}
	// 	return
	// default:
	// 	panic(&zodo.InvalidConfigError{
	// 		Message: fmt.Sprintf("storage.type: %s", storageType),
	// 	})
	// }
}
