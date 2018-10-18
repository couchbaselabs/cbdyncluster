package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var ipsCmd = &cobra.Command{
	Use:   "ips [cluster_id]",
	Short: "Fetches the IPs for a cluster",
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

		Ips := fmt.Sprintf("%s", strings.Join(addresses, ","))
		fmt.Printf("%s\n", Ips)
	},
}

func init() {
	rootCmd.AddCommand(ipsCmd)
}
