package dev

import (
	"fmt"
	"strings"
	zodo "zodo/src"

	"github.com/xuri/excelize/v2"
)

func GenerateJavaCode(path, name string, sheetIndex int) error {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return err
	}

	defer f.Close()

	sheetName := f.GetSheetName(sheetIndex)
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return err
	}

	fields := make([]string, 0)

	for _, row := range rows {
		for _, cell := range row {
			if strings.HasPrefix(cell, "{.") {
				fields = append(fields, cell[2:len(cell)-1])
			}
		}
	}

	dateHandled := false

	importLines := make([]string, 0)
	importLines = append(importLines, "import lombok.Data;")

	codeLines := make([]string, 0)
	codeLines = append(codeLines, "")
	codeLines = append(codeLines, "/**")
	codeLines = append(codeLines, fmt.Sprintf(" * %s导出DTO", sheetName))
	codeLines = append(codeLines, " */")
	codeLines = append(codeLines, "@Data")
	codeLines = append(codeLines, fmt.Sprintf("public class %s {", name))
	codeLines = append(codeLines, "")
	for _, field := range fields {
		var javaType string
		var dateFormat string
		if strings.Contains(strings.ToLower(field), "time") {
			javaType = "Date"
			dateFormat = "yyyy-MM-dd HH:mm:ss"
		} else if strings.Contains(strings.ToLower(field), "date") {
			javaType = "Date"
			dateFormat = "yyyy-MM-dd"
		} else {
			javaType = "String"
			dateFormat = ""
		}
		if dateFormat != "" {
			if !dateHandled {
				importLines = append(importLines, "import java.util.Date;")
				importLines = append(importLines, "import com.alibaba.excel.annotation.format.DateTimeFormat;")
				dateHandled = true
			}
			codeLines = append(codeLines, fmt.Sprintf("    @DateTimeFormat(\"%s\")", dateFormat))
		}
		codeLines = append(codeLines, fmt.Sprintf("    private %s %s;", javaType, field))
		codeLines = append(codeLines, "")
	}
	codeLines = append(codeLines, "}")

	lines := make([]string, 0)
	lines = append(lines, importLines...)
	lines = append(lines, codeLines...)

	for _, line := range lines {
		fmt.Println(line)
	}
	zodo.WriteToClipboard(lines)

	return nil
}
