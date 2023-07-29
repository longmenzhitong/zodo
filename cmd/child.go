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
	Use:   "child",
	Short: "Add child of todo",
	Long:  ``,
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
	rootCmd.AddCommand(childCmd)

	childCmd.Flags().BoolVarP(&overwriteOtherChildren, "overwrite", "o", false, "Overwrite other children of todo")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// childCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// childCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
