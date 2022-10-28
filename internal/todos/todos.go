package todos

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"time"
	"zodo/internal/cst"
	"zodo/internal/emails"
	"zodo/internal/errs"
	"zodo/internal/ids"
	"zodo/internal/param"
	"zodo/internal/stdout"
	"zodo/internal/times"
)

const (
	statusPending    = "Pending"
	statusProcessing = "Processing"
	statusDone       = "Done"
)

func List() {
	rows := make([]table.Row, 0)
	for _, td := range Data.List {
		if !param.All && td.Status == statusDone {
			continue
		}
		if td.ParentId != 0 {
			continue
		}

		if td.Childs == nil || len(td.Childs) == 0 {
			rows = append(rows, table.Row{
				td.Id,
				td.Content,
				td.getStatus(),
				td.getDeadLine(),
			})
		} else {
			rows = append(rows, table.Row{
				td.Id,
				td.Content,
				"",
				"",
			})
			for childId := range td.Childs {
				child := Data.Map[childId]
				if child != nil && child.Status != statusDone {
					rows = append(rows, table.Row{
						child.Id,
						fmt.Sprintf(" |-%s", child.Content),
						child.getStatus(),
						child.getDeadLine(),
					})
				}
			}
		}
	}
	stdout.PrintTable(table.Row{"Id", "Content", "Status", "Deadline"}, rows)
}

func Detail(id int) {
	td := Data.Map[id]
	if td == nil {
		return
	}

	rows := make([]table.Row, 0)
	rows = append(rows, table.Row{"Id", td.Id})
	rows = append(rows, table.Row{"Content", td.Content})
	rows = append(rows, table.Row{"Status", td.getStatus()})
	rows = append(rows, table.Row{"Deadline", td.getDeadLine()})
	rows = append(rows, table.Row{"Remark", td.Remark})
	rows = append(rows, table.Row{"Create", td.getCreateTime()})
	rows = append(rows, table.Row{"Parent", td.getParentId()})
	rows = append(rows, table.Row{"Child", td.getChilds()})
	stdout.PrintTable(table.Row{"Item", "Val"}, rows)
}

func DailyReport() error {
	var text string
	for _, td := range Data.List {
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
	return emails.Send("Daily Report", text)
}

func Add(content string) (int, error) {
	if content == "" {
		return -1, &errs.InvalidInputError{
			Input:   content,
			Message: fmt.Sprintf("content empty"),
		}
	}
	id := ids.Get()
	Data.add(todo{
		Id:         id,
		Content:    content,
		Status:     statusPending,
		CreateTime: time.Now().Format(cst.LayoutDateTime),
	})
	return id, nil
}

func Delete(ids []int) {
	for _, id := range ids {
		Data.delete(id)
	}
}

func Modify(id int, content string) {
	if content == "" {
		return
	}
	td := Data.Map[id]
	if td != nil {
		td.Content = content
	}
	Data.save()
}

func SetDeadline(id int, deadline string) {
	td := Data.Map[id]
	if td != nil {
		td.Deadline = deadline
	}
	Data.save()
}

func SetRemark(id int, remark string) {
	td := Data.Map[id]
	if td != nil {
		td.Remark = remark
	}
	Data.save()
}

func SetChild(parentId int, childIds []int) error {
	parent := Data.Map[parentId]
	if parent == nil {
		return &errs.NotFoundError{
			Target:  "parent",
			Message: fmt.Sprintf("parentId: %d", parentId),
		}
	}
	if parent.Childs == nil {
		parent.Childs = make(map[int]bool, 0)
	}
	for _, childId := range childIds {
		child := Data.Map[childId]
		if child == nil {
			return &errs.NotFoundError{
				Target:  "child",
				Message: fmt.Sprintf("childId: %d", childId),
			}
		}
		child.ParentId = parentId
		parent.Childs[childId] = true
	}
	Data.save()
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

func modifyStatus(id int, status string) {
	td := Data.Map[id]
	if td != nil {
		td.Status = status
	}
	Data.save()
}

func calcRemainDays(deadline string) (natureDays int, workDays int) {
	ddlTime, err := time.Parse(cst.LayoutYearMonthDay, deadline)
	if err != nil {
		panic(err)
	}

	return times.CalcBetweenDays(time.Now(), ddlTime)
}
