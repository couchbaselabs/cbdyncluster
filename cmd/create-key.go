package cmd

import (
	"fmt"
	"github.com/couchbaselabs/cbdynclusterd/service/cloud"
	"os"

	"github.com/spf13/cobra"
)

var createKeysCmd = &cobra.Command{
	Use:   "create-key <cluster ID>",
	Short: "Creates a api key pair for columnar instnace",
	Long:  `Creates an API key pair (api id/secret) using the provided columnar ID`,
	Run: func(cmd *cobra.Command, args []string) {
		checkConfigInitialized()

		if len(args) == 0 {
			printAndExit("No cluster-id provided")
		}

		path := "/columnar/" + args[0] + "/create-key"

		var respData cloud.ColumnarApiKeys
		err := serverRestCall("POST", path, nil, &respData, false)

		if err != nil {
			fmt.Printf("Failed to create key: %s\n", err)
			os.Exit(1)
		}
		fmt.Printf("%s,%s\n", respData.APIKeyId, respData.Secret)
	},
}

func init() {
	rootCmd.AddCommand(createKeysCmd)
}
