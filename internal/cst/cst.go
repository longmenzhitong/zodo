package cst

import "os"

const (
	PathSep                      = string(os.PathSeparator)
	LayoutDateTime               = "2006-01-02 15:04:05"
	LayoutYearMonthDay           = "2006-01-02"
	LayoutMonthDay               = "01-02"
	LayoutHourMinute             = "15:04"
	LayoutMonthDayHourMinute     = "01-02 15:04"
	LayoutYearMonthDayHourMinute = "2006-01-02 15:04"
)

func HomeDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return homeDir
}
