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

var all bool

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls [<keyword of content>]",
	Short: "Show todo list",
	Long:  `Show todo list, optionally filter by keyword of content, or show todos in all statuses`,
	Run: func(cmd *cobra.Command, args []string) {
		todo.List(argsToStr(args), all)
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)

	lsCmd.Flags().BoolVarP(&all, "all", "a", false, "Show todos in all statuses, including done and hiding todos")
}
