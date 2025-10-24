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

	"github.com/LibSEA/mixnet/ping"
	"github.com/spf13/cobra"
)

var pingOpts = struct {
    host string
    port uint16
} {
    host: "locahost",
    port: 8080,
}


// pingCmd represents the ping command
var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		os.Exit(ping.Run(ping.Options{
    		Host: pingOpts.host,
    		Port: pingOpts.port,
		}))
	},
}

func init() {
	rootCmd.AddCommand(pingCmd)

	pingCmd.PersistentFlags().StringVar(&pingOpts.host, "host", pingOpts.host, "host to connect to")
	pingCmd.PersistentFlags().Uint16Var(&pingOpts.port, "port", pingOpts.port, "port to connect to")
}
