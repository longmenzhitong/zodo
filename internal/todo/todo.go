package todo

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"time"
	"zodo/internal/cst"
	"zodo/internal/files"
	"zodo/internal/ids"
	"zodo/internal/param"
	"zodo/internal/stdout"
	"zodo/internal/times"
)

type todo struct {
	Id         int
	Content    string
	Status     string
	Deadline   string
	Remark     string
	CreateTime string
}

func (td *todo) GetStatus() string {
	switch td.Status {
	case statusPending:
		return color.HiMagentaString(td.Status)
	case statusProcessing:
		return color.HiBlueString(td.Status)
	case statusDone:
		return color.CyanString(td.Status)
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
	todos []todo
)

func init() {
	path = files.GetPath(fileName)
	files.EnsureExist(path)

	todos = make([]todo, 0)
	for _, line := range files.ReadLinesFromPath(path) {
		var td todo
		err := json.Unmarshal([]byte(line), &td)
		if err != nil {
			panic(err)
		}
		todos = append(todos, td)
	}
}

func Save() {
	lines := make([]string, 0)
	for _, td := range todos {
		js, err := json.Marshal(td)
		if err != nil {
			panic(err)
		}
		lines = append(lines, string(js))
	}
	files.RewriteLinesToPath(path, lines)
}

func List() {
	rows := make([]table.Row, 0)
	for _, td := range todos {
		if td.Status == statusDeleted {
			continue
		}

		if !param.All && td.Status == statusDone {
			continue
		}

		rows = append(rows, table.Row{
			td.Id,
			td.Content,
			td.GetStatus(),
			td.GetDeadLine(),
			td.GetCreateTime(),
		})
	}
	stdout.PrintTable(table.Row{"Id", "Content", "Status", "Deadline", "Create"}, rows)
}

func Detail(id int) {
	td := findById(id)
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
	stdout.PrintTable(table.Row{"Item", "Val"}, rows)
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
	}
	todos = append(todos, td)
}

func Modify(id int, content string) {
	if content == "" {
		return
	}
	td := findById(id)
	if td != nil {
		td.Content = content
	}
}

func Deadline(id int, deadline string) {
	td := findById(id)
	if td != nil {
		td.Deadline = deadline
	}
}

func Remark(id int, remark string) {
	td := findById(id)
	if td != nil {
		td.Remark = remark
	}
}

func Pending(id int) {
	modifyStatus(id, statusPending)
}

func Processing(id int) {
	modifyStatus(id, statusProcessing)
}

func Done(id int) {
	modifyStatus(id, statusDone)
}

func Delete(id int) {
	modifyStatus(id, statusDeleted)
}

func modifyStatus(id int, status string) {
	td := findById(id)
	if td != nil {
		td.Status = status
	}
}

func findById(id int) *todo {
	for i := range todos {
		if todos[i].Id == id {
			return &todos[i]
		}
	}
	return nil
}

func calcRemainDays(deadline string) (natureDays int, workDays int) {
	ddlTime, err := time.Parse(cst.LayoutYearMonthDay, deadline)
	if err != nil {
		panic(err)
	}

	return times.CalcBetweenDays(time.Now(), ddlTime)
}
