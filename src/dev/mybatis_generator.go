package dev

import (
	"fmt"
	"regexp"
	"strings"
	zodo "zodo/src"
)

const classNameKeyword = "class"

func GenerateMybatisCode(path string) error {
	var packageName string
	var className string
	fieldNameMap := make(map[string]string)

	lines := zodo.ReadLinesFromPath(path)
	for _, line := range lines {
		if strings.HasPrefix(line, "import") {
			continue
		}

		if strings.HasPrefix(line, "package") {
			packageName = parsePackageName(line)
		} else if isClassNameLine(line) {
			className = parseClassName(line)
		} else {
			name := parseFieldName(line)
			if name != "" {
				fieldNameMap[name] = camelToSnake(name)
			}
		}
	}

	result := make([]string, 0)
	// result map
	result = append(result, fmt.Sprintf("<resultMap id=\"%sResultMap\" type=\"%s.%s\">", className, packageName, className))
	for p, c := range fieldNameMap {
		result = append(result, fmt.Sprintf("%s<result column=\"%s\" property=\"%s\"/>", indent(), c, p))
	}
	result = append(result, fmt.Sprintf("</resultMap>"))
	// result column
	result = append(result, "")
	result = append(result, fmt.Sprintf("<sql id=\"%sResultColumn\">", className))
	var columns string
	for _, c := range fieldNameMap {
		columns += fmt.Sprintf("`%s`, ", c)
	}
	result = append(result, fmt.Sprintf("%s%s", indent(), strings.TrimSuffix(columns, ", ")))
	result = append(result, fmt.Sprintf("</sql>"))
	for _, s := range result {
		fmt.Println(s)
	}

	return zodo.WriteLinesToClipboard(result)
}

func parsePackageName(line string) string {
	r := regexp.MustCompile(`package\s+([\w.]+)\s*;`)
	m := r.FindStringSubmatch(line)
	if len(m) > 1 {
		return m[1]
	}
	return ""
}

func isClassNameLine(line string) bool {
	return strings.Contains(line, classNameKeyword) && strings.Contains(line, "{")
}

func parseClassName(line string) string {
	i := strings.Index(line, classNameKeyword)
	s := line[i+len(classNameKeyword)+1:]
	j := strings.Index(s, " ")
	return s[:j]
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
	re := regexp.MustCompile(`([A-Z])`)
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
