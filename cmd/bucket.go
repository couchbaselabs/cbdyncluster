package cmd

import (
	"fmt"
	"os"

	"github.com/couchbaselabs/cbdynclusterd/daemon"
	"github.com/spf13/cobra"
)

var addBucketCmd = &cobra.Command{
	Use:   "add-bucket <cluster ID>",
	Short: "Adds a bucket",
	Long:  `Adds a bucket to a cluster using cluster ID`,
	Run: func(cmd *cobra.Command, args []string) {
		checkConfigInitialized()

		if len(args) == 0 {
			printAndExit("No cluster-id provided")
		}

		path := "/cluster/" + args[0] + "/add-bucket"

		flags := cmd.Flags()
		var err error
		var name, bucketType string
		var ramQuota, replicaCount int
		var useHostname bool
		if name, err = flags.GetString("name"); err != nil {
			printAndExit("Invalid name")
		}
		if ramQuota, err = flags.GetInt("ram-quota"); err != nil {
			printAndExit("Invalid ram-quota")
		}
		if bucketType, err = flags.GetString("type"); err != nil {
			printAndExit("Invalid bucket type")
		}
		if replicaCount, err = flags.GetInt("replica-count"); err != nil {
			printAndExit("Invalid bucket type")
		}
		if useHostname, err = flags.GetBool("use-hostname"); err != nil {
			printAndExit("Invalid use-hostname option")
		}

		var reqData daemon.AddBucketJSON

		reqData.RamQuota = ramQuota
		reqData.UseHostname = useHostname
		reqData.Name = name
		reqData.BucketType = bucketType
		reqData.ReplicaCount = replicaCount

		err = serverRestCall("POST", path, reqData, nil, false)

		if err != nil {
			fmt.Printf("Failed to setup a cluster: %s\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(addBucketCmd)

	addBucketCmd.Flags().Int("ram-quota", 200, "ram quota")
	addBucketCmd.Flags().String("name", "", "name of the bucket.")
	addBucketCmd.Flags().String("type", "couchbase", "type of the bucket.")
	addBucketCmd.Flags().Int("replica-count", 1, "number of replicas")
	addBucketCmd.Flags().Bool("use-hostname", false, "set true to setup a cluster using hostname. default is false")
}
