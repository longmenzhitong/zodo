package todos

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"strconv"
	"strings"
	"zodo/internal/files"
	"zodo/internal/times"
)

type todo struct {
	Id         int
	Content    string
	Status     string
	Deadline   string
	Remark     string
	CreateTime string
	ParentId   int
	Children   map[int]bool
}

func (t *todo) getStatus() string {
	if t.hasChildren() {
		return ""
	}
	switch t.Status {
	case statusPending:
		return color.HiMagentaString(t.Status)
	case statusProcessing:
		return color.HiCyanString(t.Status)
	case statusDone:
		return color.HiBlueString(t.Status)
	default:
		return t.Status
	}
}

func (t *todo) getDeadLine() string {
	if t.Deadline == "" || t.hasChildren() {
		return ""
	}

	if t.Status == statusDone {
		return times.Simplify(t.Deadline)
	}

	nd, wd := calcRemainDays(t.Deadline)
	ddl := fmt.Sprintf("%s (%dnd/%dwd)", t.Deadline, nd, wd)
	ddl = times.Simplify(ddl)

	if t.Status == statusPending || t.Status == statusProcessing {
		if wd == 0 && nd == 0 {
			ddl = color.RedString(ddl)
		} else if wd == 1 || nd == 1 {
			ddl = color.HiYellowString(ddl)
		} else {
			ddl = color.GreenString(ddl)
		}
	}
	return ddl
}

func (t *todo) getCreateTime() string {
	return times.Simplify(t.CreateTime)
}

func (t *todo) getParentId() string {
	return strconv.Itoa(t.ParentId)
}

func (t *todo) getChildren() string {
	if !t.hasChildren() {
		return ""
	}
	childIds := make([]string, 0)
	for id := range t.Children {
		childIds = append(childIds, strconv.Itoa(id))
	}
	return strings.Join(childIds, ",")
}

func (t *todo) hasChildren() bool {
	return t.Children != nil && len(t.Children) > 0
}

type data struct {
	List []*todo
	Map  map[int]*todo
}

func (d *data) Refresh() {
	d.List = make([]*todo, 0)
	d.Map = make(map[int]*todo, 0)
	for _, line := range files.ReadLinesFromPath(path) {
		var td todo
		err := json.Unmarshal([]byte(line), &td)
		if err != nil {
			panic(err)
		}
		d.List = append(d.List, &td)
		d.Map[td.Id] = &td
	}
}

func (d *data) save() {
	lines := make([]string, 0)
	for _, td := range d.List {
		js, err := json.Marshal(td)
		if err != nil {
			panic(err)
		}
		lines = append(lines, string(js))
	}
	files.RewriteLinesToPath(path, lines)
}

func (d *data) add(td todo) {
	d.List = append(d.List, &td)
	d.Map[td.Id] = &td
	d.save()
}

func (d *data) delete(id int) {
	newList := make([]*todo, 0)
	for _, td := range d.List {
		if td.Id != id {
			newList = append(newList, td)
		}
	}
	d.List = newList

	toDelete := d.Map[id]
	if toDelete != nil {
		parent := d.Map[toDelete.ParentId]
		if parent != nil {
			delete(parent.Children, id)
		}
	}

	delete(d.Map, id)

	d.save()
}

func (d *data) clear() int {
	n := 0
	newList := make([]*todo, 0)
	for _, td := range d.List {
		if td.ParentId != 0 {
			if _, ok := d.Map[td.ParentId]; !ok {
				delete(d.Map, td.Id)
				n++
				continue
			}
		}
		newList = append(newList, td)
	}
	d.List = newList
	d.save()
	return n
}

const fileName = "todo"

var path string

var Data *data

func init() {
	path = files.GetPath(fileName)
	files.EnsureExist(path)

	Data = &data{}
	Data.Refresh()
}
