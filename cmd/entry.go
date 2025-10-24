/*
mixnet - tool to create and manage LibSEA mixnets
Copyright (C) 2025  Liberatory Sofware Engineering Association

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/
package cmd

import (
	"os"

	"github.com/LibSEA/mixnet/entry"
	"github.com/spf13/cobra"
)

// entryCmd represents the entry command
var entryCmd = &cobra.Command{
	Use:   "entry",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
        host, _:= cmd.Flags().GetString("host")
        port, _:= cmd.Flags().GetUint16("port")
		os.Exit(entry.Run(entry.Options{
    		Port: port,
    		Host: host,
		}))
	},
}

func init() {
	rootCmd.AddCommand(entryCmd)

	entryCmd.PersistentFlags().String(
    	"host",
    	"localhost",
    	"host to listen on",
	)
	entryCmd.PersistentFlags().Uint16("port", 8080, "port to connect to")
}
