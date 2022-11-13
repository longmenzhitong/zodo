package zodo

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/mozillazg/go-pinyin"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	statusPending    = "Pending"
	statusProcessing = "Processing"
	statusDone       = "Done"
	statusHiding     = "Hiding"
)

var statusPriority = map[string]int{
	statusHiding:     -1,
	statusDone:       0,
	statusPending:    1,
	statusProcessing: 2,
}

type todo struct {
	Id           int
	Content      string
	Status       string
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
}

func (t *todo) getStatus(colorful bool) string {
	if colorful {
		switch t.Status {
		case statusPending:
			return color.HiMagentaString(t.Status)
		case statusProcessing:
			return color.HiCyanString(t.Status)
		case statusDone:
			return color.HiBlueString(t.Status)
		case statusHiding:
			return color.HiBlackString(t.Status)
		}
	}
	return t.Status
}

func (t *todo) getRemainDays() (natureDays int, workDays int) {
	ddlTime, err := time.Parse(LayoutYearMonthDay, t.Deadline)
	if err != nil {
		panic(err)
	}
	return CalcBetweenDays(time.Now(), ddlTime)
}

func (t *todo) getDeadLineAndRemain(colorful bool) (ddl string, remain string) {
	if t.Deadline == "" {
		return "", ""
	}

	if t.Status == statusDone {
		return SimplifyTime(t.Deadline), ""
	}

	ddlTime, err := time.Parse(LayoutYearMonthDay, t.Deadline)
	if err != nil {
		panic(err)
	}

	ddl = fmt.Sprintf("%s(%s)", t.Deadline, ddlTime.Weekday().String()[:3])
	ddl = SimplifyTime(ddl)

	nd, wd := t.getRemainDays()
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

func (t *todo) getRemindTime() string {
	return SimplifyTime(t.RemindTime)
}

func (t *todo) getDoneTime() string {
	return SimplifyTime(t.DoneTime)
}

func (t *todo) getCreateTime() string {
	return SimplifyTime(t.CreateTime)
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
	if t.Status == statusHiding {
		return false
	}
	if t.Status == statusDone && !Config.Todo.ShowDone {
		return false
	}
	return true
}

const (
	todoFileName = "todo"
	todoRedisKey = "zd:todo"
)

var (
	todoPath   string
	backupPath string
)

func init() {
	todoPath = Path(todoFileName)
	backupPath = todoPath + ".backup"
}

func hitTodo(td *todo, keyword string) bool {
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
		m := cc._map()
		for childId := range td.Children {
			if hitTodo(m[childId], keyword) {
				return true
			}
		}
	}

	return false
}

func walkTodo(td *todo, tds *[]todo, level int, status []string, allStatus bool) {
	if td == nil {
		return
	}

	if !hitStatus(td, status, allStatus) {
		return
	}

	td.Level = level
	*tds = append(*tds, *td)

	if td.Children == nil || len(td.Children) == 0 {
		return
	}

	m := cc._map()
	childList := make([]*todo, 0)
	for childId := range td.Children {
		child := m[childId]
		if child == nil {
			fmt.Println(&NotFoundError{
				Target:  "child",
				Message: fmt.Sprintf("parentId: %d, childId: %d", td.Id, childId),
			})
		} else {
			childList = append(childList, child)
		}
	}
	childList = sortTodo(childList)

	for _, child := range childList {
		walkTodo(child, tds, level+1, status, allStatus)
	}
}

func hitStatus(td *todo, status []string, allStatus bool) bool {
	if td == nil {
		return false
	}
	if allStatus {
		return true
	}
	if len(status) == 0 {
		return td.isVisible()
	}
	if td.hasChildren() {
		for childId := range td.Children {
			if hitStatus(cc._map()[childId], status, allStatus) {
				return true
			}
		}
		return false
	} else {
		hit := false
		for _, s := range status {
			if strings.HasPrefix(strings.ToLower(td.Status), strings.ToLower(s)) {
				hit = true
				break
			}
		}
		return hit
	}
}

func sortTodo(tds []*todo) []*todo {
	sort.SliceStable(tds, func(i, j int) bool {
		a := tds[i]
		b := tds[j]

		if a.Status != b.Status {
			return statusPriority[a.Status] > statusPriority[b.Status]
		}

		if a.Status == statusDone && b.Status == statusDone {
			return a.Id < b.Id
		}

		if a.Deadline != b.Deadline {
			if a.Deadline != "" && b.Deadline != "" {
				ta, err := time.Parse(LayoutYearMonthDay, a.Deadline)
				if err != nil {
					panic(err)
				}
				tb, err := time.Parse(LayoutYearMonthDay, b.Deadline)
				if err != nil {
					panic(err)
				}
				return ta.Unix() < tb.Unix()
			} else {
				return a.Deadline != ""
			}
		}

		return a.Id < b.Id
	})
	return tds
}
