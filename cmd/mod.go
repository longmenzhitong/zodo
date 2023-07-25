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

// modCmd represents the mod command
var modCmd = &cobra.Command{
	Use:   "mod <id> [content]",
	Short: "Modify content of todo",
	Long:  `Modify content of todo, will copy content of todo if only id was provided.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 1 {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}

			return todo.CopyContent(id)
		}

		id, content, err := argsToIdAndStr(args)
		if err != nil {
			return err
		}

		todo.Modify(id, content)
		todo.Save()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(modCmd)
}
