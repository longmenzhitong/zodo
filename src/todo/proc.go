package todo

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"time"
	"zodo/src"
)

func List(keyword string, status []string, allStatus bool) {
	rows := make([]table.Row, 0)
	for _, td := range cc.list(keyword, status, allStatus) {
		content := td.Content
		if td.Level > 0 {
			content = fmt.Sprintf("%s|-%s", padding(td.Level), content)
		}
		stat := td.getStatus(true)
		ddl, remain := td.getDeadLineAndRemain(true)
		if td.hasChildren() && !zodo.Config.Todo.ShowParent {
			stat = ""
			ddl = ""
			remain = ""
		}
		rows = append(rows, table.Row{
			td.Id,
			content,
			stat,
			ddl,
			remain,
		})
	}
	zodo.PrintTable(table.Row{"Id", "Content", "Status", "Deadline", "Remain"}, rows)
}

func Detail(id int) error {
	td := cc._map()[id]
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
	zodo.PrintTable(table.Row{"Item", "Val"}, rows)
	return nil
}

func Add(content string) (int, error) {
	if content == "" {
		return -1, &zodo.InvalidInputError{
			Message: fmt.Sprint("empty content"),
		}
	}
	id := zodo.Id(zodo.Config.Storage.Type)
	cc.add(todo{
		Id:         id,
		Content:    content,
		Status:     statusPending,
		CreateTime: time.Now().Format(zodo.LayoutDateTime),
	})
	return id, nil
}

func Modify(id int, content string) {
	if content == "" {
		return
	}
	td := cc._map()[id]
	if td != nil {
		td.Content = content
	}
}

func Remove(ids []int) {
	for _, id := range ids {
		cc.remove(id)
	}
}

func SetDeadline(id int, deadline string) {
	td := cc._map()[id]
	if td != nil {
		td.Deadline = deadline
	}
}

func SetRemark(id int, remark string) {
	td := cc._map()[id]
	if td != nil {
		td.Remark = remark
	}
}

func SetChild(parentId int, childIds []int, append bool) error {
	m := cc._map()
	parent := m[parentId]
	if parent == nil {
		return &zodo.NotFoundError{
			Target:  "parent",
			Message: fmt.Sprintf("parentId: %d", parentId),
		}
	}
	if parent.Children != nil && !append {
		for childId := range parent.Children {
			child := m[childId]
			if child == nil {
				fmt.Println(&zodo.NotFoundError{
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
			fmt.Println(&zodo.NotFoundError{
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
	td := cc._map()[id]
	if td != nil {
		td.Status = statusPending
	}
}

func SetProcessing(id int) {
	td := cc._map()[id]
	if td != nil {
		td.Status = statusProcessing
	}
}

func SetDone(id int) {
	td := cc._map()[id]
	if td == nil {
		return
	}
	td.Status = statusDone
	td.DoneTime = time.Now().Format(zodo.LayoutDateTime)
	if !td.hasChildren() {
		return
	}
	for childId := range td.Children {
		SetDone(childId)
	}
}

func SetHiding(id int) {
	td := cc._map()[id]
	if td == nil {
		return
	}
	td.Status = statusHiding
	if !td.hasChildren() {
		return
	}
	for childId := range td.Children {
		SetHiding(childId)
	}
}

func Save() {
	cc.save()
}

func Report() error {
	cc.refresh()
	var text string
	for _, td := range cc.list("", []string{}, false) {
		status := td.getStatus(false)
		ddl, remain := td.getDeadLineAndRemain(false)
		if td.hasChildren() && !zodo.Config.Todo.ShowParent {
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

func Rollback() {
	writeTodoLines(zodo.ReadLinesFromPath(backupPath), zodo.Config.Storage.Type)
}

func Transfer() {
	switch zodo.Config.Storage.Type {
	case zodo.StorageTypeFile:
		writeTodoLines(readTodoLines(zodo.StorageTypeRedis), zodo.StorageTypeFile)
		zodo.SetId(zodo.GetId(zodo.StorageTypeRedis)+1, zodo.StorageTypeFile)
		return
	case zodo.StorageTypeRedis:
		writeTodoLines(readTodoLines(zodo.StorageTypeFile), zodo.StorageTypeRedis)
		zodo.SetId(zodo.GetId(zodo.StorageTypeFile)+1, zodo.StorageTypeRedis)
		return
	default:
		panic(&zodo.InvalidConfigError{
			Message: fmt.Sprintf("storage.type: %s", zodo.Config.Storage.Type),
		})
	}
}

func ClearDoneTodo() int {
	return cc.clearDoneTodo()
}

func DefragId() int {
	return cc.defragId()
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
