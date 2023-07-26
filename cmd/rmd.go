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
	Use:   "rmd",
	Short: "Set or copy remind time of todo",
	Long: `Use:
  * Set remind time of todo: rmd <id> [remind time], accept "yyyy-MM-dd HH:mm" or "MM-dd HH:mm" or "HH:mm" or empty
  * Copy remind time of todo: rmd -c <id>

Config:
  The remind feature requires the following two configurations:
  * reminder
  * email
  Use "conf" command to see the detail of these two configurations.

Dependency:
  The remind feature implements by cron job and email. So it depends on "server" command to start the remind cron job.
  See the help info of "server" command for further understanding.`,
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
	rootCmd.AddCommand(rmdCmd)

	rmdCmd.Flags().BoolVarP(&copyRemindTime, "copy", "c", false, "Copy remind time of todo")
}
