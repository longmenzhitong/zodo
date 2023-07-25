package cmd

import (
	"fmt"
	"strconv"
	zodo "zodo/src"
	"zodo/src/dev"
	"zodo/src/dev/jenkins"
	"zodo/src/todo"

	"gopkg.in/yaml.v3"
)

type Option struct {
	Join             JoinCommand             `command:"join" description:"Join todos: join <to-id> <from-id>"`
	Remove           RemoveCommand           `command:"rm" description:"Remove todos: rm [-r] <id>..."`
	SetRemark        SetRemarkCommand        `command:"rmk" description:"Set remark of todo: rmk <id> <remark>"`
	SetDeadline      SetDeadlineCommand      `command:"ddl" description:"Set deadline of todo: ddl <id> <deadline>"`
	RemoveDeadline   RemoveDeadlineCommand   `command:"ddl-" description:"Remove deadline of todo: ddl- <id>..."`
	SetRemind        SetRemindCommand        `command:"rmd" description:"Set remind of todo: rmd [-l] <id> <remind-time>"`
	RemoveRemind     RemoveRemindCommand     `command:"rmd-" description:"Remove remind of todo: rmd- <id>..."`
	SetChild         SetChildCommand         `command:"scd" description:"Set child of todo: scd <parent-id> <child-id>..."`
	AddChild         AddChildCommand         `command:"acd" description:"Add child of todo: acd <parent-id> <child-id>..."`
	SetPending       SetPendingCommand       `command:"pend" description:"Mark todo status as pending: pend <id>..."`
	SetProcessing    SetProcessingCommand    `command:"proc" description:"Mark todo status as processing: proc <id>..."`
	SetDone          SetDoneCommand          `command:"done" description:"Mark todo status as done: done <id>..."`
	SetHiding        SetHidingCommand        `command:"hide" description:"Mark todo status as hiding: hide <id>..."`
	Server           ServerCommand           `command:"server" description:"Server mode on"`
	Rollback         RollbackCommand         `command:"rbk" description:"Rollback to last version"`
	Tidy             TidyCommand             `command:"tidy" description:"Tidy data: tidy [-a] [-d] [-i]"`
	Config           ConfigCommand           `command:"conf" description:"Show config"`
	Statistics       StatisticsCommand       `command:"stat" description:"Show statistics of todos"`
	JenkinsDeploy    JenkinsDeployCommand    `command:"jd" description:"Jenkins deploy: jd"`
	JenkinsStatus    JenkinsStatusCommand    `command:"js" description:"Jenkins status: js"`
	JenkinsHistory   JenkinsHistoryCommand   `command:"jh" description:"Jenkins history: jh [-c <history-count>]"`
	DrawioHelper     DrawioHelperCommand     `command:"dh" description:"Drawio Helper: simplify sql for Drawio import: dh <sql-file-path>"`
	MybatisGenerator MybatisGeneratorCommand `command:"mg" description:"MyBatis Generator: generate result map and column: mg <java-file-path>"`
	ExcelHelper      ExcelHelperCommand      `command:"eh" description:"Excel helper: generate java class from excel template: eh -p <excel-template-path> [-n <java-class-name>] [-i <sheet-index>]"`
}

type JoinCommand struct {
}

func (c *JoinCommand) Execute(args []string) error {
	ids, err := argsToIds(args)
	if err != nil {
		return err
	}

	if len(ids) != 2 {
		return &zodo.InvalidInputError{
			Message: "there must be two ids",
		}
	}
	todo.Join(ids[0], ids[1])
	todo.Save()
	return nil
}

type RemoveCommand struct {
	Recursively bool `short:"r" required:"false" description:"Remove child todo recursively"`
}

func (c *RemoveCommand) Execute(args []string) error {
	ids, err := argsToIds(args)
	if err != nil {
		return err
	}
	todo.Remove(ids, c.Recursively)
	todo.Save()
	return nil
}

type SetRemarkCommand struct {
}

func (c *SetRemarkCommand) Execute(args []string) error {
	if len(args) == 1 {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}

		return todo.CopyRemark(id)
	}

	id, remark, err := argsToIdAndStr(args)
	if err != nil {
		return err
	}
	todo.SetRemark(id, remark)
	todo.Save()
	return nil
}

type SetDeadlineCommand struct {
}

func (c *SetDeadlineCommand) Execute(args []string) error {
	id, ddl, err := argsToIdAndStr(args)
	if err != nil {
		return err
	}

	ddl, err = validateDeadline(ddl)
	if err != nil {
		return err
	}

	todo.SetDeadline(id, ddl)
	todo.Save()
	return nil
}

type RemoveDeadlineCommand struct {
}

func (c *RemoveDeadlineCommand) Execute(args []string) error {
	ids, err := argsToIds(args)
	if err != nil {
		return err
	}
	for _, id := range ids {
		todo.SetDeadline(id, "")
	}
	todo.Save()
	return nil
}

type SetRemindCommand struct {
	Loop bool `short:"l" required:"false" description:"Choose loop type"`
}

func (c *SetRemindCommand) Execute(args []string) error {
	id, rmd, err := argsToIdAndStr(args)
	if err != nil {
		return err
	}
	rmd, err = validateRemind(rmd)
	if err != nil {
		return err
	}
	err = todo.SetRemind(id, rmd, c.Loop)
	if err != nil {
		return err
	}
	todo.Save()
	return nil
}

type RemoveRemindCommand struct {
}

func (c *RemoveRemindCommand) Execute(args []string) error {
	ids, err := argsToIds(args)
	if err != nil {
		return err
	}
	todo.RemoveRemind(ids)
	todo.Save()
	return nil
}

type SetChildCommand struct {
}

func (c *SetChildCommand) Execute(args []string) error {
	return setChild(args, false)
}

type AddChildCommand struct {
}

func (c *AddChildCommand) Execute(args []string) error {
	return setChild(args, true)
}

func setChild(args []string, append bool) error {
	ids, err := argsToIds(args)
	if err != nil {
		return err
	}
	if len(ids) < 2 {
		return &zodo.InvalidInputError{
			Message: fmt.Sprintf("expect: scd [parentId] [childId], got: %v", args),
		}
	}
	err = todo.SetChild(ids[0], ids[1:], append)
	if err != nil {
		return err
	}
	todo.Save()
	return nil
}

type SetPendingCommand struct {
}

func (c *SetPendingCommand) Execute(args []string) error {
	ids, err := argsToIds(args)
	if err != nil {
		return err
	}
	for _, id := range ids {
		todo.SetPending(id)
	}
	todo.Save()
	return nil
}

type SetProcessingCommand struct {
}

func (c *SetProcessingCommand) Execute(args []string) error {
	ids, err := argsToIds(args)
	if err != nil {
		return err
	}
	for _, id := range ids {
		todo.SetProcessing(id)
	}
	todo.Save()
	return nil
}

type SetDoneCommand struct {
}

func (c *SetDoneCommand) Execute(args []string) error {
	ids, err := argsToIds(args)
	if err == nil {
		for _, id := range ids {
			todo.SetDone(id)
		}
		todo.Save()
		return nil
	}

	id, remark, err := argsToIdAndStr(args)
	if err == nil {
		todo.SetDone(id)
		todo.SetRemark(id, remark)
		todo.Save()
		return nil
	}

	return &zodo.InvalidInputError{
		Message: fmt.Sprintf("expect: done [id1] [id2]... or done [id] [remark], got: %v", args),
	}
}

type SetHidingCommand struct {
}

func (c *SetHidingCommand) Execute(args []string) error {
	ids, err := argsToIds(args)
	if err != nil {
		return err
	}
	for _, id := range ids {
		todo.SetHiding(id)
	}
	todo.Save()
	return nil
}

type ServerCommand struct {
}

func (c *ServerCommand) Execute([]string) error {
	if zodo.Config.DailyReport.Enabled {
		todo.StartDailyReport()
	}
	if zodo.Config.Reminder.Enabled {
		todo.StartReminder()
	}
	select {}
}

type RollbackCommand struct {
}

func (c *RollbackCommand) Execute([]string) error {
	todo.Rollback()
	return nil
}

type TidyCommand struct {
	All      bool `short:"a" required:"false" description:"Execute all tidy works"`
	DoneTodo bool `short:"d" required:"false" description:"Clear done todos"`
	Id       bool `short:"i" required:"false" description:"Defrag ids"`
}

func (c *TidyCommand) Execute([]string) error {
	changed := false
	if c.All || c.DoneTodo {
		count := todo.ClearDoneTodo()
		if count > 0 {
			zodo.PrintDoneMsg("Clear %d done todos.\n", count)
			changed = true
		}
	}
	if c.All || c.Id {
		from, to := todo.DefragId()
		if from != to {
			zodo.PrintDoneMsg("Defrag ids from %d to %d.\n", from, to)
			changed = true
		}
	}
	if changed {
		todo.Save()
	}
	return nil
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
