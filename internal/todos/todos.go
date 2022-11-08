package todos

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"time"
	"zodo/internal/conf"
	"zodo/internal/cst"
	"zodo/internal/emails"
	"zodo/internal/errs"
	"zodo/internal/ids"
	"zodo/internal/stdout"
	"zodo/internal/times"
)

const (
	statusPending    = "Pending"
	statusProcessing = "Processing"
	statusDone       = "Done"
)

var statusPriority = map[string]int{
	statusDone:       0,
	statusPending:    1,
	statusProcessing: 2,
}

func List(keyword string) {
	rows := make([]table.Row, 0)
	for _, td := range list(keyword) {
		content := td.Content
		if td.Level > 0 {
			content = fmt.Sprintf("%s|-%s", padding(td.Level, "  "), content)
		}
		ddl, remain := td.getDeadLineAndRemain(true)
		rows = append(rows, table.Row{
			td.Id,
			content,
			td.getStatus(true),
			ddl,
			remain,
		})
	}
	stdout.PrintTable(table.Row{"Id", "Content", "Status", "Deadline", "Remain"}, rows)
}

func Detail(id int) {
	td := _map()[id]
	if td == nil {
		return
	}

	rows := make([]table.Row, 0)
	rows = append(rows, table.Row{"Id", td.Id})
	rows = append(rows, table.Row{"Content", td.Content})
	rows = append(rows, table.Row{"Status", td.getStatus(true)})
	ddl, remain := td.getDeadLineAndRemain(true)
	rows = append(rows, table.Row{"Deadline", ddl})
	rows = append(rows, table.Row{"Remain", remain})
	rows = append(rows, table.Row{"Remark", td.Remark})
	rows = append(rows, table.Row{"RemindTime", td.RemindTime})
	rows = append(rows, table.Row{"RemindStatus", td.RemindStatus})
	rows = append(rows, table.Row{"LoopType", td.LoopType})
	rows = append(rows, table.Row{"Create", td.getCreateTime()})
	rows = append(rows, table.Row{"Parent", td.getParentId()})
	rows = append(rows, table.Row{"Children", td.getChildren()})
	stdout.PrintTable(table.Row{"Item", "Val"}, rows)
}

func DailyReport() error {
	load()
	var text string
	for _, td := range list("") {
		ddl, remain := td.getDeadLineAndRemain(false)
		if td.Level == 0 {
			text += "\n"
			if ddl != "" {
				text += fmt.Sprintf("* %s  %s, deadline %s, remain %s\n", td.Content, td.getStatus(false), ddl, remain)
			} else {
				text += fmt.Sprintf("* %s  %s\n", td.Content, td.getStatus(false))
			}
		} else {
			if ddl != "" {
				text += fmt.Sprintf("%s  |- %s  %s, deadline %s, remain %s\n",
					padding(td.Level, "    "), td.Content, td.getStatus(false), ddl, remain)
			} else {
				text += fmt.Sprintf("%s  |- %s  %s\n",
					padding(td.Level, "    "), td.Content, td.getStatus(false))
			}
		}
	}
	if text != "" {
		return emails.Send("Daily Report", text)
	}
	return nil
}

func Save() {
	save()
}

func Add(content string) (int, error) {
	if content == "" {
		return -1, &errs.InvalidInputError{
			Message: fmt.Sprint("content empty"),
		}
	}
	id := ids.GetAndSet(conf.Data.Storage.Type)
	add(todo{
		Id:         id,
		Content:    content,
		Status:     statusPending,
		CreateTime: time.Now().Format(cst.LayoutDateTime),
	})
	return id, nil
}

func Delete(ids []int) {
	for _, id := range ids {
		_delete(id)
	}
}

func Modify(id int, content string) {
	if content == "" {
		return
	}
	td := _map()[id]
	if td != nil {
		td.Content = content
	}
}

func Rollback() {
	rollback()
}

func Transfer() {
	transfer()
}

func SetDeadline(id int, deadline string) {
	td := _map()[id]
	if td != nil {
		td.Deadline = deadline
	}
}

func SetRemark(id int, remark string) {
	td := _map()[id]
	if td != nil {
		td.Remark = remark
	}
}

func SetChild(parentId int, childIds []int, append bool) error {
	m := _map()
	parent := m[parentId]
	if parent == nil {
		return &errs.NotFoundError{
			Target:  "parent",
			Message: fmt.Sprintf("parentId: %d", parentId),
		}
	}
	if parent.Children != nil && !append {
		for childId, _ := range parent.Children {
			child := m[childId]
			if child == nil {
				fmt.Println(&errs.NotFoundError{
					Target:  "child",
					Message: fmt.Sprintf("childId: %d", childId),
				})
			} else {
				child.ParentId = 0
			}
		}
	}
	if parent.Children == nil || !append {
		parent.Children = make(map[int]bool, 0)
	}
	for _, childId := range childIds {
		child := m[childId]
		if child == nil {
			fmt.Println(&errs.NotFoundError{
				Target:  "child",
				Message: fmt.Sprintf("childId: %d", childId),
			})
			continue
		}

		oldParent := m[child.ParentId]
		if oldParent != nil {
			delete(oldParent.Children, childId)
		}

		child.ParentId = parentId
		parent.Children[childId] = true
	}
	return nil
}

func SetPending(id int) {
	modifyStatus(id, statusPending)
}

func SetProcessing(id int) {
	modifyStatus(id, statusProcessing)
}

func SetDone(id int) {
	modifyStatus(id, statusDone)
}

func modifyStatus(id int, status string) {
	td := _map()[id]
	if td != nil {
		td.Status = status
	}
}

func calcRemainDays(deadline string) (natureDays int, workDays int) {
	ddlTime, err := time.Parse(cst.LayoutYearMonthDay, deadline)
	if err != nil {
		panic(err)
	}

	return times.CalcBetweenDays(time.Now(), ddlTime)
}
