package file

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"zodo/internal/cst"
)

const (
	dirName = ".zodo"
)

var (
	Dir string
)

func init() {
	Dir = cst.HomeDir() + cst.PathSep + dirName
	if _, err := os.Stat(Dir); errors.Is(err, os.ErrNotExist) {
		err = os.Mkdir(Dir, os.ModePerm)
		if err != nil {
			panic(fmt.Errorf("mkdir %s error: %v", Dir, err))
		}
	}
}

func ReadLinesFromPath(path string) []string {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = f.Close()
	}()

	return readLinesFromFile(f)
}

func readLinesFromFile(f *os.File) []string {
	scanner := bufio.NewScanner(f)
	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			lines = append(lines, line)
		}
	}

	if scanner.Err() != nil {
		panic(scanner.Err())
	}

	return lines
}

func RewriteLinesToPath(path string, lines []string) {
	WriteLinesToPath(path, lines, os.O_RDWR|os.O_TRUNC)
}

func WriteLinesToPath(path string, lines []string, mod int) {
	f, err := os.OpenFile(path, mod, 0)
	if err != nil {
		panic(err)
	}
	writeLinesToFile(f, lines)
}

func writeLinesToFile(f *os.File, lines []string) {
	w := bufio.NewWriter(f)
	for _, line := range lines {
		_, err := fmt.Fprintln(w, line)
		if err != nil {
			panic(err)
		}
	}

	err := w.Flush()
	if err != nil {
		panic(err)
	}
}

func EnsureExist(path string) {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		_, err = os.Create(path)
		if err != nil {
			panic(fmt.Errorf("create %s error: %v", path, err))
		}
	}
}
