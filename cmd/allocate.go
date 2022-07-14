package cmd

import (
	"fmt"
	"os"

	"github.com/couchbaselabs/cbdynclusterd/daemon"
	"github.com/spf13/cobra"
)

var allocateCmd = &cobra.Command{
	Use:   "allocate",
	Short: "Allocates a new cluster",
	//Long:  `Allocates a new cluster`,
	Run: func(cmd *cobra.Command, args []string) {
		checkConfigInitialized()

		isSimpleInvoke := false
		isSimpleInvoke = isSimpleInvoke || cmd.Flags().Changed("num-nodes")
		isSimpleInvoke = isSimpleInvoke || cmd.Flags().Changed("server-version")
		isSimpleInvoke = isSimpleInvoke || cmd.Flags().Changed("use-ce")

		numNodes, _ := cmd.Flags().GetInt("num-nodes")
		serverVersion, _ := cmd.Flags().GetString("server-version")
		useCE, _ := cmd.Flags().GetBool("use-ce")
		platform, _ := cmd.Flags().GetString("platform")
		OS, _ := cmd.Flags().GetString("os")
		arch, _ := cmd.Flags().GetString("arch")
		serverlessMode, _ := cmd.Flags().GetBool("serverless-mode")

		if numNodes < 0 || numNodes > 24 {
			fmt.Printf("Must allocate between 1 and 24 nodes\n")
			os.Exit(1)
		}

		var reqData daemon.CreateClusterJSON
		for i := 0; i < numNodes; i++ {
			reqData.Nodes = append(reqData.Nodes, daemon.CreateClusterNodeJSON{
				ServerVersion:       serverVersion,
				UseCommunityEdition: useCE,
				Platform:            platform,
				OS:                  OS,
				Arch:                arch,
				ServerlessMode:      serverlessMode,
			})
		}

		var respData daemon.NewClusterJSON
		err := serverRestCall("POST", "/clusters", reqData, &respData, false)
		if err != nil {
			fmt.Printf("Failed to allocate cluster: %s\n", err)
			os.Exit(1)
		}

		fmt.Printf("%s\n", respData.ID)
	},
}

func init() {
	rootCmd.AddCommand(allocateCmd)

	allocateCmd.Flags().Int("num-nodes", 3, "The number of nodes to initialize")
	allocateCmd.Flags().String("server-version", "5.5.0", "The server version to use when allocating the nodes.")
	allocateCmd.Flags().Bool("use-ce", false, "Use the Community edition (CE) of the Couchbase Server.")
	allocateCmd.Flags().String("platform", "docker", "The platform to use when allocating the nodes.")
	allocateCmd.Flags().String("os", "centos7", "The operating system to use")
	allocateCmd.Flags().String("arch", "x86_64", "The CPU architecture to use")
	allocateCmd.Flags().Bool("serverless-mode", false, "Start up the cluster in serverless mode")
}
