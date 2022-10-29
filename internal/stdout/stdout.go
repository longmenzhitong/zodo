package stdout

import (
	"github.com/jedib0t/go-pretty/v6/table"
	"os"
	"zodo/internal/conf"
)

func PrintTable(header table.Row, rows []table.Row) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetAllowedRowLength(conf.Data.Table.MaxLen)
	t.AppendHeader(header)
	for _, row := range rows {
		t.AppendRow(row)
		t.AppendSeparator()
	}
	t.Render()
}
