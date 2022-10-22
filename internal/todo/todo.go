package todo

import (
	"encoding/json"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"time"
	"zodo/internal/color"
	"zodo/internal/cst"
	"zodo/internal/file"
	"zodo/internal/ids"
	"zodo/internal/stdout"
	"zodo/internal/times"
)

const (
	fileName = "todo"
)

const (
	statusPending    = "Pending"
	statusProcessing = "Processing"
	statusDone       = "Done"
	statusDeleted    = "Deleted"
)

type Todo struct {
	Id         int
	Content    string
	Status     string
	Deadline   string
	CreateTime string
}

func (t *Todo) GetStatus() string {
	switch t.Status {
	case statusPending:
		return color.Purple(t.Status)
	case statusProcessing:
		return color.Cyan(t.Status)
	default:
		return t.Status
	}
}

func (t *Todo) GetDeadLine() string {
	if t.Deadline == "" {
		return ""
	}
	nd, wd := calcRemainDays(t.Deadline)
	ddl := fmt.Sprintf("%s (%dnd/%dwd)", t.Deadline, nd, wd)
	ddl = times.Simplify(ddl)

	if t.Status == statusPending || t.Status == statusProcessing {
		if wd == 0 {
			ddl = color.Red(ddl)
		} else if wd == 1 {
			ddl = color.Yellow(ddl)
		} else {
			ddl = color.Green(ddl)
		}
	}
	return ddl
}

var (
	Path  string
	Todos []Todo
)

func init() {
	Path = file.Dir + cst.PathSep + fileName
	file.EnsureExist(Path)

	Todos = make([]Todo, 0)
	for _, line := range file.ReadLinesFromPath(Path) {
		var td Todo
		err := json.Unmarshal([]byte(line), &td)
		if err != nil {
			panic(err)
		}
		Todos = append(Todos, td)
	}
}

func Save() {
	lines := make([]string, 0)
	for _, td := range Todos {
		js, err := json.Marshal(td)
		if err != nil {
			panic(err)
		}
		lines = append(lines, string(js))
	}
	file.RewriteLinesToPath(Path, lines)
}

func List() {
	rows := make([]table.Row, 0)
	for _, td := range Todos {
		if td.Status == statusDeleted || td.Status == statusDone {
			continue
		}

		rows = append(rows, table.Row{
			td.Id,
			td.Content,
			td.GetStatus(),
			td.GetDeadLine(),
			times.Simplify(td.CreateTime),
		})
	}
	stdout.PrintTable(table.Row{"Id", "Content", "Status", "Deadline", "Create"}, rows)
}

func Add(content string) {
	td := Todo{
		Id:         ids.Get(),
		Content:    content,
		Status:     statusPending,
		CreateTime: time.Now().Format(cst.LayoutDateTime),
	}
	Todos = append(Todos, td)
}

func Modify(id int, content string) {
	td := findById(id)
	if td == nil {
		return
	}
	td.Content = content
}

func Deadline(id int, deadline string) {
	td := findById(id)
	if td == nil {
		return
	}
	td.Deadline = deadline
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
	if td == nil {
		return
	}
	td.Status = status
}

func findById(id int) *Todo {
	for i := range Todos {
		if Todos[i].Id == id {
			return &Todos[i]
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
