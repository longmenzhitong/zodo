package zodo

import (
	"os"
	"strings"
)

const (
	PathSep                      = string(os.PathSeparator)
	LayoutDate                   = "2006-01-02"
	LayoutDateTime               = "2006-01-02 15:04:05"
	LayoutMonthDay               = "01-02"
	LayoutHourMinute             = "15:04"
	LayoutMonthDayHourMinute     = "01-02 15:04"
	LayoutYearMonthDayHourMinute = "2006-01-02 15:04"
)

func homeDir() string {
	d, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return d
}

func ProjectDir() string {
	return homeDir() + PathSep + ".config" + PathSep + "zodo"
}

func Path(filename string) string {
	return ProjectDir() + PathSep + filename
}

func CurrentPath(filename string) string {
	return "." + PathSep + filename
}

func CurrentDirName() string {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	i := strings.LastIndex(pwd, PathSep)
	return pwd[i+1:]
}
