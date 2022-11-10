package zodo

import "os"

const (
	Yes                          = "y"
	PathSep                      = string(os.PathSeparator)
	LayoutDateTime               = "2006-01-02 15:04:05"
	LayoutYearMonthDay           = "2006-01-02"
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
	return homeDir() + PathSep + "zodo-data"
}

func Path(filename string) string {
	return ProjectDir() + PathSep + filename
}
