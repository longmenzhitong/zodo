package todo

import (
	"encoding/json"
	"github.com/jedib0t/go-pretty/v6/table"
	"time"
	"zodo/internal/cst"
	"zodo/internal/file"
	"zodo/internal/ids"
	"zodo/internal/stdout"
)

const (
	fileName = "todo"
)

const (
	statusPending   = "Pending"
	statusDone      = "Done"
	statusAbandoned = "Abandoned"
	statusDeleted   = "Deleted"
)

type Todo struct {
	Id         int
	Content    string
	Status     string
	CreateTime string
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

func FindById(id int) *Todo {
	for i := range Todos {
		if Todos[i].Id == id {
			return &Todos[i]
		}
	}
	return nil
}

func List() {
	rows := make([]table.Row, 0)
	for _, td := range Todos {
		if td.Status == statusDeleted {
			continue
		}
		rows = append(rows, table.Row{td.Id, td.Content, td.Status, td.CreateTime})
	}
	stdout.PrintTable(table.Row{"Id", "Content", "Status", "Create"}, rows)
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

func Modify(id int, content string, status string) {
	td := FindById(id)
	if td == nil {
		return
	}
	if content != "" {
		td.Content = content
	}
	if status != "" {
		td.Status = status
	}
}

func Pending(id int) {
	Modify(id, "", statusPending)
}

func Done(id int) {
	Modify(id, "", statusDone)
}
func Abandon(id int) {
	Modify(id, "", statusAbandoned)
}

func Delete(id int) {
	Modify(id, "", statusDeleted)
}
