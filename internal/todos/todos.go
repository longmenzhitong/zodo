package todos

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"sort"
	"time"
	"zodo/internal/conf"
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
		if td.ParentId == 0 {
			walkTree(td, &rows, "")
		}
	}
	stdout.PrintTable(table.Row{"Id", "Content", "Status", "Deadline"}, rows)
}

func walkTree(td *todo, rows *[]table.Row, tab string) {
	if td == nil {
		return
	}
	if !param.All && td.Status == statusDone {
		return
	}
	content := td.Content
	if td.ParentId != 0 {
		content = fmt.Sprintf("%s|-%s", tab, content)
	}
	*rows = append(*rows, table.Row{
		td.Id,
		content,
		td.getStatus(),
		td.getDeadLine(),
	})
	if td.Children == nil || len(td.Children) == 0 {
		return
	}
	childList := make([]*todo, 0)
	for childId, _ := range td.Children {
		child := Data.Map[childId]
		if child == nil {
			fmt.Println(&errs.NotFoundError{
				Target:  "child",
				Message: fmt.Sprintf("parentId: %d, childId: %d", td.Id, childId),
			})
		} else {
			childList = append(childList, Data.Map[childId])
		}
	}
	sort.Slice(childList, func(i, j int) bool {
		a := childList[i]
		b := childList[j]
		if a.Deadline != "" && b.Deadline != "" {
			ta, err := time.Parse(cst.LayoutYearMonthDay, a.Deadline)
			if err != nil {
				panic(err)
			}
			tb, err := time.Parse(cst.LayoutYearMonthDay, b.Deadline)
			if err != nil {
				panic(err)
			}
			return ta.Unix() < tb.Unix()
		}
		if a.Deadline == "" && b.Deadline == "" {
			if a.Status != b.Status {
				return a.Status == statusProcessing
			}
			return a.Id < b.Id
		}
		return a.Deadline != ""
	})
	for _, child := range childList {
		walkTree(child, rows, tab+"  ")
	}
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
	rows = append(rows, table.Row{"Children", td.getChildren()})
	stdout.PrintTable(table.Row{"Item", "Val"}, rows)
}

func DailyReport() error {
	Data.load()
	var text string
	for _, td := range Data.List {
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

func Add(content string) (int, error) {
	if content == "" {
		return -1, &errs.InvalidInputError{
			Input:   content,
			Message: fmt.Sprintf("content empty"),
		}
	}
	id := ids.GetAndSet(conf.Data.Storage.Type)
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

func Transfer() {
	Data.Transfer()
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

func SetChild(parentId int, childIds []int, append bool) error {
	parent := Data.Map[parentId]
	if parent == nil {
		return &errs.NotFoundError{
			Target:  "parent",
			Message: fmt.Sprintf("parentId: %d", parentId),
		}
	}
	if parent.Children != nil && !append {
		for childId, _ := range parent.Children {
			child := Data.Map[childId]
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
		child := Data.Map[childId]
		if child == nil {
			fmt.Println(&errs.NotFoundError{
				Target:  "child",
				Message: fmt.Sprintf("childId: %d", childId),
			})
			continue
		}

		oldParent := Data.Map[child.ParentId]
		if oldParent != nil {
			delete(oldParent.Children, childId)
		}

		child.ParentId = parentId
		parent.Children[childId] = true
	}
	Data.save()
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
