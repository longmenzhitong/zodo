package command

import (
	"fmt"
	"zodo/internal/conf"
	"zodo/internal/errs"
	"zodo/internal/task"
	"zodo/internal/todos"
)

type Option struct {
	Server         ServerCommand         `command:"sv" description:"Enter server mode"`
	List           ListCommand           `command:"ls" description:"Show todo list: list [-a] <keyword>"`
	Detail         DetailCommand         `command:"cat" description:"Show todo detail: cat <id>..."`
	Add            AddCommand            `command:"add" description:"Add todo: add [-p <parent-id>] [-d <deadline>] <content>"`
	Modify         ModifyCommand         `command:"mod" description:"Modify todo: mod <id> <content>"`
	Remove         RemoveCommand         `command:"rm" description:"Remove todo: rm <id>..."`
	DailyReport    DailyReportCommand    `command:"dr" description:"Send daily report email"`
	Rollback       RollbackCommand       `command:"rbk" description:"Rollback to last version"`
	Transfer       TransferCommand       `command:"tr" description:"Transfer between file and redis"`
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
}

type ServerCommand struct {
}

func (c *ServerCommand) Execute([]string) error {
	if conf.Data.DailyReport.Enabled {
		task.StartDailyReport()
	}
	if conf.Data.Reminder.Enabled {
		task.StartReminder()
	}
	select {}
}

type ListCommand struct {
	All bool `short:"a" required:"false" description:"Show all todos"`
}

func (c *ListCommand) Execute(args []string) error {
	keyword := argsToStr(args)
	todos.List(keyword, c.All)
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
		todos.Detail(id)
		fmt.Println()
	}
	return nil
}

type AddCommand struct {
	ParentId int    `short:"p" required:"false" description:"parent-id"`
	Deadline string `short:"d" required:"false" description:"deadline"`
}

func (c *AddCommand) Execute(args []string) error {
	id, err := todos.Add(argsToStr(args))
	if err != nil {
		return err
	}

	if c.ParentId != 0 {
		err = todos.SetChild(c.ParentId, []int{id}, true)
		if err != nil {
			return err
		}
	}

	if c.Deadline != "" {
		ddl, err := validateDeadline(c.Deadline)
		if err != nil {
			return err
		}
		todos.SetDeadline(id, ddl)
	}

	todos.Save()
	return nil
}

type ModifyCommand struct {
}

func (c *ModifyCommand) Execute(args []string) error {
	id, content, err := argsToIdAndStr(args)
	if err != nil {
		return err
	}
	todos.Modify(id, content)
	todos.Save()
	return nil
}

type RemoveCommand struct {
}

func (c *RemoveCommand) Execute(args []string) error {
	ids, err := argsToIds(args)
	if err != nil {
		return err
	}
	todos.Delete(ids)
	todos.Save()
	return nil
}

type DailyReportCommand struct {
}

func (c *DailyReportCommand) Execute([]string) error {
	return todos.DailyReport()
}

type RollbackCommand struct {
}

func (c *RollbackCommand) Execute([]string) error {
	todos.Rollback()
	return nil
}

type TransferCommand struct {
}

func (c *TransferCommand) Execute([]string) error {
	todos.Transfer()
	return nil
}

type SetRemarkCommand struct {
}

func (c *SetRemarkCommand) Execute(args []string) error {
	id, remark, err := argsToIdAndStr(args)
	if err != nil {
		return err
	}
	todos.SetRemark(id, remark)
	todos.Save()
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

	todos.SetDeadline(id, ddl)
	todos.Save()
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
		todos.SetDeadline(id, "")
	}
	todos.Save()
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
	err = todos.SetRemind(id, rmd, c.Loop)
	if err != nil {
		return err
	}
	todos.Save()
	return nil
}

type RemoveRemindCommand struct {
}

func (c *RemoveRemindCommand) Execute(args []string) error {
	ids, err := argsToIds(args)
	if err != nil {
		return err
	}
	todos.RemoveRemind(ids)
	todos.Save()
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
		return &errs.InvalidInputError{
			Message: fmt.Sprintf("expect: scd [parentId] [childId], got: %v", args),
		}
	}
	err = todos.SetChild(ids[0], ids[1:], append)
	if err != nil {
		return err
	}
	todos.Save()
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
		todos.SetPending(id)
	}
	todos.Save()
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
		todos.SetProcessing(id)
	}
	todos.Save()
	return nil
}

type SetDoneCommand struct {
}

func (c *SetDoneCommand) Execute(args []string) error {
	ids, err := argsToIds(args)
	if err == nil {
		for _, id := range ids {
			todos.SetDone(id)
		}
		todos.Save()
		return nil
	}

	id, remark, err := argsToIdAndStr(args)
	if err == nil {
		todos.SetDone(id)
		todos.SetRemark(id, remark)
		todos.Save()
		return nil
	}

	return &errs.InvalidInputError{
		Message: fmt.Sprintf("expect: done [id1] [id2]... or done [id] [remark], got: %v", args),
	}
}
