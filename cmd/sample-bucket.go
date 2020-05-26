package cmd

import (
	"fmt"
	"os"

	"github.com/couchbaselabs/cbdynclusterd/daemon"
	"github.com/spf13/cobra"
)

var addSampleBucketCmd = &cobra.Command{
	Use:   "add-sample-bucket <cluster ID>",
	Short: "Loads a sample bucket",
	Long:  `Loads the sample bucket specified to a cluster using cluster ID`,
	Run: func(cmd *cobra.Command, args []string) {
		checkConfigInitialized()

		if len(args) == 0 {
			printAndExit("No cluster-id provided")
		}

		path := "/cluster/" + args[0] + "/add-sample-bucket"

		flags := cmd.Flags()
		var err error
		var name string
		var useHostname bool
		if name, err = flags.GetString("name"); err != nil {
			printAndExit("Invalid name")
		}
		if useHostname, err = flags.GetBool("use-hostname"); err != nil {
			printAndExit("Invalid use-hostname option")
		}

		var reqData daemon.AddSampleBucketJSON

		reqData.UseHostname = useHostname
		reqData.SampleBucket = name

		err = serverRestCall("POST", path, reqData, nil, false)

		if err != nil {
			fmt.Printf("Failed to load sample bucket: %s\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(addSampleBucketCmd)

	addSampleBucketCmd.Flags().String("name", "", "name of the bucket.")
	addSampleBucketCmd.Flags().Bool("use-hostname", false, "set true to setup a cluster using hostname. default is false")
}
