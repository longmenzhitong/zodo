package todo

import (
	zodo "zodo/src"
)

func readTodoLines() []string {
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
	return
}
