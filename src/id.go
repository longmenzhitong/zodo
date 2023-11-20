package zodo

import (
	"strconv"
)

const (
	idFileName = "id"
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
	return getIdFromPath(i.path)
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

func (i *id) writeNext(id int) {
	RewriteLinesToPath(i.path, []string{strconv.Itoa(id)})
	return
}
