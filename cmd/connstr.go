package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var connstrSslFlag bool

var connstrCmd = &cobra.Command{
	Use:   "connstr [cluster_id]",
	Short: "Fetches the connection string for a cluster",
	//Long:  `Fetches the connection string for a cluster`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		checkConfigInitialized()

		cluster, err := getCluster(args[0])
		if err != nil {
			fmt.Printf("Failed to list all clusters: %s\n", err)
			os.Exit(1)
		}

		var addresses []string
		for _, node := range cluster.Nodes {
			addresses = append(addresses, node.IPv4Address)
		}

		scheme := "couchbase"
		if connstrSslFlag {
			scheme = "couchbases"
		}

		connStr := fmt.Sprintf("%s://%s", scheme, strings.Join(addresses, ","))
		fmt.Printf("%s\n", connStr)
	},
}

func init() {
	rootCmd.AddCommand(connstrCmd)

	connstrCmd.Flags().BoolVar(&connstrSslFlag, "ssl", false, "gets the SSL variant of the connection string")
}
