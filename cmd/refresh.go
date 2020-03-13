package cmd

import (
	"fmt"
	"os"

	"github.com/couchbaselabs/cbdynclusterd/daemon"
	"github.com/spf13/cobra"
)

var refreshCmd = &cobra.Command{
	Use:   "refresh <cluster_id> <timeout>",
	Short: "Updates timeout of the cluster",
	Long:  "Timeout is in duration format of Golang which is possibly signed sequence of decimal numbers, each with optional fraction and a unit suffix, such as \"300s\", \"+1.5h\" or \"2h25m\"",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		checkConfigInitialized()

		clusterId := args[0]
		_, err := getCluster(clusterId)
		if err != nil {
			fmt.Printf("Failed to find the cluster %s: %s\n", clusterId, err)
			os.Exit(1)
		}

		var reqData daemon.RefreshJSON
		reqData.Timeout = args[1]
		err = serverRestCall("PUT", "/cluster/"+clusterId, reqData, nil, false)
		if err != nil {
			fmt.Printf("Failed to refresh cluster %s: %s\n", clusterId, err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(refreshCmd)
}
