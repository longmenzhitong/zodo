package todos

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/go-redis/redis"
	"github.com/mozillazg/go-pinyin"
	"sort"
	"strconv"
	"strings"
	"time"
	"zodo/internal/conf"
	"zodo/internal/cst"
	"zodo/internal/errs"
	"zodo/internal/files"
	"zodo/internal/ids"
	"zodo/internal/param"
	"zodo/internal/redish"
	"zodo/internal/times"
)

type todo struct {
	Id           int
	Content      string
	Status       string
	Deadline     string
	Remark       string
	RemindTime   string
	RemindStatus string
	LoopType     string
	CreateTime   string
	ParentId     int
	Children     map[int]bool
	Level        int
}

func (t *todo) getStatus(colorful bool) string {
	if t.hasChildren() {
		return ""
	}
	if colorful {
		switch t.Status {
		case statusPending:
			return color.HiMagentaString(t.Status)
		case statusProcessing:
			return color.HiCyanString(t.Status)
		case statusDone:
			return color.HiBlueString(t.Status)
		}
	}
	return t.Status
}

func (t *todo) getDeadLineAndRemain(colorful bool) (ddl string, remain string) {
	if t.Deadline == "" || t.hasChildren() {
		return "", ""
	}

	if t.Status == statusDone {
		return times.Simplify(t.Deadline), ""
	}

	ddlTime, err := time.Parse(cst.LayoutYearMonthDay, t.Deadline)
	if err != nil {
		panic(err)
	}

	ddl = fmt.Sprintf("%s(%s)", t.Deadline, ddlTime.Weekday().String())
	ddl = times.Simplify(ddl)

	nd, wd := calcRemainDays(t.Deadline)
	remain = fmt.Sprintf("%dnd/%dwd", nd, wd)

	if colorful && (t.Status == statusPending || t.Status == statusProcessing) {
		if wd == 0 && nd == 0 {
			ddl = color.RedString(ddl)
			remain = color.RedString(remain)
		} else if wd == 1 || nd == 1 {
			ddl = color.HiYellowString(ddl)
			remain = color.HiYellowString(remain)
		} else {
			ddl = color.GreenString(ddl)
			remain = color.GreenString(remain)
		}
	}
	return
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

const fileName = "todo"
const key = "zd:todo"

var path string
var backupPath string
var _list []*todo

func init() {
	path = files.GetPath(fileName)
	backupPath = path + ".backup"
	load()
}

func load() {
	_list = make([]*todo, 0)
	for _, line := range readLines(conf.Data.Storage.Type) {
		var td todo
		err := json.Unmarshal([]byte(line), &td)
		if err != nil {
			panic(err)
		}
		_list = append(_list, &td)
	}
}

func readLines(storageType string) []string {
	if conf.IsFileStorage(storageType) {
		return files.ReadLinesFromPath(path)
	}
	if conf.IsRedisStorage(storageType) {
		var lines []string
		cmd := redish.Client().Get(key)
		linesJson, err := cmd.Result()
		if errors.Is(err, redis.Nil) {
			return lines
		}
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal([]byte(linesJson), &lines)
		if err != nil {
			panic(err)
		}
		return lines
	}
	panic(&errs.InvalidConfigError{
		Config:  "storage.type",
		Message: storageType,
	})
}

func save() {
	backup()
	lines := make([]string, 0)
	for _, td := range _list {
		js, err := json.Marshal(td)
		if err != nil {
			panic(err)
		}
		lines = append(lines, string(js))
	}
	writeLines(lines, conf.Data.Storage.Type)
}

func writeLines(lines []string, storageType string) {
	if conf.IsFileStorage(storageType) {
		files.RewriteLinesToPath(path, lines)
		return
	}
	if conf.IsRedisStorage(storageType) {
		linesJson, err := json.Marshal(lines)
		if err != nil {
			panic(err)
		}
		redish.Client().Set(key, linesJson, 0)

		if conf.Data.Storage.Redis.Localize {
			writeLines(lines, conf.StorageTypeFile)
		}

		return
	}
	panic(&errs.InvalidConfigError{
		Config:  "storage.type",
		Message: storageType,
	})
}

func backup() {
	lines := readLines(conf.Data.Storage.Type)
	files.RewriteLinesToPath(backupPath, lines)
}

func rollback() {
	lines := files.ReadLinesFromPath(backupPath)
	writeLines(lines, conf.Data.Storage.Type)
}

func transfer() {
	if conf.IsFileStorage() {
		lines := readLines(conf.StorageTypeRedis)
		writeLines(lines, conf.StorageTypeFile)

		id := ids.Get(conf.StorageTypeRedis)
		ids.Set(id+1, conf.StorageTypeFile)
		return
	}
	if conf.IsRedisStorage() {
		lines := readLines(conf.StorageTypeFile)
		writeLines(lines, conf.StorageTypeRedis)

		id := ids.Get(conf.StorageTypeFile)
		ids.Set(id+1, conf.StorageTypeRedis)
		return
	}
	panic(&errs.InvalidConfigError{
		Config:  "storage.type",
		Message: conf.Data.Storage.Type,
	})
}

func list(keyword string) []todo {
	tds := make([]todo, 0)
	for _, td := range _list {
		if td.ParentId == 0 && hit(td, keyword) {
			walk(td, &tds, 0)
		}
	}
	return tds
}

func hit(td *todo, keyword string) bool {
	if td == nil {
		return false
	}
	if !param.All && td.Status == statusDone {
		return false
	}
	if keyword == "" {
		return true
	}

	content := strings.ToLower(td.Content)
	keyword = strings.ToLower(keyword)
	if strings.Contains(content, keyword) {
		return true
	}

	pa := pinyin.NewArgs()
	pyArrays := pinyin.Pinyin(content, pa)
	var pyStr string
	for _, pyArray := range pyArrays {
		for _, py := range pyArray {
			pyStr += py
		}
	}
	if strings.Contains(pyStr, keyword) {
		return true
	}

	if td.hasChildren() {
		for childId := range td.Children {
			if hit(_map()[childId], keyword) {
				return true
			}
		}
	}

	return false
}

func walk(td *todo, tds *[]todo, level int) {
	if td == nil {
		return
	}
	if !param.All && td.Status == statusDone {
		return
	}

	td.Level = level
	*tds = append(*tds, *td)

	if td.Children == nil || len(td.Children) == 0 {
		return
	}

	childList := make([]*todo, 0)
	for childId := range td.Children {
		child := _map()[childId]
		if child == nil {
			fmt.Println(&errs.NotFoundError{
				Target:  "child",
				Message: fmt.Sprintf("parentId: %d, childId: %d", td.Id, childId),
			})
		} else {
			childList = append(childList, child)
		}
	}
	childList = _sort(childList)

	for _, child := range childList {
		walk(child, tds, level+1)
	}
}

func _sort(tds []*todo) []*todo {
	sort.Slice(tds, func(i, j int) bool {
		a := tds[i]
		b := tds[j]
		if a.Deadline != "" && b.Deadline != "" {
			ta, err := time.Parse(cst.LayoutYearMonthDay, a.Deadline)
			if err != nil {
				panic(err)
			}
			tb, err := time.Parse(cst.LayoutYearMonthDay, b.Deadline)
			if err != nil {
				panic(err)
			}
			return ta.Unix() < tb.Unix()
		}
		if a.Deadline == "" && b.Deadline == "" {
			if a.Status != b.Status {
				return a.Status == statusProcessing
			}
			return a.Id < b.Id
		}
		return a.Deadline != ""
	})
	return tds
}

func padding(level int, unit string) string {
	var res string
	for i := 0; i < level; i++ {
		res += unit
	}
	return res
}

func _map() map[int]*todo {
	res := make(map[int]*todo, 0)
	for _, td := range _list {
		res[td.Id] = td
	}
	return res
}

func add(td todo) {
	_list = append(_list, &td)
}

func _delete(id int) {
	m := _map()
	toDelete := m[id]
	if toDelete == nil {
		return
	}

	newList := make([]*todo, 0)
	for _, td := range _list {
		if td.Id != id {
			newList = append(newList, td)
		}
	}
	_list = newList

	parent := m[toDelete.ParentId]
	if parent != nil {
		delete(parent.Children, id)
	}

	if toDelete.hasChildren() {
		for childId := range toDelete.Children {
			_delete(childId)
		}
	}
}
