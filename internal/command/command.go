package command

import (
	"fmt"
	"zodo/internal/conf"
	"zodo/internal/errs"
	"zodo/internal/task"
	"zodo/internal/todos"
)

type Option struct {
	Server         ServerCommand         `command:"server"`
	List           ListCommand           `command:"ls"`
	Detail         DetailCommand         `command:"cat"`
	Add            AddCommand            `command:"add"`
	Modify         ModifyCommand         `command:"mod"`
	Remove         RemoveCommand         `command:"rm"`
	DailyReport    DailyReportCommand    `command:"dr"`
	Rollback       RollbackCommand       `command:"rbk"`
	Transfer       TransferCommand       `command:"tr"`
	SetRemark      SetRemarkCommand      `command:"rmk"`
	SetDeadline    SetDeadlineCommand    `command:"ddl"`
	RemoveDeadline RemoveDeadlineCommand `command:"ddl-"`
	SetRemind      SetRemindCommand      `command:"rmd"`
	SetLoopRemind  SetLoopRemindCommand  `command:"rmd+"`
	RemoveRemind   RemoveRemindCommand   `command:"rmd-"`
	SetChild       SetChildCommand       `command:"scd"`
	AddChild       AddChildCommand       `command:"acd"`
	SetPending     SetPendingCommand     `command:"pend"`
	SetProcessing  SetProcessingCommand  `command:"proc"`
	SetDone        SetDoneCommand        `command:"done"`
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
	ParentId int    `short:"p" required:"false" description:"parent id"`
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
}

func (c *SetRemindCommand) Execute(args []string) error {
	return setRemind(args, false)
}

type SetLoopRemindCommand struct {
}

func (c *SetLoopRemindCommand) Execute(args []string) error {
	return setRemind(args, true)
}

func setRemind(args []string, loop bool) error {
	id, rmd, err := argsToIdAndStr(args)
	if err != nil {
		return err
	}
	rmd, err = validateRemind(rmd)
	if err != nil {
		return err
	}
	err = todos.SetRemind(id, rmd, loop)
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
