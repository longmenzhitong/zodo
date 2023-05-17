package zodo

import (
	"errors"
	"gopkg.in/yaml.v3"
	"os"
)

const (
	StorageTypeRedis = "redis"
	StorageTypeFile  = "file"
)

const configFileName = "conf"

var Config config

type config struct {
	Todo struct {
		Padding          int  `yaml:"padding"`
		ShowDone         bool `yaml:"showDone"`
		ShowParentStatus bool `yaml:"showParentStatus"`
		CopyIdAfterAdd   bool `yaml:"copyIdAfterAdd"`
	} `yaml:"todo"`
	Storage struct {
		Type  string `yaml:"type"`
		Redis struct {
			Address  string `yaml:"address"`
			Password string `yaml:"password"`
			Db       int    `yaml:"db"`
			Localize bool   `yaml:"localize"`
		} `yaml:"redis"`
	} `yaml:"storage"`
	DailyReport struct {
		Enabled bool   `yaml:"enabled"`
		Cron    string `yaml:"cron"`
	} `yaml:"dailyReport"`
	Reminder struct {
		Enabled bool   `yaml:"enabled"`
		Cron    string `yaml:"cron"`
	} `yaml:"reminder"`
	Email struct {
		Server string   `yaml:"server"`
		Port   int      `yaml:"port"`
		Auth   string   `yaml:"auth"`
		From   string   `yaml:"from"`
		To     []string `yaml:"to"`
	} `yaml:"email"`
	Table struct {
		MaxLen int `yaml:"maxLen"`
	} `yaml:"table"`
	Jenkins struct {
		Url      string `yaml:"url"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	}
}

func (c *config) Init() {
	dir := ProjectDir()
	if _, err := os.Stat(dir); errors.Is(err, os.ErrNotExist) {
		err = os.Mkdir(dir, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	configPath := Path(configFileName)
	if _, err := os.Stat(configPath); errors.Is(err, os.ErrNotExist) {
		Config = config{}
		Config.Todo.Padding = 3
		Config.Todo.ShowDone = false
		Config.Todo.ShowParentStatus = false
		Config.Todo.CopyIdAfterAdd = true
		Config.Storage.Type = StorageTypeFile
		Config.DailyReport.Enabled = false
		Config.Reminder.Enabled = false
		Config.Table.MaxLen = 200
		out, err := yaml.Marshal(Config)
		if err != nil {
			panic(err)
		}
		RewriteLinesToPath(configPath, []string{string(out)})
		return
	}

	f, err := os.ReadFile(configPath)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(f, &c)
	if err != nil {
		panic(err)
	}
}
