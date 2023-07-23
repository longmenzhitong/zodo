package zodo

import (
	"fmt"
	"strconv"
)

const (
	idFileName = "id"
	idRedisKey = "zd:id"
)

var (
	idPath     string
	backupPath string
)

func init() {
	idPath = Path(idFileName)
	backupPath = idPath + ".backup"
}

func Id(storageType string) int {
	id := GetId(storageType)
	SetId(id+1, storageType)
	return id
}

func GetId(storageType string) int {
	switch storageType {
	case StorageTypeFile:
		return getIdFromPath(idPath)
	case StorageTypeRedis:
		return getIdFromRedis()
	default:
		panic(&InvalidConfigError{
			Message: fmt.Sprintf("storage.type: %s", storageType),
		})
	}
}

func getIdFromPath(path string) int {
	var id int
	lines := ReadLinesFromPath(path)
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

func getIdFromRedis() int {
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
}

func SetId(id int, storageType string) {
	BackupId(storageType)

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

func BackupId(storageType string) {
	curId := GetId(storageType)
	RewriteLinesToPath(backupPath, []string{strconv.Itoa(curId)})
}

func RollbackId(storageType string) {
	SetId(getIdFromPath(backupPath), storageType)
}
