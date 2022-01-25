package cmd

import (
	"fmt"
	"os"

	"github.com/couchbaselabs/cbdynclusterd/daemon"
	"github.com/spf13/cobra"
)

var createCloudCmd = &cobra.Command{
	Use:   "create-cloud",
	Short: "Creates a new cloud cluster",
	Run: func(cmd *cobra.Command, args []string) {
		checkConfigInitialized()

		flags := cmd.Flags()
		var nodes []string
		var err error
		if nodes, err = flags.GetStringArray("node"); err != nil {
			printAndExit("Invalid node")
		}
		var reqData daemon.CreateCloudClusterJSON
		for i := 0; i < len(nodes); i++ {
			reqData.Services = append(reqData.Services, (nodes)[i])
		}

		var respData daemon.ClusterJSON
		err = serverRestCall("POST", "/create-cloud", reqData, &respData, false)
		if err != nil {
			fmt.Printf("Failed to allocate cluster: %s\n", err)
			os.Exit(1)
		}

		fmt.Printf("%s\n", respData.ID)
	},
}

func init() {
	rootCmd.AddCommand(createCloudCmd)

	createCloudCmd.Flags().StringArray("node", nil, "Comma separated services.")
}
