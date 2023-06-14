package dev

import (
	"fmt"
	"github.com/atotto/clipboard"
	"regexp"
	"strings"
	"zodo/src"
)

const (
	packageNamePattern  = `package\s+([\w.]+)\s*;`
	classNamePattern    = `class\s+([\w.]+)\s*{`
	camelToSnakePattern = `([A-Z])`
)

func GenerateMybatisCode(path string) error {
	var packageName string
	var className string
	fieldNameMap := make(map[string]string)

	lines := zodo.ReadLinesFromPath(path)
	for _, line := range lines {
		if strings.HasPrefix(line, "import") {
			continue
		}

		var name string
		if strings.HasPrefix(line, "package") {
			name = parseName(packageNamePattern, line)
			if name != "" {
				packageName = name
			}
		} else if strings.Contains(line, "class") {
			name = parseName(classNamePattern, line)
			if name != "" {
				className = name
			}
		} else {
			name = parseFieldName(line)
			if name != "" {
				fieldNameMap[name] = camelToSnake(name)
			}
		}
	}

	result := make([]string, 0)
	// result map
	result = append(result, fmt.Sprintf("<resultMap id=\"%s\" type=\"%s.%s\">", className, packageName, className))
	for p, c := range fieldNameMap {
		result = append(result, fmt.Sprintf("%s<result column=\"%s\" property=\"%s\"/>", indent(), c, p))
	}
	result = append(result, fmt.Sprintf("</resultMap>"))
	// result column
	result = append(result, "")
	result = append(result, fmt.Sprintf("<sql id=\"%s_Column_List\">", className))
	var columns string
	for _, c := range fieldNameMap {
		columns += fmt.Sprintf("`%s`, ", c)
	}
	result = append(result, fmt.Sprintf("%s%s", indent(), strings.TrimSuffix(columns, ", ")))
	result = append(result, fmt.Sprintf("</sql>"))
	for _, s := range result {
		fmt.Println(s)
	}

	// 复制到剪切板
	var text string
	for _, s := range result {
		text += s + "\n"
	}
	err := clipboard.WriteAll(text)
	if err != nil {
		return err
	}

	fmt.Printf("\n(Copied.)\n")
	return nil
}

func parseName(pattern, line string) string {
	r := regexp.MustCompile(pattern)
	m := r.FindStringSubmatch(line)
	if len(m) > 1 {
		return m[1]
	}
	return ""
}

func parseFieldName(line string) string {
	if !strings.HasSuffix(line, ";") {
		return ""
	}
	i := strings.LastIndex(line, " ")
	if i == -1 {
		return ""
	}

	return line[i+1 : len(line)-1]
}

func camelToSnake(camel string) string {
	re := regexp.MustCompile(camelToSnakePattern)
	converted := re.ReplaceAllString(camel, "_$1")
	return strings.ToLower(converted)
}

func indent() string {
	c := 4
	var s string
	for i := 0; i < c; i++ {
		s += " "
	}
	return s
}
