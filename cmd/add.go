/*
Copyright © 2023 zhihaoyu <longmenzhitong@gmail.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"strconv"

	zodo "zodo/src"
	"zodo/src/todo"

	"github.com/spf13/cobra"
)

var parentId int
var deadline string
var remindTime string
var remark string

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add <content of new todo>",
	Short: "Add new todo",
	Long: `Add new todo, optionally specify parent ID, deadline, remind time, and remark of new todo.
This command will write the id of new todo into clipboard. You can set the config "todo.copyIdAfterAdd" to false to disable this feature.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := todo.Add(argsToStr(args))
		if err != nil {
			return err
		}

		if parentId != -1 {
			err = todo.SetChild(parentId, []int{id}, true)
			if err != nil {
				return err
			}
		}

		if deadline != "" {
			ddl, err := validateDeadline(deadline)
			if err != nil {
				return err
			}
			todo.SetDeadline(id, ddl)
		}

		if remindTime != "" {
			rmd, err := validateRemind(remindTime)
			if err != nil {
				return err
			}
			err = todo.SetRemind(id, rmd)
			if err != nil {
				return err
			}
		}

		if remark != "" {
			todo.SetRemark(id, remark)
		}

		todo.Save()

		if zodo.Config.Todo.CopyIdAfterAdd {
			err := zodo.WriteLineToClipboard(strconv.Itoa(id))
			if err != nil {
				return err
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().IntVarP(&parentId, "parent", "p", -1, "Specify parent ID of new todo")
	addCmd.Flags().StringVarP(&deadline, "deadline", "d", "", `Specify deadline of new todo, accept "yyyy-MM-dd" or "MM-dd"`)
	addCmd.Flags().StringVarP(&remindTime, "remind", "r", "", `Specify remind time of new todo, accept "yyyy-MM-dd HH:mm" or "MM-dd HH:mm" or "HH:mm"`)
	addCmd.Flags().StringVarP(&remark, "remark", "R", "", "Specify remark of new todo")
}
