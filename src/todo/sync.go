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
	zodo.PrintStartMsg("Pushing data...\n")
	err := pushData()
	if err != nil {
		return err
	}

	zodo.PrintStartMsg("Pushing ID...\n")
	err = pushID()
	if err != nil {
		return err
	}

	zodo.PrintDoneMsg("Done.\n")
	return nil
}

func pushData() error {
	switch zodo.Config.Sync.Type {
	case zodo.SyncTypeRedis:
		data := zodo.ReadLinesFromPath(path)
		dataJson, err := json.Marshal(data)
		if err != nil {
			return err
		}
		zodo.Redis().Set(redisKeyData, dataJson, 0)
		return nil
	case zodo.SyncTypeS3:
		return zodo.PushToS3(path, s3ObjectKeyData)
	default:
		return invalidSyncTypeConfigError()
	}
}

func pushID() error {
	switch zodo.Config.Sync.Type {
	case zodo.SyncTypeRedis:
		id := zodo.Id.GetNext()
		zodo.Redis().Set(redisKeyId, id, 0)
		return nil
	case zodo.SyncTypeS3:
		return zodo.PushToS3(zodo.Id.Path, s3ObjectKeyId)
	default:
		return invalidSyncTypeConfigError()
	}
}

func Pull() error {
	zodo.PrintStartMsg("Pulling data...\n")
	err := pullData()
	if err != nil {
		return err
	}

	zodo.PrintStartMsg("Pulling ID...\n")
	err = pullID()
	if err != nil {
		return err
	}

	zodo.PrintDoneMsg("Done.\n")
	return nil
}

func pullData() error {
	switch zodo.Config.Sync.Type {
	case zodo.SyncTypeRedis:
		// 拉取
		var data []string
		cmd := zodo.Redis().Get(redisKeyData)
		dataJson, err := cmd.Result()
		if err != nil {
			return err
		}
		err = json.Unmarshal([]byte(dataJson), &data)
		if err != nil {
			return err
		}
		if len(data) > 0 {
			// 备份并写入
			Cache.set(data)
			Cache.save()
		}
		return nil
	case zodo.SyncTypeS3:
		// 备份
		Cache.save()
		// 拉取并写入
		return zodo.PullFromS3(path, s3ObjectKeyData)
	default:
		return invalidSyncTypeConfigError()
	}
}

func pullID() error {
	switch zodo.Config.Sync.Type {
	case zodo.SyncTypeRedis:
		// 拉取
		var id int
		cmd := zodo.Redis().Get(redisKeyId)
		idStr, err := cmd.Result()
		if err != nil {
			return err
		}
		id, err = strconv.Atoi(idStr)
		if err != nil {
			return err
		}
		if id > 0 {
			// 备份并写入
			zodo.Id.SetNext(id)
		}
		return nil
	case zodo.SyncTypeS3:
		// 备份
		zodo.Id.Backup()
		// 拉取并写入
		return zodo.PullFromS3(zodo.Id.Path, s3ObjectKeyId)
	default:
		return invalidSyncTypeConfigError()
	}
}

func invalidSyncTypeConfigError() error {
	return &zodo.InvalidConfigError{
		Message: fmt.Sprintf("sync.type: %s, expect 'redis' or 's3'", zodo.Config.Sync.Type),
	}
}
