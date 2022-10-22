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
	Remark     string
	CreateTime string
}

func (td *Todo) GetStatus() string {
	switch td.Status {
	case statusPending:
		return color.Purple(td.Status)
	case statusProcessing:
		return color.Cyan(td.Status)
	default:
		return td.Status
	}
}

func (td *Todo) GetDeadLine() string {
	if td.Deadline == "" {
		return ""
	}
	nd, wd := calcRemainDays(td.Deadline)
	ddl := fmt.Sprintf("%s (%dnd/%dwd)", td.Deadline, nd, wd)
	ddl = times.Simplify(ddl)

	if td.Status == statusPending || td.Status == statusProcessing {
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

func (td *Todo) GetCreateTime() string {
	return times.Simplify(td.CreateTime)
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
		if td.Status == statusDeleted {
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
	stdout.PrintTable(table.Row{"Item", "Value"}, rows)
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
