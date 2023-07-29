package cmd

import (
	"fmt"
	zodo "zodo/src"
	"zodo/src/dev"
	"zodo/src/dev/jenkins"
	"zodo/src/todo"

	"gopkg.in/yaml.v3"
)

type Option struct {
	Config           ConfigCommand           `command:"conf" description:"Show config"`
	Statistics       StatisticsCommand       `command:"stat" description:"Show statistics of todos"`
	JenkinsDeploy    JenkinsDeployCommand    `command:"jd" description:"Jenkins deploy: jd"`
	JenkinsStatus    JenkinsStatusCommand    `command:"js" description:"Jenkins status: js"`
	JenkinsHistory   JenkinsHistoryCommand   `command:"jh" description:"Jenkins history: jh [-c <history-count>]"`
	DrawioHelper     DrawioHelperCommand     `command:"dh" description:"Drawio Helper: simplify sql for Drawio import: dh <sql-file-path>"`
	MybatisGenerator MybatisGeneratorCommand `command:"mg" description:"MyBatis Generator: generate result map and column: mg <java-file-path>"`
	ExcelHelper      ExcelHelperCommand      `command:"eh" description:"Excel helper: generate java class from excel template: eh -p <excel-template-path> [-n <java-class-name>] [-i <sheet-index>]"`
}

type ConfigCommand struct {
}

func (c *ConfigCommand) Execute([]string) error {
	out, err := yaml.Marshal(zodo.Config)
	if err != nil {
		return err
	}
	fmt.Println(string(out))
	return nil
}

type StatisticsCommand struct {
}

func (c *StatisticsCommand) Execute([]string) error {
	todo.Statistics()
	return nil
}

type DrawioHelperCommand struct {
}

func (c *DrawioHelperCommand) Execute(args []string) error {
	path := argsToStr(args)
	dev.SimplifySql(path)
	return nil
}

type JenkinsDeployCommand struct {
}

func (c *JenkinsDeployCommand) Execute([]string) error {
	return jenkins.Deploy()
}

type JenkinsStatusCommand struct {
}

func (c *JenkinsStatusCommand) Execute([]string) error {
	return jenkins.Status()
}

type JenkinsHistoryCommand struct {
	Count int `short:"c" required:"false" description:"History count, default: 5"`
}

func (c *JenkinsHistoryCommand) Execute([]string) error {
	if c.Count <= 0 {
		c.Count = 5
	}
	return jenkins.History(c.Count)
}

type MybatisGeneratorCommand struct {
}

func (c *MybatisGeneratorCommand) Execute(args []string) error {
	return dev.GenerateMybatisCode(argsToStr(args))
}

type ExcelHelperCommand struct {
	Path       string `short:"p" required:"true" description:"Path of excel template"`
	Name       string `short:"n" required:"false" description:"Name of java class, default: ExportDTO"`
	SheetIndex int    `short:"i" required:"false" description:"Index of excel sheet, default: 0"`
}

func (c *ExcelHelperCommand) Execute(args []string) error {
	if c.Name == "" {
		c.Name = "ExportDTO"
	}
	return dev.GenerateJavaCode(c.Path, c.Name, c.SheetIndex)
}
