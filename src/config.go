package zodo

import (
	"errors"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

const (
	StorageTypeRedis = "redis"
	StorageTypeFile  = "file"
)

type config struct {
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
		MaxLen  int `yaml:"maxLen"`
		Padding int `yaml:"padding"`
	} `yaml:"table"`
}

const configFileName = "conf"

var (
	configPath string
	Config     config
)

func InitConfig() {
	dir := ProjectDir()
	if _, err := os.Stat(dir); errors.Is(err, os.ErrNotExist) {
		err = os.Mkdir(dir, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	configPath = Path(configFileName)

	if _, err := os.Stat(configPath); errors.Is(err, os.ErrNotExist) {
		initYaml(configPath)
	}

	parseYaml(configPath)
}

func initYaml(path string) {
	RewriteLinesToPath(path, []string{
		"storage:",
		"  type: file",
		"  redis:",
		"    address:",
		"    password:",
		"    db:",
		"    localize:",
		"dailyReport:",
		"  enabled: false",
		"  cron:",
		"reminder:",
		"  enabled: false",
		"  cron:",
		"email:",
		"  server:",
		"  port:",
		"  auth:",
		"  from:",
		"  to:",
		"table:",
		"  maxLen: 200",
		"  padding: 3",
	})
}

func parseYaml(path string) {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(f, &Config)
	if err != nil {
		panic(err)
	}
	if Config.Table.MaxLen == 0 {
		Config.Table.MaxLen = 150
	}
}
