package todos

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/go-redis/redis"
	"strconv"
	"strings"
	"zodo/internal/conf"
	"zodo/internal/files"
	"zodo/internal/redish"
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
	var lines []string
	if conf.IsFileStorage() {
		lines = files.ReadLinesFromPath(path)
	} else {
		cmd := redish.Client.Get(key)
		linesJson, err := cmd.Result()
		if errors.Is(err, redis.Nil) {
			return
		}
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal([]byte(linesJson), &lines)
		if err != nil {
			panic(err)
		}
	}
	for _, line := range lines {
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
	if conf.IsFileStorage() {
		files.RewriteLinesToPath(path, lines)
	} else {
		linesJson, err := json.Marshal(lines)
		if err != nil {
			panic(err)
		}
		redish.Client.Set(key, linesJson, 0)
	}
}

func (d *data) add(td todo) {
	d.List = append(d.List, &td)
	d.Map[td.Id] = &td
	d.save()
}

func (d *data) delete(id int) {
	toDelete := d.Map[id]
	if toDelete == nil {
		return
	}

	newList := make([]*todo, 0)
	for _, td := range d.List {
		if td.Id != id {
			newList = append(newList, td)
		}
	}
	d.List = newList

	parent := d.Map[toDelete.ParentId]
	if parent != nil {
		delete(parent.Children, id)
	}

	if toDelete.hasChildren() {
		for childId := range toDelete.Children {
			d.delete(childId)
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

const key = "zd:todo"

var path string

var Data *data

func init() {
	path = files.GetPath(fileName)
	files.EnsureExist(path)

	Data = &data{}
	Data.Refresh()
}
