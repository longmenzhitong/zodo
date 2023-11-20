package zodo

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	defaultPadding        = 2
	defaultTableMaxLength = 200
)

const (
	defaultColorStatusPending    = ColorMagenta
	defaultColorStatusProcessing = ColorCyan
	defaultColorStatusDone       = ColorBlue
	defaultColorStatusHiding     = ColorBlack
	defaultColorDeadlineNormal   = ColorGreen
	defaultColorDeadlineNervous  = ColorYellow
	defaultColorDeadlineOverdue  = ColorRed
)

const defaultEditor = "vim"

const SyncTypeRedis = "redis"

const configFileName = "config.yml"

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
		ShowDone         bool   `yaml:"showDone"`
		ShowParentStatus bool   `yaml:"showParentStatus"`
		CopyIdAfterAdd   bool   `yaml:"copyIdAfterAdd"`
		Editor           string `yaml:"editor"`
	} `yaml:"todo"`
	Sync struct {
		Type  string `yaml:"type"`
		Redis struct {
			Address  string `yaml:"address"`
			Password string `yaml:"password"`
			Db       int    `yaml:"db"`
		} `yaml:"redis"`
	} `yaml:"sync"`
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
		c.Todo.Color.Status.Pending = defaultColorStatusPending
	}
	if c.Todo.Color.Status.Processing == "" {
		c.Todo.Color.Status.Processing = defaultColorStatusProcessing
	}
	if c.Todo.Color.Status.Done == "" {
		c.Todo.Color.Status.Done = defaultColorStatusDone
	}
	if c.Todo.Color.Status.Hiding == "" {
		c.Todo.Color.Status.Hiding = defaultColorStatusHiding
	}
	if c.Todo.Color.Deadline.Normal == "" {
		c.Todo.Color.Deadline.Normal = defaultColorDeadlineNormal
	}
	if c.Todo.Color.Deadline.Nervous == "" {
		c.Todo.Color.Deadline.Nervous = defaultColorDeadlineNervous
	}
	if c.Todo.Color.Deadline.Overdue == "" {
		c.Todo.Color.Deadline.Overdue = defaultColorDeadlineOverdue
	}

	if c.Todo.Editor == "" {
		c.Todo.Editor = defaultEditor
	}

}
