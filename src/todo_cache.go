package zodo

import (
	"encoding/json"
)

type cache struct {
	data []*todo
}

func (c *cache) refresh() {
	newData := make([]*todo, 0)
	for _, line := range readTodoLines(Config.Storage.Type) {
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
	RewriteLinesToPath(backupPath, readTodoLines(Config.Storage.Type))

	lines := make([]string, 0)
	for _, td := range c.data {
		js, err := json.Marshal(td)
		if err != nil {
			panic(err)
		}
		lines = append(lines, string(js))
	}
	writeTodoLines(lines, Config.Storage.Type)
}

func (c *cache) list(keyword string, status []string, allStatus bool) []todo {
	tds := make([]todo, 0)
	for _, td := range sortTodo(c.data) {
		if td.ParentId == 0 && hitTodo(td, keyword) {
			walkTodo(td, &tds, 0, status, allStatus)
		}
	}
	return tds
}

func (c *cache) _map() map[int]*todo {
	m := make(map[int]*todo, 0)
	for _, td := range c.data {
		m[td.Id] = td
	}
	return m
}

func (c *cache) add(td todo) {
	c.data = append(c.data, &td)
}

func (c *cache) remove(id int) {
	m := c._map()
	toRemove := m[id]
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

	parent := m[toRemove.ParentId]
	if parent != nil {
		delete(parent.Children, id)
	}

	if toRemove.hasChildren() {
		for childId := range toRemove.Children {
			c.remove(childId)
		}
	}
}

func (c *cache) clear() {
	for _, td := range c.data {
		if td.Status == statusDone {
			c.remove(td.Id)
		}
	}
}

var cc cache

func InitCache() {
	cc = cache{}
	cc.refresh()
}
