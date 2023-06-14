package dev

import (
	"os"
	"strings"
	"zodo/src"
)

const fileName = "simplified.sql"

const comment = "COMMENT"

var ignorePrefixes = []string{
	"INDEX", "UNIQUE INDEX", "KEY", "UNIQUE KEY", "create_time", "update_time",
}

func SimplifySql(path string) {
	sqls := zodo.ReadLinesFromPath(path)
	handled := make([]string, 0)
	for _, sql := range sqls {
		sql = strings.TrimSpace(sql)
		if hasIgnorePrefix(sql) {
			continue
		}
		if strings.HasPrefix(sql, "`") {
			sql = strings.TrimPrefix(sql, "`")
			i := strings.Index(sql, "`")
			if !strings.Contains(sql, comment) {
				sql = sql[:i] + ","
			} else {
				j := strings.Index(sql, comment)
				sql = sql[:i] + sql[j+len(comment):]
			}
		}
		handled = append(handled, sql)
	}
	zodo.WriteLinesToPath(zodo.CurrentPath(fileName), handled, os.O_RDWR|os.O_TRUNC)
}

func hasIgnorePrefix(sql string) bool {
	for _, p := range ignorePrefixes {
		if strings.Contains(sql, p) {
			return true
		}
	}
	return false
}
