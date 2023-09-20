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
	zodo "zodo/src"
	"zodo/src/todo"

	"github.com/spf13/cobra"
)

// ddlCmd represents the ddl command
var ddlCmd = &cobra.Command{
	Use:   "ddl <id> [deadline]",
	Short: "Set deadline of todo",
	Long: `Set deadline of todo.

Param:
  [deadline] can be "yyyy-MM-dd" or "MM-dd" or empty.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		id, deadline, err := argsToIdAndOptionalStr(args)
		if err != nil {
			return err
		}

		if deadline == "" {
			td := todo.Get(id)
			deadline, err = zodo.EditByVim(td.Deadline)
			if err != nil {
				return err
			}
		}

		if deadline != "" {
			deadline, err = validateDeadline(deadline)
			if err != nil {
				return err
			}
		}

		todo.SetDeadline(id, deadline)
		todo.Save()
		return nil
	},
}

func init() {
	RootCmd.AddCommand(ddlCmd)
}
