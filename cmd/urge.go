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

// urgeCmd represents the urge command
var urgeCmd = &cobra.Command{
	Use:   "urge <id>",
	Short: "Raise the priority of specified todo",
	Long: `Raise the priority of specified todo. 
Todos with higher priority will be at the top of the list.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ids, err := argsToIds(args)
		if err != nil {
			return err
		}
		todo.Urge(ids[0])
		todo.Save()
		return nil
	},
}

func init() {
	RootCmd.AddCommand(urgeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// urgeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// urgeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
