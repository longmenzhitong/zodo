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
	"fmt"
	zodo "zodo/src"
	"zodo/src/todo"

	"github.com/spf13/cobra"
)

var overwriteOtherChildren bool

// childCmd represents the child command
var childCmd = &cobra.Command{
	Use:   "child <parentId> <childId>...",
	Short: "Add child of todo",
	Long: `Add child of todo.

Note:
  By default, parent todos which have at least one child do not show their status
  in the list. Set the config "todo.showParentStatus" to "true" to change it.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ids, err := argsToIds(args)
		if err != nil {
			return err
		}
		if len(ids) < 2 {
			return &zodo.InvalidInputError{
				Message: fmt.Sprintf("expect: <parentId> <childId>..., got: %v", args),
			}
		}
		err = todo.SetChild(ids[0], ids[1:], !overwriteOtherChildren)
		if err != nil {
			return err
		}
		todo.Save()
		return nil
	},
}

func init() {
	RootCmd.AddCommand(childCmd)

	childCmd.Flags().BoolVarP(&overwriteOtherChildren, "overwrite", "o", false, "Overwrite other children of todo")
}
