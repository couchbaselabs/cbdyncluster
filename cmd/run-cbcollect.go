package cmd

import (
	"fmt"
	"os"

	"github.com/couchbaselabs/cbdynclusterd/daemon"
	"github.com/spf13/cobra"
)

var cbcollectCmd = &cobra.Command{
	Use:   "cbcollect <cluster ID>",
	Short: "Runs cbcollect-info on a cluster",
	Long:  `Runs cbcollect-info on a cluster`,
	Run: func(cmd *cobra.Command, args []string) {
		checkConfigInitialized()

		if len(args) == 0 {
			printAndExit("No cluster-id provided")
		}

		path := "/cluster/" + args[0] + "/cbcollect"

		flags := cmd.Flags()
		outOption, err := flags.GetString("out-dir")
		if outOption == "" {
			printAndExit("No out-dir provided")
		}

		var respData daemon.CBCollectResultJSON

		err = serverRestCall("GET", path, nil, &respData, false)

		if err != nil {
			fmt.Printf("Failed to run cbcollect-info: %s\n", err)
			os.Exit(1)
		}

		for hostname, bytes := range respData.Collections {
			writeBytes(fmt.Sprintf("%s/cbcollect-%s.zip", outOption, hostname), bytes)
		}
	},
}

func init() {
	rootCmd.AddCommand(cbcollectCmd)

	cbcollectCmd.Flags().String("out-dir", ".", "The directory to write cbcollects to")
}
