package cmd

import (
	"fmt"
	"os"

	"github.com/couchbaselabs/cbdynclusterd/daemon"
	"github.com/spf13/cobra"
)

var addCollectionCmd = &cobra.Command{
	Use:   "add-collection <cluster ID>",
	Short: "Adds a collection",
	Long:  `Adds a collection to a cluster using cluster ID`,
	Run: func(cmd *cobra.Command, args []string) {
		checkConfigInitialized()

		if len(args) == 0 {
			printAndExit("No cluster-id provided")
		}

		path := "/cluster/" + args[0] + "/add-collection"

		flags := cmd.Flags()
		var err error
		var name, scopeName, bucketName string
		var useHostname bool
		if name, err = flags.GetString("name"); err != nil {
			printAndExit("Invalid name")
		}
		if scopeName, err = flags.GetString("scope-name"); err != nil {
			printAndExit("Invalid scope-name")
		}
		if bucketName, err = flags.GetString("bucket-name"); err != nil {
			printAndExit("Invalid bucket-name")
		}
		if useHostname, err = flags.GetBool("use-hostname"); err != nil {
			printAndExit("Invalid use-hostname option")
		}

		var reqData daemon.AddCollectionJSON

		reqData.UseHostname = useHostname
		reqData.Name = name
		reqData.ScopeName = scopeName
		reqData.BucketName = bucketName

		err = serverRestCall("POST", path, reqData, nil, false)

		if err != nil {
			fmt.Printf("Failed to add collection to cluster: %s\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(addCollectionCmd)

	addCollectionCmd.Flags().String("name", "", "name of the collection.")
	addCollectionCmd.Flags().String("scope-name", "_default", "name of the scope.")
	addCollectionCmd.Flags().String("bucket-name", "default", "name of the scope.")
	addCollectionCmd.Flags().Bool("use-hostname", false, "set true to use a cluster using hostname. default is false")
}
