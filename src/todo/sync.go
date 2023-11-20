package todo

import (
	"encoding/json"
	"fmt"
	"strconv"
	zodo "zodo/src"
)

const (
	redisKeyData = "zodo:data"
	redisKeyId   = "zodo:id"
)

const (
	s3ObjectKeyData = "zodo_data"
	s3ObjectKeyId   = "zodo_id"
)

func Push() error {
	data := zodo.ReadLinesFromPath(path)
	id := zodo.Id.GetNext()

	switch zodo.Config.Sync.Type {
	case zodo.SyncTypeRedis:
		// 推送数据
		dataJson, err := json.Marshal(data)
		if err != nil {
			return err
		}
		zodo.Redis().Set(redisKeyData, dataJson, 0)

		// 推送ID
		zodo.Redis().Set(redisKeyId, id, 0)
		return nil
	case zodo.SyncTypeS3:
		// 推送数据
		err := zodo.PushToS3(path, s3ObjectKeyData)
		if err != nil {
			return err
		}

		// 推送ID
		err = zodo.PushToS3(zodo.Id.Path, s3ObjectKeyId)
		if err != nil {
			return err
		}

		return nil
	default:
		return invalidSyncTypeConfigError()
	}
}

func Pull() error {
	switch zodo.Config.Sync.Type {
	case zodo.SyncTypeRedis:
		var data []string
		var id int
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
		if len(data) > 0 && id > 0 {
			// 备份并写入数据
			Cache.set(data)
			Cache.save()
			// 备份并写入ID
			zodo.Id.SetNext(id)
		}
		return nil
	case zodo.SyncTypeS3:
		// 备份数据
		Cache.save()
		// 备份ID
		zodo.Id.Backup()
		// 拉取并写入数据
		zodo.PullFromS3(path, s3ObjectKeyData)
		// 拉取并写入ID
		zodo.PullFromS3(zodo.Id.Path, s3ObjectKeyId)
		return nil
	default:
		return invalidSyncTypeConfigError()
	}
}

func invalidSyncTypeConfigError() error {
	return &zodo.InvalidConfigError{
		Message: fmt.Sprintf("sync.type: %s, expect 'redis' or 's3'", zodo.Config.Sync.Type),
	}
}
