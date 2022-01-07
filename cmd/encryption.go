package cmd

import (
	"fmt"
	"os"

	"github.com/couchbaselabs/cbdynclusterd/daemon"
	"github.com/spf13/cobra"
)

var setupClusterEncryptionCmd = &cobra.Command{
	Use:   "setup-cluster-encryption <cluster ID>",
	Short: "Setup cluster encryption",
	Long:  "Setup cluster encryption using cluster ID",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		checkConfigInitialized()

		path := "/cluster/" + args[0] + "/setup-cluster-encryption"

		var err error
		var level string

		if level, err = cmd.Flags().GetString("level"); err != nil {
			printAndExit("Invalid level option")
		}

		var reqData daemon.SetupClusterEncryptionJSON

		reqData.Level = level

		err = serverRestCall("POST", path, reqData, nil, false)

		if err != nil {
			fmt.Printf("Failed to setup cluster encryption: %s\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(setupClusterEncryptionCmd)

	setupClusterEncryptionCmd.Flags().String("level", "", "Set cluster encryption level, default is cluster default (control)")
}
