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

var decreasePriority bool
var restorePriority bool

// urgeCmd represents the urge command
var urgeCmd = &cobra.Command{
	Use:   "urge <id>...",
	Short: "Raise the priority of specified todo",
	Long: `Raise the priority of specified todo. 
Todos with higher priority will be at the top of the list.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ids, err := argsToIds(args)
		if err != nil {
			return err
		}

		for _, id := range ids {
			if restorePriority {
				todo.SetPriority(id, 0)
			} else {
				var p int
				if decreasePriority {
					p = -1
				} else {
					p = 1
				}
				todo.AddPriority(id, p)
			}
		}

		todo.Save()
		return nil
	},
}

func init() {
	RootCmd.AddCommand(urgeCmd)

	urgeCmd.Flags().BoolVarP(&decreasePriority, "decrease", "d", false, "Decrease the priority of specified todo")
	urgeCmd.Flags().BoolVarP(&restorePriority, "restore", "r", false, "Restore the priority of specified todo")
}
