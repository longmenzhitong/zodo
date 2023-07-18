package cmd

import (
	"fmt"
	"strconv"
	zodo "zodo/src"
	"zodo/src/dev"
	"zodo/src/dev/jenkins"
	"zodo/src/todo"

	"github.com/atotto/clipboard"
	"gopkg.in/yaml.v3"
)

type Option struct {
	List             ListCommand             `command:"ls" description:"Show todo list: list [-a] [-s <status-prefix>] [<keyword>]"`
	Detail           DetailCommand           `command:"cat" description:"Show todo detail: cat <id>..."`
	Add              AddCommand              `command:"add" description:"Add todo: add [-p <parent-id>] [-d <deadline>] [-r <remind-time>] <content>"`
	Modify           ModifyCommand           `command:"mod" description:"Modify todo: mod <id> <content>"`
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
	Transfer         TransferCommand         `command:"trans" description:"Transfer between file and redis"`
	Tidy             TidyCommand             `command:"tidy" description:"Tidy data: tidy [-a] [-d] [-i]"`
	Config           ConfigCommand           `command:"conf" description:"Show config"`
	Statistics       StatisticsCommand       `command:"stat" description:"Show statistics of todos"`
	JenkinsDeploy    JenkinsDeployCommand    `command:"jd" description:"Jenkins deploy: jd [-j <job-name>] [-s <server-name>] [-b <build-branch>] [-c]"`
	JenkinsStatus    JenkinsStatusCommand    `command:"js" description:"Jenkins status: js [-j <job-name>]"`
	JenkinsHistory   JenkinsHistoryCommand   `command:"jh" description:"Jenkins history: jh [-j <job-name>] [-c <history-count>]"`
	DrawioHelper     DrawioHelperCommand     `command:"dh" description:"Drawio Helper: simplify sql for Drawio import: dh <sql-file-path>"`
	MybatisGenerator MybatisGeneratorCommand `command:"mg" description:"MyBatis Generator: generate result map and column: mg <java-file-path>"`
	ExcelHelper      ExcelHelperCommand      `command:"eh" description:"Excel helper: generate java class from excel template: eh -p <excel-template-path> [-n <java-class-name>] [-i <sheet-index>]"`
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
	ParentId int    `short:"p" required:"false" description:"Specify parent id for new todo"`
	Deadline string `short:"d" required:"false" description:"Specify deadline for new todo, format: yyyy-MM-dd | MM-dd"`
	Remind   string `short:"r" required:"false" description:"Specify remind time for new todo, format: yyyy-MM-dd HH:mm | MM-dd HH:mm | HH:mm"`
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
	if len(args) == 1 {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}

		todo.CopyContent(id)
		return nil
	}

	id, content, err := argsToIdAndStr(args)
	if err != nil {
		return err
	}

	todo.Modify(id, content)
	todo.Save()
	return nil
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

		todo.CopyRemark(id)
		return nil
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

type TransferCommand struct {
}

func (c *TransferCommand) Execute([]string) error {
	todo.Transfer()
	return nil
}

type TidyCommand struct {
	All      bool `short:"a" required:"false" description:"Execute all tidy works"`
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
	Job       string `short:"j" required:"false" description:"Jenkins job, default: name of current directory"`
	Server    string `short:"s" required:"false" description:"Jenkins parameter [SERVERNAME]"`
	Branch    string `short:"b" required:"false" description:"Jenkins parameter [BUILD_BRANCH], default: branch of current directory"`
	CheckCode bool   `short:"c" required:"false" description:"Jenkins parameter [IS_CHECK_CODE]"`
}

func (c *JenkinsDeployCommand) Execute([]string) error {
	if c.Job == "" {
		c.Job = jenkins.DefaultJob()
	}
	if c.Branch == "" {
		b, err := jenkins.DefaultBranch()
		if err != nil {
			return err
		}
		c.Branch = b
	}
	return jenkins.Deploy(c.Job, c.Server, c.Branch, c.CheckCode)
}

type JenkinsStatusCommand struct {
	Job string `short:"j" required:"false" description:"Jenkins job, default: name of current directory"`
}

func (c *JenkinsStatusCommand) Execute([]string) error {
	if c.Job == "" {
		c.Job = jenkins.DefaultJob()
	}
	return jenkins.Status(c.Job)
}

type JenkinsHistoryCommand struct {
	Job   string `short:"j" required:"false" description:"Jenkins job, default: name of current directory"`
	Count int    `short:"c" required:"false" description:"History count, default: 5"`
}

func (c *JenkinsHistoryCommand) Execute([]string) error {
	if c.Job == "" {
		c.Job = jenkins.DefaultJob()
	}
	if c.Count <= 0 {
		c.Count = 5
	}
	return jenkins.History(c.Job, c.Count)
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
