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
	"zodo/src/todo"

	"github.com/spf13/cobra"
)

var copy bool

// rmkCmd represents the rmk command
var rmkCmd = &cobra.Command{
	Use:   "rmk <id> [remark of todo]",
	Short: "Set or copy remark of todo",
	Long: `Set or copy remark of todo:
* Set remark of todo: rmk <id> [remark of todo]
* Copy remark of todo: rmk -c <id>`,
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}

		// copy remark
		if copy {
			return todo.CopyRemark(id)
		}

		// set remark
		var remark string
		if len(args) == 1 {
			remark = ""
		} else {
			_, remark, err = argsToIdAndStr(args)
			if err != nil {
				return err
			}
		}

		todo.SetRemark(id, remark)
		todo.Save()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(rmkCmd)

	rmkCmd.Flags().BoolVarP(&copy, "copy", "c", false, "Copy remark of todo")
}
