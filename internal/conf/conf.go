package conf

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"zodo/internal/cst"
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
	path = files.Dir + cst.PathSep + fileName

	if _, err := os.Stat(path); err != nil {
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
