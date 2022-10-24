package conf

import (
	"errors"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"zodo/internal/files"
)

type config struct {
	Git struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	}
}

const (
	fileName = "conf"
)

var (
	All  config
	path string
)

func init() {
	path = files.GetPath(fileName)

	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		initYaml(path)
		return
	}

	parseYaml(path)
}

func initYaml(path string) {
	files.EnsureExist(path)
	files.RewriteLinesToPath(path, []string{
		"git:",
		"  username: ",
		"  password: ",
	})
}

func parseYaml(path string) {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(f, &All)
	if err != nil {
		panic(err)
	}
}
