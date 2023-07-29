package todo

import (
	"fmt"
	"time"
	zodo "zodo/src"

	"github.com/jedib0t/go-pretty/v6/table"
)

func List(keyword string, allStatus bool) {
	rows := make([]table.Row, 0)
	showDeadline := false
	for _, td := range Cache.list(keyword, allStatus) {
		content := td.Content
		if td.Remark != "" {
			content += zodo.ColoredString(zodo.ColorBlue, "*")
		}
		if td.Level > 0 {
			content = fmt.Sprintf("%s|- %s", padding(td.Level), content)
		}
		stat := td.getStatus(true)

		ddl, remain := td.getDeadLineAndRemain(true)
		if td.hasChildren() && !zodo.Config.Todo.ShowParentStatus {
			stat = ""
			ddl = ""
			remain = ""
		}
		if ddl != "" {
			showDeadline = true
		}

		row := table.Row{td.Id, content, stat}
		if showDeadline {
			row = append(row, ddl)
			row = append(row, remain)
		}
		rows = append(rows, row)
	}

	title := table.Row{"Id", "Content", "Status"}
	if showDeadline {
		title = append(title, "Deadline")
		title = append(title, "Remain")
	}
	zodo.PrintTable(&title, rows)
}

func Detail(id int) error {
	td := Cache.get(id)
	if td == nil {
		return &zodo.NotFoundError{
			Target:  "todo",
			Message: fmt.Sprintf("id: %d", id),
		}
	}

	rows := make([]table.Row, 0)
	rows = append(rows, table.Row{"Id", td.Id})
	rows = append(rows, table.Row{"Content", td.Content})
	rows = append(rows, table.Row{"Status", td.getStatus(true)})
	ddl, remain := td.getDeadLineAndRemain(true)
	rows = append(rows, table.Row{"Deadline", ddl})
	rows = append(rows, table.Row{"Remain", remain})
	rows = append(rows, table.Row{"Remark", td.Remark})
	rows = append(rows, table.Row{"RemindTime", td.getRemindTime()})
	rows = append(rows, table.Row{"RemindStatus", td.RemindStatus})
	rows = append(rows, table.Row{"LoopType", td.LoopType})
	rows = append(rows, table.Row{"DoneTime", td.getDoneTime()})
	rows = append(rows, table.Row{"CreateTime", td.getCreateTime()})
	rows = append(rows, table.Row{"Parent", td.getParentId()})
	rows = append(rows, table.Row{"Children", td.getChildren()})
	zodo.PrintTable(&table.Row{"Item", "Val"}, rows)
	return nil
}

func Add(content string) (int, error) {
	if content == "" {
		return -1, &zodo.InvalidInputError{
			Message: fmt.Sprint("empty content"),
		}
	}
	id := zodo.Id.SetGetNext()
	Cache.add(todo{
		Id:         id,
		Content:    content,
		Status:     StatusPending,
		CreateTime: time.Now().Format(zodo.LayoutDateTime),
	})
	return id, nil
}

func Modify(id int, content string) {
	if content == "" {
		return
	}

	td := Cache.get(id)
	if td != nil {
		td.Content = content
	}
}

func CopyContent(id int) error {
	td := Cache.get(id)
	if td != nil {
		return zodo.WriteLineToClipboard(td.Content)
	}
	return nil
}

func Remove(ids []int, recursively bool) {
	zodo.Id.Backup()

	for _, id := range ids {
		Cache.remove(id, recursively)
	}
}

func SetDeadline(id int, deadline string) {
	td := Cache.get(id)
	if td != nil {
		td.Deadline = deadline
	}
}

func CopyDeadline(id int) error {
	td := Cache.get(id)
	if td != nil {
		return zodo.WriteLineToClipboard(td.Deadline)
	}
	return nil
}

func SetRemark(id int, remark string) {
	td := Cache.get(id)
	if td != nil {
		td.Remark = remark
	}
}

func CopyRemark(id int) error {
	td := Cache.get(id)
	if td != nil {
		return zodo.WriteLineToClipboard(td.Remark)
	}
	return nil
}

func SetChild(parentId int, childIds []int, append bool) error {
	parent := Cache.get(parentId)
	if !append && parent.hasChildren() {
		for childId := range parent.Children {
			Cache.get(childId).ParentId = 0
		}
	}
	if !append || parent.Children == nil {
		parent.Children = make(map[int]bool, 0)
	}
	for _, childId := range childIds {
		child := Cache.get(childId)
		if child.ParentId != 0 {
			delete(Cache.get(child.ParentId).Children, childId)
		}

		child.ParentId = parentId
		parent.Children[childId] = true

		// for swap parent and child
		if parent.ParentId == childId {
			parent.ParentId = 0
		}
		if child.Children[parentId] {
			delete(child.Children, parentId)
		}
	}
	return nil
}

func SetStatus(id int, status Status) {
	td := Cache.get(id)
	if td == nil {
		return
	}
	td.Status = status
	if !td.hasChildren() {
		return
	}
	// FIXME: 应该影响子任务的状态吗？
	for childId := range td.Children {
		SetStatus(childId, status)
	}
}

func Save() {
	Cache.save()
}

func Report() error {
	Cache.refresh()
	var text string
	for _, td := range Cache.list("", false) {
		status := td.getStatus(false)
		ddl, remain := td.getDeadLineAndRemain(false)
		if td.hasChildren() && !zodo.Config.Todo.ShowParentStatus {
			status = ""
			ddl = ""
			remain = ""
		}
		if td.Level == 0 {
			text += "\n"
			if ddl != "" {
				text += fmt.Sprintf("* %s  %s, deadline %s, remain %s\n", td.Content, status, ddl, remain)
			} else {
				text += fmt.Sprintf("* %s  %s\n", td.Content, status)
			}
		} else {
			if ddl != "" {
				text += fmt.Sprintf("%s  |- %s  %s, deadline %s, remain %s\n", padding(td.Level), td.Content, status, ddl, remain)
			} else {
				text += fmt.Sprintf("%s  |- %s  %s\n", padding(td.Level), td.Content, status)
			}
		}
	}
	if text != "" {
		return zodo.SendEmail("Daily Report", text)
	}
	return nil
}

func padding(level int) string {
	var p string
	for i := 0; i < zodo.Config.Todo.Padding; i++ {
		p += " "
	}
	var res string
	for i := 0; i < level; i++ {
		res += p
	}
	return res
}

func Rollback() {
	writeTodoLines(zodo.ReadLinesFromPath(backupPath), zodo.Config.Storage.Type)
	zodo.Id.Rollback()
}

func ClearDoneTodo() int {
	return Cache.clearDoneTodo()
}

func DefragId() (int, int) {
	return Cache.defragId()
}

func Statistics() {
	proc := 0
	pend := 0
	done := 0
	for _, td := range Cache.data {
		if td.hasChildren() && !zodo.Config.Todo.ShowParentStatus {
			continue
		}
		switch td.Status {
		case StatusProcessing:
			proc++
		case StatusPending:
			pend++
		case StatusDone:
			done++
		}
	}

	rows := make([]table.Row, 0)
	rows = append(rows, table.Row{zodo.ColoredString(zodo.Config.Todo.Color.Status.Processing, string(StatusProcessing)), proc})
	rows = append(rows, table.Row{zodo.ColoredString(zodo.Config.Todo.Color.Status.Pending, string(StatusPending)), pend})
	rows = append(rows, table.Row{zodo.ColoredString(zodo.Config.Todo.Color.Status.Done, string(StatusDone)), done})
	rows = append(rows, table.Row{"NextId", zodo.Id.GetNext()})
	zodo.PrintTable(&table.Row{"Item", "Value"}, rows)
}
