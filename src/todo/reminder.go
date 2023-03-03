package todo

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"time"
	"zodo/src"
)

const (
	rmdStatusWaiting  = "Waiting"
	rmdStatusFinished = "Finished"
)

const (
	loopOnce         = "Once"
	loopPerDay       = "Per Day"
	loopPerWorkDay   = "Per Work Day"
	loopPerMonday    = "Per Monday"
	loopPerTuesday   = "Per Tuesday"
	loopPerWednesday = "Per Wednesday"
	loopPerThursday  = "Per Thursday"
	loopPerFriday    = "Per Friday"
	loopPerSaturday  = "Per Saturday"
	loopPerSunday    = "Per Sunday"
)

var loopTypes = []string{
	loopOnce,
	loopPerDay,
	loopPerWorkDay,
	loopPerMonday,
	loopPerTuesday,
	loopPerWednesday,
	loopPerThursday,
	loopPerFriday,
	loopPerSaturday,
	loopPerSunday,
}

func SetRemind(id int, rmdTime string, loop bool) error {
	td := cc._map()[id]
	if td == nil {
		return &zodo.NotFoundError{
			Target:  "todo",
			Message: fmt.Sprintf("id: %d", id),
		}
	}

	if loop {
		rows := make([]table.Row, 0)
		for i := 0; i < len(loopTypes); i++ {
			rows = append(rows, table.Row{i, loopTypes[i]})
		}
		zodo.PrintTable(table.Row{"Num", "Type"}, rows)
		num, err := zodo.ReadInt(0, len(loopTypes), "Enter number of remind type:")
		if err != nil {
			return err
		}
		td.LoopType = loopTypes[num]
	} else {
		td.LoopType = loopOnce
	}

	td.RemindTime = rmdTime
	td.RemindStatus = rmdStatusWaiting

	return nil
}

func RemoveRemind(ids []int) {
	m := cc._map()
	for _, id := range ids {
		td := m[id]
		if td != nil {
			td.RemindTime = ""
			td.RemindStatus = ""
			td.LoopType = ""
		}
	}
}

func Remind() error {
	cc.refresh()
	var text string
	m := cc._map()
	for _, td := range cc.list("", []string{}, true) {
		if !isNeedRemind(td.RemindTime, td.LoopType, td.RemindStatus, time.Now()) {
			continue
		}

		if td.LoopType == loopOnce {
			m[td.Id].RemindStatus = rmdStatusFinished
		}

		ddl, remain := td.getDeadLineAndRemain(false)
		text += "\n"
		if ddl != "" {
			text += fmt.Sprintf("* %s  %s, deadline %s, remain %s\n", td.Content, td.getStatus(false), ddl, remain)
		} else if td.Status != statusHiding {
			text += fmt.Sprintf("* %s  %s\n", td.Content, td.getStatus(false))
		} else {
			text += fmt.Sprintf("* %s\n", td.Content)
		}
	}
	if text != "" {
		err := zodo.SendEmail("Reminder", text)
		if err != nil {
			return err
		}
		cc.save()
	}
	return nil
}

func isNeedRemind(rmdTime, loopType, rmdStatus string, checkTime time.Time) bool {
	if rmdTime == "" || loopType == "" || rmdStatus == "" {
		return false
	}
	if rmdStatus == rmdStatusFinished {
		return false
	}
	t, err := time.ParseInLocation(zodo.LayoutDateTime, rmdTime, time.Local)
	if err != nil {
		panic(err)
	}
	if loopType == loopOnce {
		return checkTime.Equal(t) || checkTime.After(t)
	}
	if t.Hour() != checkTime.Hour() || t.Minute() != checkTime.Minute() || t.Second() != checkTime.Second() {
		return false
	}
	wd := checkTime.Weekday()
	switch loopType {
	case loopPerDay:
		return true
	case loopPerWorkDay:
		return wd != time.Saturday && wd != time.Sunday
	case loopPerMonday:
		return wd == time.Monday
	case loopPerTuesday:
		return wd == time.Tuesday
	case loopPerWednesday:
		return wd == time.Wednesday
	case loopPerThursday:
		return wd == time.Thursday
	case loopPerFriday:
		return wd == time.Friday
	case loopPerSaturday:
		return wd == time.Saturday
	case loopPerSunday:
		return wd == time.Sunday
	}
	panic(&zodo.InvalidInputError{
		Message: fmt.Sprintf("remindTime: %s, loopType: %s, remindStatus: %s, checkTime: %v",
			rmdTime, loopType, rmdStatus, checkTime),
	})
}
