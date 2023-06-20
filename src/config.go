package zodo

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	StorageTypeRedis = "redis"
	StorageTypeFile  = "file"
)

const (
	defaultPadding               = 2
	defaultTableMaxLength        = 200
	defaultPendingStatusColor    = ColorMagenta
	defaultProcessingStatusColor = ColorCyan
	defaultDoneStatusColor       = ColorBlue
	defaultHidingStatusColor     = ColorBlack
	defaultNormalDeadlineColor   = ColorGreen
	defaultNervousDeadlineColor  = ColorYellow
	defaultOverdueDeadlineColor  = ColorRed
	defaultPollingIntervalSecond = 1
)

const configFileName = "conf"

var Config config

type config struct {
	Todo struct {
		Padding        int `yaml:"padding"`
		TableMaxLength int `yaml:"tableMaxLength"`
		Color          struct {
			Status struct {
				Pending    string `yaml:"pending"`
				Processing string `yaml:"processing"`
				Done       string `yaml:"done"`
				Hiding     string `yaml:"hiding"`
			} `yaml:"status"`
			Deadline struct {
				Normal  string `yaml:"normal"`
				Nervous string `yaml:"nervous"`
				Overdue string `yaml:"overdue"`
			} `yaml:"deadline"`
		} `yaml:"color"`
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
	Jenkins struct {
		Url                   string `yaml:"url"`
		Username              string `yaml:"username"`
		Password              string `yaml:"password"`
		PrintStatus           bool   `yaml:"printStatus"`
		PollingIntervalSecond int    `yaml:"pollingIntervalSecond"`
	} `yaml:"jenkins"`
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
		Config.check()
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
	c.check()
}

func (c *config) check() {
	if c.Todo.Padding <= 0 {
		c.Todo.Padding = defaultPadding
	}
	if c.Todo.TableMaxLength <= 0 {
		c.Todo.TableMaxLength = defaultTableMaxLength
	}
	if c.Todo.Color.Status.Pending == "" {
		c.Todo.Color.Status.Pending = defaultPendingStatusColor
	}
	if c.Todo.Color.Status.Processing == "" {
		c.Todo.Color.Status.Processing = defaultProcessingStatusColor
	}
	if c.Todo.Color.Status.Done == "" {
		c.Todo.Color.Status.Done = defaultDoneStatusColor
	}
	if c.Todo.Color.Status.Hiding == "" {
		c.Todo.Color.Status.Hiding = defaultHidingStatusColor
	}
	if c.Todo.Color.Deadline.Normal == "" {
		c.Todo.Color.Deadline.Normal = defaultNormalDeadlineColor
	}
	if c.Todo.Color.Deadline.Nervous == "" {
		c.Todo.Color.Deadline.Nervous = defaultNervousDeadlineColor
	}
	if c.Todo.Color.Deadline.Overdue == "" {
		c.Todo.Color.Deadline.Overdue = defaultOverdueDeadlineColor
	}

	if c.Storage.Type == "" {
		c.Storage.Type = StorageTypeFile
	}

	if c.Jenkins.PollingIntervalSecond <= 0 {
		c.Jenkins.PollingIntervalSecond = defaultPollingIntervalSecond
	}
}
