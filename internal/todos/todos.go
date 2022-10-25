package todos

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"strconv"
	"strings"
	"time"
	"zodo/internal/cst"
	"zodo/internal/emails"
	"zodo/internal/errs"
	"zodo/internal/files"
	"zodo/internal/ids"
	"zodo/internal/param"
	"zodo/internal/stdout"
	"zodo/internal/times"
)

const (
	typeChild  = "Child"
	typeParent = "Parent"
)

type todo struct {
	Id         int
	Content    string
	Status     string
	Deadline   string
	Remark     string
	CreateTime string
	Type       string
	Parent     map[int]bool
	Child      map[int]bool
}

func (td *todo) GetStatus() string {
	switch td.Status {
	case statusPending:
		return color.HiMagentaString(td.Status)
	case statusProcessing:
		return color.HiCyanString(td.Status)
	case statusDone:
		return color.HiBlueString(td.Status)
	default:
		return td.Status
	}
}

func (td *todo) GetDeadLine() string {
	if td.Deadline == "" {
		return ""
	}

	if td.Status == statusDone {
		return times.Simplify(td.Deadline)
	}

	nd, wd := calcRemainDays(td.Deadline)
	ddl := fmt.Sprintf("%s (%dnd/%dwd)", td.Deadline, nd, wd)
	ddl = times.Simplify(ddl)

	if td.Status == statusPending || td.Status == statusProcessing {
		if wd == 0 {
			ddl = color.RedString(ddl)
		} else if wd == 1 {
			ddl = color.HiYellowString(ddl)
		} else {
			ddl = color.GreenString(ddl)
		}
	}
	return ddl
}

func (td *todo) GetCreateTime() string {
	return times.Simplify(td.CreateTime)
}

func (td *todo) GetParent() string {
	if td.Parent == nil {
		return ""
	}
	parentIds := make([]string, 0)
	for id := range td.Parent {
		parentIds = append(parentIds, strconv.Itoa(id))
	}
	return strings.Join(parentIds, ",")
}

func (td *todo) GetChild() string {
	if td.Child == nil {
		return ""
	}
	childIds := make([]string, 0)
	for id := range td.Child {
		childIds = append(childIds, strconv.Itoa(id))
	}
	return strings.Join(childIds, ",")
}

const (
	fileName = "todo"
)

const (
	statusPending    = "Pending"
	statusProcessing = "Processing"
	statusDone       = "Done"
	statusDeleted    = "Deleted"
)

var (
	path  string
	tds   []*todo
	tdMap map[int]*todo
)

func init() {
	path = files.GetPath(fileName)
	files.EnsureExist(path)

	Load()
}

func Load() {
	tds = make([]*todo, 0)
	tdMap = make(map[int]*todo, 0)
	for _, line := range files.ReadLinesFromPath(path) {
		var td todo
		err := json.Unmarshal([]byte(line), &td)
		if err != nil {
			panic(err)
		}
		tds = append(tds, &td)
		tdMap[td.Id] = &td
	}
}

func List() {
	rows := make([]table.Row, 0)
	for _, td := range tds {
		if td.Status == statusDeleted {
			continue
		}

		if !param.All && td.Status == statusDone {
			continue
		}

		if td.Type == typeChild {
			continue
		}

		if td.Child == nil || len(td.Child) == 0 {
			rows = append(rows, table.Row{
				td.Id,
				td.Content,
				td.GetStatus(),
				td.GetDeadLine(),
			})
		} else {
			rows = append(rows, table.Row{
				td.Id,
				td.Content,
				"",
				"",
			})
			for childId := range td.Child {
				child := tdMap[childId]
				if child != nil {
					rows = append(rows, table.Row{
						child.Id,
						fmt.Sprintf("  * %s", child.Content),
						child.GetStatus(),
						child.GetDeadLine(),
					})
				}
			}
		}
	}
	stdout.PrintTable(table.Row{"Id", "Content", "Status", "Deadline"}, rows)
}

func Detail(id int) {
	td := tdMap[id]
	if td == nil {
		return
	}

	rows := make([]table.Row, 0)
	rows = append(rows, table.Row{"Id", td.Id})
	rows = append(rows, table.Row{"Content", td.Content})
	rows = append(rows, table.Row{"Status", td.GetStatus()})
	rows = append(rows, table.Row{"Deadline", td.GetDeadLine()})
	rows = append(rows, table.Row{"Remark", td.Remark})
	rows = append(rows, table.Row{"Create", td.GetCreateTime()})
	rows = append(rows, table.Row{"Type", td.Type})
	rows = append(rows, table.Row{"Parent", td.GetParent()})
	rows = append(rows, table.Row{"Child", td.GetChild()})
	stdout.PrintTable(table.Row{"Item", "Val"}, rows)
}

func DailyReport() {
	var text string
	for _, td := range tds {
		if td.Status == statusDeleted {
			continue
		}

		if !param.All && td.Status == statusDone {
			continue
		}

		text += fmt.Sprintf("%s [%s]\n", td.Content, td.Status)
		if td.Deadline != "" {
			text += fmt.Sprintf("%s is the deadline.\n", times.Simplify(td.Deadline))
		}
		if td.Remark != "" {
			text += fmt.Sprintf("%s\n", td.Remark)
		}
		text += fmt.Sprintf("Created on %s.\n", times.Simplify(td.CreateTime))
		text += fmt.Sprintf("====================\n")
	}
	emails.Send("Daily Report", text)
}

func Add(content string) {
	if content == "" {
		return
	}
	td := todo{
		Id:         ids.Get(),
		Content:    content,
		Status:     statusPending,
		CreateTime: time.Now().Format(cst.LayoutDateTime),
		Type:       typeParent,
	}
	tds = append(tds, &td)
	save()
}

func Modify(id int, content string) {
	if content == "" {
		return
	}
	td := tdMap[id]
	if td != nil {
		td.Content = content
	}
	save()
}

func SetDeadline(id int, deadline string) {
	td := tdMap[id]
	if td != nil {
		td.Deadline = deadline
	}
	save()
}

func SetRemark(id int, remark string) {
	td := tdMap[id]
	if td != nil {
		td.Remark = remark
	}
	save()
}

func SetChild(parentId int, childIds []int) error {
	parent := tdMap[parentId]
	if parent == nil {
		return &errs.NotFoundError{
			Target:  "parent",
			Message: fmt.Sprintf("parentId: %d", parentId),
		}
	}
	if parent.Child == nil {
		parent.Child = make(map[int]bool, 0)
	}
	for _, childId := range childIds {
		child := tdMap[childId]
		if child == nil {
			return &errs.NotFoundError{
				Target:  "child",
				Message: fmt.Sprintf("childId: %d", childId),
			}
		}
		child.Type = typeChild
		if child.Parent == nil {
			child.Parent = make(map[int]bool, 0)
		}
		child.Parent[parentId] = true
		parent.Child[childId] = true
	}
	save()
	return nil
}

// TODO 子任务状态的变更可能会影响父任务
func SetPending(id int) {
	modifyStatus(id, statusPending)
}

func SetProcessing(id int) {
	modifyStatus(id, statusProcessing)
}

func SetDone(id int) {
	modifyStatus(id, statusDone)
}

func SetDeleted(id int) {
	modifyStatus(id, statusDeleted)
}

func modifyStatus(id int, status string) {
	td := tdMap[id]
	if td != nil {
		td.Status = status
	}
	save()
}

func save() {
	lines := make([]string, 0)
	for _, td := range tds {
		js, err := json.Marshal(td)
		if err != nil {
			panic(err)
		}
		lines = append(lines, string(js))
	}
	files.RewriteLinesToPath(path, lines)
}

func calcRemainDays(deadline string) (natureDays int, workDays int) {
	ddlTime, err := time.Parse(cst.LayoutYearMonthDay, deadline)
	if err != nil {
		panic(err)
	}

	return times.CalcBetweenDays(time.Now(), ddlTime)
}
