package dev

import (
	"fmt"
	"os"
	"strings"
	zodo "zodo/src"
)

const fileName = "simplified.sql"

const comment = "COMMENT"

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

		// 处理表名
		if isCreateTableLine(sql) {
			createTableLineNum = len(handled)
		}
		if isTableNameLine(sql) {
			tableName := getTableName(sql)
			handled[createTableLineNum] = appendTableName(handled[createTableLineNum], tableName)
		}

		// 处理字段
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
	for _, s := range handled {
		fmt.Println(s)
	}
	zodo.WriteToClipboard(handled)
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
	return fmt.Sprintf("%s(%s)%s", sql[:i-1], tableName, sql[i:])
}
