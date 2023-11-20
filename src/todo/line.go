package todo

import (
	zodo "zodo/src"
)

func readTodoLines() []string {
	return zodo.ReadLinesFromPath(path)
}

func writeTodoLines(lines []string) {
	zodo.RewriteLinesToPath(path, lines)
	return
}
