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

var doneTodos bool
var fragmentIds bool

// tidyCmd represents the tidy command
var tidyCmd = &cobra.Command{
	Use:   "tidy",
	Short: "Do some tidy work",
	Long:  `Do some tidy work.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		changed := false
		if all || doneTodos {
			count := todo.ClearDoneTodo()
			if count > 0 {
				zodo.PrintDoneMsg("Clear %d done todos.\n", count)
				changed = true
			}
		}
		if all || fragmentIds {
			from, to := todo.DefragId()
			if from != to {
				zodo.PrintDoneMsg("Defrag ids from %d to %d.\n", from, to)
				changed = true
			}
		}
		if changed {
			todo.Save()
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(tidyCmd)

	tidyCmd.Flags().BoolVarP(&all, "all", "a", false, "Do all tidy work")
	tidyCmd.Flags().BoolVarP(&doneTodos, "done", "d", false, "Tidy done todos")
	tidyCmd.Flags().BoolVarP(&fragmentIds, "fragment", "f", false, "Tidy fragment ids")
}
