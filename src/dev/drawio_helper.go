package dev

import (
	"fmt"
	"os"
	"strings"
	zodo "zodo/src"
)

const fileName = "simplified.sql"

const commentKeyword = "COMMENT"

var ignoreKeywords = []string{
	"INDEX", "UNIQUE INDEX", "KEY", "UNIQUE KEY", "create_time", "update_time",
}

func SimplifySql(path string) {
	sqls := zodo.ReadLinesFromPath(path)
	handled := make([]string, 0)
	var createTableLineNum int
	for _, sql := range sqls {
		sql = strings.TrimSpace(sql)
		if hasIngnoreKeyword(sql) {
			continue
		}

		if isCreateTableLine(sql) {
			createTableLineNum = len(handled)
		} else if isTableNameLine(sql) {
			if strings.Contains(sql, commentKeyword) {
				tableName := getTableName(sql)
				handled[createTableLineNum] = appendTableName(handled[createTableLineNum], tableName)
			}
			sql += "\n"
		} else if strings.HasPrefix(sql, "`") {
			sql = strings.TrimPrefix(sql, "`")
			i := strings.Index(sql, "`")
			if !strings.Contains(sql, commentKeyword) {
				sql = sql[:i] + ","
			} else {
				j := strings.Index(sql, commentKeyword)
				sql = sql[:i] + sql[j+len(commentKeyword):]
			}
		}

		handled = append(handled, sql)
	}
	zodo.WriteLinesToPath(zodo.CurrentPath(fileName), handled, os.O_RDWR|os.O_TRUNC)
	for _, s := range handled {
		fmt.Println(s)
	}
	zodo.WriteLinesToClipboard(handled)
}

func hasIngnoreKeyword(sql string) bool {
	for _, k := range ignoreKeywords {
		if strings.Contains(sql, k) {
			return true
		}
	}
	return false
}

func isCreateTableLine(sql string) bool {
	return strings.HasPrefix(sql, "CREATE")
}

func isTableNameLine(sql string) bool {
	return strings.Contains(sql, "ENGINE")
}

func getTableName(sql string) string {
	i := strings.LastIndex(sql, "'")
	j := strings.LastIndex(sql[:i-1], "'")
	return sql[j+1 : i]
}

func appendTableName(sql, tableName string) string {
	i := strings.LastIndex(sql, "`")
	return fmt.Sprintf("%s(%s)%s", sql[:i], tableName, sql[i:])
}
