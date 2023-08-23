package todo

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
	zodo "zodo/src"

	"github.com/mozillazg/go-pinyin"
)

type Status string

const (
	StatusPending    Status = "Pending"
	StatusProcessing Status = "Processing"
	StatusDone       Status = "Done"
	StatusHiding     Status = "Hiding"
)

type todo struct {
	Id           int
	Content      string
	Status       Status
	Deadline     string
	Remark       string
	RemindTime   string
	RemindStatus string
	LoopType     string
	DoneTime     string
	CreateTime   string
	ParentId     int
	Children     map[int]bool
	Level        int
	Priority     int
}

func (t *todo) getStatus(colorful bool) string {
	status := string(t.Status)
	if colorful {
		switch t.Status {
		case StatusPending:
			return zodo.ColoredString(zodo.Config.Todo.Color.Status.Pending, status)
		case StatusProcessing:
			return zodo.ColoredString(zodo.Config.Todo.Color.Status.Processing, status)
		case StatusDone:
			return zodo.ColoredString(zodo.Config.Todo.Color.Status.Done, status)
		case StatusHiding:
			return zodo.ColoredString(zodo.Config.Todo.Color.Status.Hiding, status)
		}
	}
	return string(t.Status)
}

func (t *todo) getRemainDays() (natureDays int, workDays int) {
	ddlTime, err := time.Parse(zodo.LayoutDate, t.Deadline)
	if err != nil {
		panic(err)
	}
	return zodo.CalcBetweenDays(time.Now(), ddlTime)
}

func (t *todo) getDeadLineAndRemain(colorful bool) (ddl string, remain string) {
	if t.Deadline == "" {
		return "", ""
	}

	if t.Status == StatusDone {
		return zodo.SimplifyTime(t.Deadline), ""
	}

	ddlTime, err := time.Parse(zodo.LayoutDate, t.Deadline)
	if err != nil {
		panic(err)
	}

	ddl = fmt.Sprintf("%s(%s)", t.Deadline, ddlTime.Weekday().String()[:3])
	ddl = zodo.SimplifyTime(ddl)

	nd, wd := t.getRemainDays()
	remain = fmt.Sprintf("%dnd/%dwd", nd, wd)

	if colorful && (t.Status == StatusPending || t.Status == StatusProcessing) {
		if wd <= 0 && nd <= 0 {
			ddl = zodo.ColoredString(zodo.Config.Todo.Color.Deadline.Overdue, ddl)
			remain = zodo.ColoredString(zodo.Config.Todo.Color.Deadline.Overdue, remain)
		} else if wd == 1 || nd == 1 {
			ddl = zodo.ColoredString(zodo.Config.Todo.Color.Deadline.Nervous, ddl)
			remain = zodo.ColoredString(zodo.Config.Todo.Color.Deadline.Nervous, remain)
		} else {
			ddl = zodo.ColoredString(zodo.Config.Todo.Color.Deadline.Normal, ddl)
			remain = zodo.ColoredString(zodo.Config.Todo.Color.Deadline.Normal, remain)
		}
	}
	return
}

func (t *todo) getRemindTime() string {
	return zodo.SimplifyTime(t.RemindTime)
}

func (t *todo) getDoneTime() string {
	return zodo.SimplifyTime(t.DoneTime)
}

func (t *todo) getCreateTime() string {
	return zodo.SimplifyTime(t.CreateTime)
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

func (t *todo) isVisible() bool {
	if t.Status == StatusHiding {
		return false
	}
	if t.Status == StatusDone && !zodo.Config.Todo.ShowDone {
		return false
	}
	return true
}

const (
	fileName = "todo"
	redisKey = "zd:todo"
)

var (
	path       string
	backupPath string
)

func init() {
	path = zodo.Path(fileName)
	backupPath = path + ".backup"
}

func hitKeyword(td *todo, keyword string) bool {
	if td == nil {
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
			if hitKeyword(Cache.get(childId), keyword) {
				return true
			}
		}
	}

	return false
}

func walkTodo(td *todo, tds *[]todo, level int, allStatus bool) {
	if td == nil {
		return
	}

	if !hitStatus(td, allStatus) {
		return
	}

	td.Level = level
	*tds = append(*tds, *td)

	if td.Children == nil || len(td.Children) == 0 {
		return
	}

	childList := make([]*todo, 0)
	for childId := range td.Children {
		child := Cache.get(childId)
		if child == nil {
			fmt.Println(&zodo.NotFoundError{
				Target:  "child",
				Message: fmt.Sprintf("parentId: %d, childId: %d", td.Id, childId),
			})
		} else {
			childList = append(childList, child)
		}
	}
	childList = sortTodo(childList)

	for _, child := range childList {
		walkTodo(child, tds, level+1, allStatus)
	}
}

func hitStatus(td *todo, allStatus bool) bool {
	if td == nil {
		return false
	}
	if allStatus {
		return true
	}
	return td.isVisible()
}

func sortTodo(tds []*todo) []*todo {
	sort.SliceStable(tds, func(i, j int) bool {
		a := tds[i]
		b := tds[j]

		if a.Priority != b.Priority {
			return a.Priority > b.Priority
		}

		return a.Id < b.Id
	})
	return tds
}
