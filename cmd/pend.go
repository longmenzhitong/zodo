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

// pendCmd represents the pend command
var pendCmd = &cobra.Command{
	Use:   "pend <id>...",
	Short: `Set todo status to "Pending"`,
	Long:  `Set todo status to "Pending".`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ids, err := argsToIds(args)
		if err != nil {
			return err
		}

		for _, id := range ids {
			todo.SetStatus(id, todo.StatusPending)
			todo.AddPriority(id, -1)
		}
		todo.Save()
		return nil
	},
}

func init() {
	RootCmd.AddCommand(pendCmd)
}
