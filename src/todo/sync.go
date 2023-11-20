package todo

import (
	"encoding/json"
	"fmt"
	"time"
	zodo "zodo/src"
)

const (
	redisKeyData      = "zodo:data"
	redisKeyId        = "zodo:id"
	redisKeyTimestamp = "zodo:timestamp"
)

func Push() error {
	data := zodo.ReadLinesFromPath(path)
	id := zodo.Id.GetNext()
	ts := time.Now().Unix()

	switch zodo.Config.Sync.Type {
	case zodo.SyncTypeRedis:
		// 同步数据
		dataJson, err := json.Marshal(data)
		if err != nil {
			panic(err)
		}
		zodo.Redis().Set(redisKeyData, dataJson, 0)

		// 同步ID
		zodo.Redis().Set(redisKeyId, id, 0)

		// 同步时间戳
		zodo.Redis().Set(redisKeyTimestamp, ts, 0)
		return nil
	default:
		return &zodo.InvalidConfigError{
			Message: fmt.Sprintf("sync.type: %s", zodo.Config.Sync.Type),
		}
	}
}
