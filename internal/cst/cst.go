package cst

import "os"

const (
	PathSep        = string(os.PathSeparator)
	LayoutDateTime = "2006-01-02 15:04:05"
	LayoutDate     = "2006-01-02"
	Placeholder    = "-"
)

func HomeDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return homeDir
}
