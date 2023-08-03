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
	"zodo/src/todo"

	"github.com/spf13/cobra"
)

var copyRemindTime bool

// rmdCmd represents the rmd command
var rmdCmd = &cobra.Command{
	Use:   "rmd <id> [remindTime]",
	Short: "Set or copy remind time of todo",
	Long: `Set or copy remind time of todo.

Param:
  [remindTime] can be "yyyy-MM-dd HH:mm" or "MM-dd HH:mm" or "HH:mm" or empty.

Need:
  The remind feature needs:
  1. Runing ZODO in server mode by using the "server" command;
  2. Set the config "reminder.enabled" to "true" and set a cron;
  3. An email server and corresponding config "email".`,
	RunE: func(cmd *cobra.Command, args []string) error {
		id, remindTime, err := argsToIdAndOptionalStr(args)
		if err != nil {
			return err
		}

		if copyRemindTime {
			return todo.CopyRemindTime(id)
		}

		if remindTime != "" {
			remindTime, err = validateRemind(remindTime)
			if err != nil {
				return err
			}
		}
		err = todo.SetRemind(id, remindTime)
		if err != nil {
			return err
		}
		todo.Save()
		return nil
	},
}

func init() {
	RootCmd.AddCommand(rmdCmd)

	rmdCmd.Flags().BoolVarP(&copyRemindTime, "copy", "c", false, "Copy remind time of todo")
}
