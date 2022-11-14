package zodo

import (
	"fmt"
	"github.com/atotto/clipboard"
	"gopkg.in/yaml.v3"
	"strconv"
)

type Option struct {
	List           ListCommand           `command:"ls" description:"Show todo list: list [-a] <keyword>"`
	Detail         DetailCommand         `command:"cat" description:"Show todo detail: cat <id>..."`
	Add            AddCommand            `command:"add" description:"Add todo: add [-p <parent-id>] [-d <deadline>] <content>"`
	Modify         ModifyCommand         `command:"mod" description:"Modify todo: mod <id> <content>"`
	Remove         RemoveCommand         `command:"rm" description:"Remove todo: rm <id>..."`
	SetRemark      SetRemarkCommand      `command:"rmk" description:"Set remark of todo: rmk <id> <remark>"`
	SetDeadline    SetDeadlineCommand    `command:"ddl" description:"Set deadline of todo: ddl <id> <deadline>"`
	RemoveDeadline RemoveDeadlineCommand `command:"ddl-" description:"Remove deadline of todo: ddl- <id>..."`
	SetRemind      SetRemindCommand      `command:"rmd" description:"Set remind of todo: rmd [-l] <id> <remind-time>"`
	RemoveRemind   RemoveRemindCommand   `command:"rmd-" description:"Remove remind of todo: rmd- <id>..."`
	SetChild       SetChildCommand       `command:"scd" description:"Set child of todo: scd <parent-id> <child-id>..."`
	AddChild       AddChildCommand       `command:"acd" description:"Add child of todo: acd <parent-id> <child-id>..."`
	SetPending     SetPendingCommand     `command:"pend" description:"Mark todo status as pending: pend <id>..."`
	SetProcessing  SetProcessingCommand  `command:"proc" description:"Mark todo status as processing: proc <id>..."`
	SetDone        SetDoneCommand        `command:"done" description:"Mark todo status as done: done <id>..."`
	SetHiding      SetHidingCommand      `command:"hide" description:"Mark todo status as hiding: hide <id>..."`
	Server         ServerCommand         `command:"server" description:"Enter server mode"`
	Report         ReportCommand         `command:"report" description:"Send report email"`
	Rollback       RollbackCommand       `command:"rollback" description:"Rollback to last version"`
	Transfer       TransferCommand       `command:"transfer" description:"Transfer between file and redis"`
	Clear          ClearCommand          `command:"clr" description:"Clear done todos"`
	Config         ConfigCommand         `command:"conf" description:"Show configs"`
}

type ListCommand struct {
	AllStatus bool     `short:"a" required:"false" description:"Show all status todos"`
	Status    []string `short:"s" required:"false" description:"Search by status prefix"`
}

func (c *ListCommand) Execute(args []string) error {
	List(argsToStr(args), c.Status, c.AllStatus)
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
		err = Detail(id)
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
	id, err := Add(argsToStr(args))
	if err != nil {
		return err
	}

	if c.ParentId != 0 {
		err = SetChild(c.ParentId, []int{id}, true)
		if err != nil {
			return err
		}
	}

	if c.Deadline != "" {
		ddl, err := validateDeadline(c.Deadline)
		if err != nil {
			return err
		}
		SetDeadline(id, ddl)
	}

	if c.Remind != "" {
		rmd, err := validateRemind(c.Remind)
		if err != nil {
			return err
		}
		err = SetRemind(id, rmd, true)
		if err != nil {
			return err
		}
	}

	Save()

	if Config.Todo.CopyIdAfterAdd {
		err = clipboard.WriteAll(strconv.Itoa(id))
		if err != nil {
			return err
		} else {
			fmt.Println("Id copied.")
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
	Modify(id, content)
	Save()
	return nil
}

type RemoveCommand struct {
}

func (c *RemoveCommand) Execute(args []string) error {
	ids, err := argsToIds(args)
	if err != nil {
		return err
	}
	Remove(ids)
	Save()
	return nil
}

type SetRemarkCommand struct {
}

func (c *SetRemarkCommand) Execute(args []string) error {
	id, remark, err := argsToIdAndStr(args)
	if err != nil {
		return err
	}
	SetRemark(id, remark)
	Save()
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

	SetDeadline(id, ddl)
	Save()
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
		SetDeadline(id, "")
	}
	Save()
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
	err = SetRemind(id, rmd, c.Loop)
	if err != nil {
		return err
	}
	Save()
	return nil
}

type RemoveRemindCommand struct {
}

func (c *RemoveRemindCommand) Execute(args []string) error {
	ids, err := argsToIds(args)
	if err != nil {
		return err
	}
	RemoveRemind(ids)
	Save()
	return nil
}

type SetChildCommand struct {
}

func (c *SetChildCommand) Execute(args []string) error {
	err := setChild(args, false)
	if err != nil {
		return err
	}
	Save()
	return nil
}

type AddChildCommand struct {
}

func (c *AddChildCommand) Execute(args []string) error {
	err := setChild(args, true)
	if err != nil {
		return err
	}
	Save()
	return nil
}

func setChild(args []string, append bool) error {
	ids, err := argsToIds(args)
	if err != nil {
		return err
	}
	if len(ids) < 2 {
		return &InvalidInputError{
			Message: fmt.Sprintf("expect: scd [parentId] [childId], got: %v", args),
		}
	}
	err = SetChild(ids[0], ids[1:], append)
	if err != nil {
		return err
	}
	Save()
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
		SetPending(id)
	}
	Save()
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
		SetProcessing(id)
	}
	Save()
	return nil
}

type SetDoneCommand struct {
}

func (c *SetDoneCommand) Execute(args []string) error {
	ids, err := argsToIds(args)
	if err == nil {
		for _, id := range ids {
			SetDone(id)
		}
		Save()
		return nil
	}

	id, remark, err := argsToIdAndStr(args)
	if err == nil {
		SetDone(id)
		SetRemark(id, remark)
		Save()
		return nil
	}

	return &InvalidInputError{
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
		SetHiding(id)
	}
	Save()
	return nil
}

type ServerCommand struct {
}

func (c *ServerCommand) Execute([]string) error {
	if Config.DailyReport.Enabled {
		StartDailyReport()
	}
	if Config.Reminder.Enabled {
		StartReminder()
	}
	select {}
}

type ReportCommand struct {
}

func (c *ReportCommand) Execute([]string) error {
	return Report()
}

type RollbackCommand struct {
}

func (c *RollbackCommand) Execute([]string) error {
	Rollback()
	return nil
}

type TransferCommand struct {
}

func (c *TransferCommand) Execute([]string) error {
	Transfer()
	return nil
}

type ClearCommand struct {
}

func (c *ClearCommand) Execute([]string) error {
	count := Clear()
	if count > 0 {
		Save()
	}
	fmt.Printf("%d cleared.\n", count)
	return nil
}

type ConfigCommand struct {
}

func (c *ConfigCommand) Execute([]string) error {
	out, err := yaml.Marshal(Config)
	if err != nil {
		return err
	}
	fmt.Println(string(out))
	return nil
}
