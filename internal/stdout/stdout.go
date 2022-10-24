package stdout

import (
	"github.com/jedib0t/go-pretty/v6/table"
	"os"
)

func PrintTable(header table.Row, rows []table.Row) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	// TODO 表格长度从配置读取
	t.SetAllowedRowLength(150)
	t.AppendHeader(header)
	for _, row := range rows {
		t.AppendRow(row)
		t.AppendSeparator()
	}
	t.Render()
}
