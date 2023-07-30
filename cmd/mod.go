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

var copyContent bool

// modCmd represents the mod command
var modCmd = &cobra.Command{
	Use:   "mod <id> [content]",
	Short: "Modify or copy content of todo",
	Long: `Modify or copy content of todo.

Param:
  [content] can be empty only when the "-c" flag is used.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}

		// copy content
		if copyContent {
			return todo.CopyContent(id)
		}

		// modify content
		_, content, err := argsToIdAndStr(args)
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

	modCmd.Flags().BoolVarP(&copyContent, "copy", "c", false, "Copy content of todo")
}
