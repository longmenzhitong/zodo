package zodo

import (
	"strconv"
)

const (
	idFileName = "id"
	idRedisKey = "zd:id"
)

var Id id

type id struct {
	path       string
	backupPath string

	next int
}

func (i *id) Init() {
	i.path = Path(idFileName)
	i.backupPath = i.path + ".backup"

	i.next = i.readNext()
}

func (i *id) SetGetNext() int {
	i.SetNext(i.next + 1)
	return i.GetNext()
}

func (i *id) GetNext() int {
	return i.next
}

func (i *id) SetNext(id int) {
	i.Backup()
	i.next = id
	i.writeNext(id)
}

func (i *id) Backup() {
	RewriteLinesToPath(i.backupPath, []string{strconv.Itoa(i.GetNext())})
}

func (i *id) Rollback() {
	i.writeNext(getIdFromPath(i.backupPath))
}

func (i *id) readNext() int {
	// TODO 自动同步
	return getIdFromPath(i.path)
	// switch i.storageType {
	// case StorageTypeFile:
	// 	return getIdFromPath(i.path)
	// case StorageTypeRedis:
	// 	return getIdFromRedis()
	// default:
	// 	panic(&InvalidConfigError{
	// 		Message: fmt.Sprintf("storage.type: %s", i.storageType),
	// 	})
	// }
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

func (i *id) writeNext(id int) {
	RewriteLinesToPath(i.path, []string{strconv.Itoa(id)})
	// TODO 自动同步
	return
	// switch i.storageType {
	// case StorageTypeFile:
	// 	RewriteLinesToPath(i.path, []string{strconv.Itoa(id)})
	// 	return
	// case StorageTypeRedis:
	// 	Redis().Set(idRedisKey, id, 0)
	// 	if Config.Storage.Redis.Localize {
	// 		RewriteLinesToPath(i.path, []string{strconv.Itoa(id)})
	// 	}
	// 	return
	// default:
	// 	panic(&InvalidConfigError{
	// 		Message: fmt.Sprintf("storage.type: %s", i.storageType),
	// 	})
	// }
}
