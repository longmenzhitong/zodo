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

func List(keyword string) {
	rows := make([]table.Row, 0)
	for _, td := range data.list(keyword) {
		content := td.Content
		if td.Level > 0 {
			content = fmt.Sprintf("%s|-%s", padding(td.Level), content)
		}
		ddl, remain := td.getDeadLineAndRemain()
		rows = append(rows, table.Row{
			td.Id,
			content,
			td.getStatus(),
			ddl,
			remain,
		})
	}
	stdout.PrintTable(table.Row{"Id", "Content", "Status", "Deadline", "Remain"}, rows)
}

func Detail(id int) {
	td := data.Map[id]
	if td == nil {
		return
	}

	rows := make([]table.Row, 0)
	rows = append(rows, table.Row{"Id", td.Id})
	rows = append(rows, table.Row{"Content", td.Content})
	rows = append(rows, table.Row{"Status", td.getStatus()})
	ddl, remain := td.getDeadLineAndRemain()
	rows = append(rows, table.Row{"Deadline", ddl})
	rows = append(rows, table.Row{"Remain", remain})
	rows = append(rows, table.Row{"Remark", td.Remark})
	rows = append(rows, table.Row{"Create", td.getCreateTime()})
	rows = append(rows, table.Row{"Parent", td.getParentId()})
	rows = append(rows, table.Row{"Children", td.getChildren()})
	stdout.PrintTable(table.Row{"Item", "Val"}, rows)
}

func DailyReport() error {
	data.load()
	var text string
	for _, td := range data.List {
		if td.Status == statusDone {
			continue
		}

		// TODO 内容格式优化
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
	return emails.Send("Daily Report", text)
}

func Save() {
	data.save()
}

func Add(content string) (int, error) {
	if content == "" {
		return -1, &errs.InvalidInputError{
			Input:   content,
			Message: fmt.Sprintf("content empty"),
		}
	}
	id := ids.GetAndSet(conf.Data.Storage.Type)
	data.add(todo{
		Id:         id,
		Content:    content,
		Status:     statusPending,
		CreateTime: time.Now().Format(cst.LayoutDateTime),
	})
	return id, nil
}

func Delete(ids []int) {
	for _, id := range ids {
		data.delete(id)
	}
}

func Modify(id int, content string) {
	if content == "" {
		return
	}
	td := data.Map[id]
	if td != nil {
		td.Content = content
	}
}

func Transfer() {
	data.transfer()
}

func SetDeadline(id int, deadline string) {
	td := data.Map[id]
	if td != nil {
		td.Deadline = deadline
	}
}

func SetRemark(id int, remark string) {
	td := data.Map[id]
	if td != nil {
		td.Remark = remark
	}
}

func SetChild(parentId int, childIds []int, append bool) error {
	parent := data.Map[parentId]
	if parent == nil {
		return &errs.NotFoundError{
			Target:  "parent",
			Message: fmt.Sprintf("parentId: %d", parentId),
		}
	}
	if parent.Children != nil && !append {
		for childId, _ := range parent.Children {
			child := data.Map[childId]
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
		child := data.Map[childId]
		if child == nil {
			fmt.Println(&errs.NotFoundError{
				Target:  "child",
				Message: fmt.Sprintf("childId: %d", childId),
			})
			continue
		}

		oldParent := data.Map[child.ParentId]
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
	td := data.Map[id]
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
