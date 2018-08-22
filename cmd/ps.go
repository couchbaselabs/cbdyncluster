package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var psAllFlag bool

var psCmd = &cobra.Command{
	Use:   "ps",
	Short: "Lists all clusters",
	//Long:  `Lists all clusters`,
	Run: func(cmd *cobra.Command, args []string) {
		checkConfigInitialized()

		clusters, err := getAllClusters()
		if err != nil {
			fmt.Printf("Failed to list all clusters: %s\n", err)
			os.Exit(1)
		}

		fmt.Printf("Clusters:\n")
		for _, cluster := range clusters {
			fmt.Printf("  %s [Owner: %s, Creator: %s, Timeout: %s]\n", cluster.ID, cluster.Owner, cluster.Creator, cluster.Timeout.Sub(time.Now()).Round(time.Second))
			for _, node := range cluster.Nodes {
				fmt.Printf("    %-16s  %-20s %-10s %-20s\n", node.ContainerID, node.Name, node.InitialServerVersion, node.IPv4Address)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(psCmd)

	psCmd.Flags().BoolVarP(&psAllFlag, "all", "a", false, "list all clusters")
}
