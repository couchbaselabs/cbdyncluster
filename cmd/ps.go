package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/couchbaselabs/cbdynclusterd/daemon"
	"github.com/spf13/cobra"
)

var psAllFlag bool

var psCmd = &cobra.Command{
	Use:   "ps",
	Short: "Lists all clusters",
	//Long:  `Allocates a new cluster`,
	Run: func(cmd *cobra.Command, args []string) {
		checkConfigInitialized()

		var respData daemon.GetClustersJSON
		err := serverRestCall("GET", "/clusters", nil, &respData, psAllFlag)
		if err != nil {
			log.Printf("Failed to list clusters: %s\n", err)
			os.Exit(1)
		}

		var clusters []*daemon.Cluster
		for _, jsonCluster := range respData {
			cluster, err := daemon.UnjsonifyCluster(&jsonCluster)
			if err != nil {
				fmt.Printf("Failed to parse list clusters response: %s", err)
				os.Exit(1)
			}

			clusters = append(clusters, cluster)
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
