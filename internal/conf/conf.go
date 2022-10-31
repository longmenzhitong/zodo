package conf

import (
	"errors"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"zodo/internal/files"
)

const (
	StorageTypeRedis = "redis"
	StorageTypeFile  = "file"
)

type data struct {
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
	files.RewriteLinesToPath(path, []string{
		"storage:",
		"  type:",
		"  redis:",
		"    address:",
		"    password:",
		"    db:",
		"    localize:",
		"dailyReport:",
		"  enabled:",
		"  cron:",
		"email:",
		"  server:",
		"  port:",
		"  auth:",
		"  from:",
		"  to:",
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

func IsFileStorage(storageType ...string) bool {
	if storageType != nil && len(storageType) > 0 {
		return storageType[0] == StorageTypeFile
	} else {
		return Data.Storage.Type == StorageTypeFile
	}
}

func IsRedisStorage(storageType ...string) bool {
	if storageType != nil && len(storageType) > 0 {
		return storageType[0] == StorageTypeRedis
	} else {
		return Data.Storage.Type == StorageTypeRedis
	}
}
