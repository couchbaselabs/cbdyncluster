package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var connstrSslFlag, connstrSrvFlag bool

var connstrCmd = &cobra.Command{
	Use:   "connstr [cluster_id]",
	Short: "Fetches the connection string for a cluster",
	//Long:  `Fetches the connection string for a cluster`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		checkConfigInitialized()

		connStr, err := getConnString(args[0], connstrSslFlag, connstrSrvFlag)

		if err != nil {
			fmt.Printf("Failed to get connection string: %s\n", err)
			os.Exit(1)
		}

		fmt.Printf("%s\n", connStr)
	},
}

func init() {
	rootCmd.AddCommand(connstrCmd)

	connstrCmd.Flags().BoolVar(&connstrSslFlag, "ssl", false, "gets the SSL variant of the connection string")
	connstrCmd.Flags().BoolVar(&connstrSrvFlag, "srv", false, "gets the DNS SRV variant of the connection string (only supported on AWS and Capella)")
}
