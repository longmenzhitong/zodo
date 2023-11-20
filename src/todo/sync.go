package todo

import (
	"encoding/json"
	"fmt"
	"strconv"
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
		// 推送数据
		dataJson, err := json.Marshal(data)
		if err != nil {
			panic(err)
		}
		zodo.Redis().Set(redisKeyData, dataJson, 0)

		// 推送ID
		zodo.Redis().Set(redisKeyId, id, 0)

		// 推送时间戳
		zodo.Redis().Set(redisKeyTimestamp, ts, 0)
		return nil
	default:
		return &zodo.InvalidConfigError{
			Message: fmt.Sprintf("sync.type: %s", zodo.Config.Sync.Type),
		}
	}
}

func Pull() error {
	var data []string
	var id int

	switch zodo.Config.Sync.Type {
	case zodo.SyncTypeRedis:
		// 拉取数据
		cmd := zodo.Redis().Get(redisKeyData)
		dataJson, err := cmd.Result()
		if err != nil {
			return err
		}
		err = json.Unmarshal([]byte(dataJson), &data)
		if err != nil {
			return err
		}

		// 拉取ID
		cmd = zodo.Redis().Get(redisKeyId)
		idStr, err := cmd.Result()
		if err != nil {
			return err
		}
		id, err = strconv.Atoi(idStr)
		if err != nil {
			return err
		}
	default:
		return &zodo.InvalidConfigError{
			Message: fmt.Sprintf("sync.type: %s", zodo.Config.Sync.Type),
		}
	}

	if len(data) > 0 && id > 0 {
		Cache.set(data)
		Cache.save()
		zodo.Id.SetNext(id)
	}

	return nil
}
