package todo

import (
	"encoding/json"
	"fmt"
	"sort"
	zodo "zodo/src"
)

var cc cache

type cache struct {
	data []*todo
}

func (c *cache) refresh() {
	newData := make([]*todo, 0)
	for _, line := range readTodoLines(zodo.Config.Storage.Type) {
		var td todo
		err := json.Unmarshal([]byte(line), &td)
		if err != nil {
			panic(err)
		}
		newData = append(newData, &td)
	}
	c.data = newData
}

func (c *cache) save() {
	// backup first
	zodo.RewriteLinesToPath(backupPath, readTodoLines(zodo.Config.Storage.Type))

	lines := make([]string, 0)
	for _, td := range c.data {
		js, err := json.Marshal(td)
		if err != nil {
			panic(err)
		}
		lines = append(lines, string(js))
	}
	writeTodoLines(lines, zodo.Config.Storage.Type)
}

func (c *cache) list(keyword string, status []string, allStatus bool) []todo {
	tds := make([]todo, 0)
	for _, td := range sortTodo(c.data) {
		if td.ParentId == 0 && hitKeyword(td, keyword) {
			walkTodo(td, &tds, 0, status, allStatus)
		}
	}
	return tds
}

func (c *cache) get(id int) *todo {
	for _, td := range c.data {
		if td.Id == id {
			return td
		}
	}
	return nil
}

func (c *cache) add(td todo) {
	if c.get(td.Id) != nil {
		panic(&zodo.InvalidInputError{Message: fmt.Sprintf("id duplicated: %d", td.Id)})
	}

	c.data = append(c.data, &td)
}

func (c *cache) remove(id int) {
	toRemove := cc.get(id)
	if toRemove == nil {
		return
	}

	newList := make([]*todo, 0)
	for _, td := range c.data {
		if td.Id != id {
			newList = append(newList, td)
		}
	}
	c.data = newList

	parent := cc.get(toRemove.ParentId)
	if parent != nil {
		delete(parent.Children, id)
	}

	if toRemove.hasChildren() {
		for childId := range toRemove.Children {
			c.remove(childId)
		}
	}
}

func (c *cache) clearDoneTodo() int {
	count := 0
	for _, td := range c.data {
		if td.Status == statusDone {
			c.remove(td.Id)
			count++
		}
	}
	return count
}

func (c *cache) defragId() int {
	sort.Slice(c.data, func(i, j int) bool {
		return c.data[i].Id < c.data[j].Id
	})
	m := make(map[int]int, 0)
	for i, td := range c.data {
		m[td.Id] = i + 1
	}
	for _, td := range c.data {
		td.Id = m[td.Id]
		if td.ParentId != 0 {
			td.ParentId = m[td.ParentId]
		}
		if td.hasChildren() {
			newChildren := make(map[int]bool, 0)
			for childId := range td.Children {
				newChildren[m[childId]] = true
			}
			td.Children = newChildren
		}
	}
	oldNextId := zodo.GetId(zodo.Config.Storage.Type)
	newNextId := len(c.data) + 1
	zodo.SetId(newNextId, zodo.Config.Storage.Type)
	return oldNextId - newNextId
}

func InitCache() {
	cc.refresh()
}
