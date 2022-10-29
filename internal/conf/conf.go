package conf

import (
	"errors"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"zodo/internal/files"
)

type data struct {
	Git struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Email    string `yaml:"email"`
	} `yaml:"git"`
	Reminder struct {
		DailyReport struct {
			Enabled bool   `yaml:"enabled"`
			Cron    string `yaml:"cron"`
		} `yaml:"dailyReport"`
		Email struct {
			Server string   `yaml:"server"`
			Auth   string   `yaml:"auth"`
			From   string   `yaml:"from"`
			To     []string `yaml:"to"`
		} `yaml:"email"`
	} `yaml:"reminder"`
	Table struct {
		MaxLen int `yaml:"maxLen"`
	} `yaml:"table"`
}

const fileName = "conf"

var Data data

var path string

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
		"  username:",
		"  password:",
		"  email:",
		"reminder:",
		"  dailyReport:",
		"    enabled:",
		"    cron:",
		"  email:",
		"    server:",
		"    auth:",
		"    from:",
		"    to:",
		"table:",
		"  maxLen:",
	})
}

func parseYaml(path string) {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(f, &Data)
	if err != nil {
		panic(err)
	}
	if Data.Table.MaxLen == 0 {
		Data.Table.MaxLen = 150
	}
}
