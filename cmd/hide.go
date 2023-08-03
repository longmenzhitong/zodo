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

// hideCmd represents the hide command
var hideCmd = &cobra.Command{
	Use:   "hide <id>...",
	Short: `Set todo status to "Hiding"`,
	Long: `Set todo status to "Hiding".

Note:
  Different from "Done" status, only "-a" flag of "list" command 
  can show the "Hiding" status todos in the list.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ids, err := argsToIds(args)
		if err != nil {
			return err
		}

		for _, id := range ids {
			todo.SetStatus(id, todo.StatusHiding)
		}
		todo.Save()
		return nil
	},
}

func init() {
	RootCmd.AddCommand(hideCmd)
}
