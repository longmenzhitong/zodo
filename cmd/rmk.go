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

// rmkCmd represents the rmk command
var rmkCmd = &cobra.Command{
	Use:   "rmk <id> [remark]",
	Short: "Set remark of todo",
	Long:  `Set remark of todo.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		id, remark, err := argsToIdAndOptionalStr(args)
		if err != nil {
			return err
		}

		if remark == "" {
			td := todo.Get(id)
			remark, err = zodo.EditByVim(td.Remark)
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
	RootCmd.AddCommand(rmkCmd)
}
