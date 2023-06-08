package cmd

import (
	"fmt"
	"github.com/atotto/clipboard"
	"gopkg.in/yaml.v3"
	"strconv"
	"zodo/src"
	"zodo/src/todo"
)

type Option struct {
	List             ListCommand             `command:"ls" description:"Show todo list: list [-a] [-s <status-prefix>] <keyword>"`
	Detail           DetailCommand           `command:"cat" description:"Show todo detail: cat <id>..."`
	Add              AddCommand              `command:"add" description:"Add todo: add [-p <parent-id>] [-d <deadline>] [-r <remind-time>] <content>"`
	Modify           ModifyCommand           `command:"mod" description:"Modify todo: mod <id> <content>"`
	Remove           RemoveCommand           `command:"rm" description:"Remove todo: rm <id>..."`
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
	Server           ServerCommand           `command:"server" description:"Enter server mode"`
	Report           ReportCommand           `command:"report" description:"Send report email"`
	Rollback         RollbackCommand         `command:"rbk" description:"Rollback to last version"`
	Transfer         TransferCommand         `command:"trans" description:"Transfer between file and redis"`
	Tidy             TidyCommand             `command:"tidy" description:"Tidy data: tidy [-a] [-d] [-i]"`
	Config           ConfigCommand           `command:"conf" description:"Show config"`
	Info             InfoCommand             `command:"info" description:"Show info"`
	SimplifySql      SimplifySqlCommand      `command:"ss" description:"Simplify sql for drawio"`
	Tea              TeaCommand              `command:"tea" description:"Wait for a tea: tea <minutes-to-wait>"`
	Jenkins          JenkinsCommand          `command:"jk" description:"Deploy by the Jenkins: jk [-s <service>] [-e <env>] [-b <branch>] [-c] [-S]"`
	MybatisGenerator MybatisGeneratorCommand `command:"mbg" description:"Generate MyBatis code: mbg -p <path>"`
}

type ListCommand struct {
	AllStatus bool     `short:"a" required:"false" description:"Show all status todos"`
	Status    []string `short:"s" required:"false" description:"Search by status prefix"`
}

func (c *ListCommand) Execute(args []string) error {
	todo.List(argsToStr(args), c.Status, c.AllStatus)
	return nil
}

type DetailCommand struct {
}

func (c *DetailCommand) Execute(args []string) error {
	ids, err := argsToIds(args)
	if err != nil {
		return err
	}
	for _, id := range ids {
		err = todo.Detail(id)
		if err != nil {
			return err
		}
		fmt.Println()
	}
	return nil
}

type AddCommand struct {
	ParentId int    `short:"p" required:"false" description:"parent-id"`
	Deadline string `short:"d" required:"false" description:"deadline"`
	Remind   string `short:"r" required:"false" description:"remind-time"`
}

func (c *AddCommand) Execute(args []string) error {
	id, err := todo.Add(argsToStr(args))
	if err != nil {
		return err
	}

	if c.ParentId != 0 {
		err = todo.SetChild(c.ParentId, []int{id}, true)
		if err != nil {
			return err
		}
	}

	if c.Deadline != "" {
		ddl, err := validateDeadline(c.Deadline)
		if err != nil {
			return err
		}
		todo.SetDeadline(id, ddl)
	}

	if c.Remind != "" {
		rmd, err := validateRemind(c.Remind)
		if err != nil {
			return err
		}
		err = todo.SetRemind(id, rmd, true)
		if err != nil {
			return err
		}
	}

	todo.Save()

	if zodo.Config.Todo.CopyIdAfterAdd {
		err = clipboard.WriteAll(strconv.Itoa(id))
		if err != nil {
			return err
		}
	}
	return nil
}

type ModifyCommand struct {
}

func (c *ModifyCommand) Execute(args []string) error {
	id, content, err := argsToIdAndStr(args)
	if err != nil {
		return err
	}
	todo.Modify(id, content)
	todo.Save()
	return nil
}

type RemoveCommand struct {
}

func (c *RemoveCommand) Execute(args []string) error {
	ids, err := argsToIds(args)
	if err != nil {
		return err
	}
	todo.Remove(ids)
	todo.Save()
	return nil
}

type SetRemarkCommand struct {
}

func (c *SetRemarkCommand) Execute(args []string) error {
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

type ReportCommand struct {
}

func (c *ReportCommand) Execute([]string) error {
	return todo.Report()
}

type RollbackCommand struct {
}

func (c *RollbackCommand) Execute([]string) error {
	todo.Rollback()
	return nil
}

type TransferCommand struct {
}

func (c *TransferCommand) Execute([]string) error {
	todo.Transfer()
	return nil
}

type TidyCommand struct {
	All      bool `short:"a" required:"false" description:"All tidy works"`
	DoneTodo bool `short:"d" required:"false" description:"Clear done todos"`
	Id       bool `short:"i" required:"false" description:"Defrag ids"`
}

func (c *TidyCommand) Execute([]string) error {
	count := 0
	if c.All || c.DoneTodo {
		count += todo.ClearDoneTodo()
	}
	if c.All || c.Id {
		count += todo.DefragId()
	}
	if count > 0 {
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

type InfoCommand struct {
}

func (c *InfoCommand) Execute([]string) error {
	todo.Info()
	return nil
}

type SimplifySqlCommand struct {
}

func (c *SimplifySqlCommand) Execute(args []string) error {
	path := argsToStr(args)
	zodo.SimplifySql(path)
	return nil
}

type TeaCommand struct {
}

func (c *TeaCommand) Execute(args []string) error {
	minutes, err := strconv.Atoi(argsToStr(args))
	if err != nil {
		return err
	}
	err = todo.Tea(minutes)
	if err != nil {
		return err
	}
	todo.Save()
	return nil
}

type JenkinsCommand struct {
	Service    string `short:"s" required:"false" description:"Service name or current dir name by default"`
	Env        string `short:"e" required:"false" description:"Service environment"`
	Branch     string `short:"b" required:"false" description:"Git branch or current git branch by default"`
	CheckCode  bool   `short:"c" required:"false" description:"Check code option"`
	StatusOnly bool   `short:"S" required:"false" description:"Print build status only"`
}

func (c *JenkinsCommand) Execute([]string) error {
	return zodo.Deploy(c.Service, c.Env, c.Branch, c.CheckCode, c.StatusOnly)
}

type MybatisGeneratorCommand struct {
	Path string `short:"p" required:"true" descriptoin:"POJO path"`
}

func (c *MybatisGeneratorCommand) Execute([]string) error {
	return zodo.GenerateMybatisCode(c.Path)
}
