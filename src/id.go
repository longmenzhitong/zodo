package zodo

import (
	"strconv"
)

const (
	idFileName = "id"
)

var Id id

type id struct {
	Path       string
	backupPath string

	next int
}

func (i *id) Init() {
	i.Path = Path(idFileName)
	i.backupPath = i.Path + ".backup"

	i.next = i.readNext()
}

func (i *id) GetSetNext() int {
	next := i.GetNext()
	i.SetNext(next + 1)
	return next
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
	return getIdFromPath(i.Path)
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
	RewriteLinesToPath(i.Path, []string{strconv.Itoa(id)})
	return
}
