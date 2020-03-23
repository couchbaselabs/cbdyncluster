package cmd

import (
	"fmt"
	"os"

	"strings"

	"github.com/couchbaselabs/cbdynclusterd/daemon"
	"github.com/couchbaselabs/cbdynclusterd/helper"
	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup <cluster ID>",
	Short: "Setup a cluster",
	Long:  "Setup a cluster using cluster ID",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		checkConfigInitialized()

		path := "/cluster/" + args[0] + "/setup"

		flags := cmd.Flags()
		var nodes []string
		var err error
		var storageMode, bucketOption, userOption string
		var useHostname bool
		var ramQuota int
		var useDevPreview bool
		if nodes, err = flags.GetStringArray("node"); err != nil {
			printAndExit("Invalid node")
		}
		if storageMode, err = cmd.Flags().GetString("storage-mode"); err != nil {
			printAndExit("Invalid storage-mode")
		}
		if ramQuota, err = cmd.Flags().GetInt("ram-quota"); err != nil {
			printAndExit("Invalid ram-quota")
		}
		if bucketOption, err = cmd.Flags().GetString("bucket"); err != nil {
			printAndExit("Invalid bucket option")
		}
		if userOption, err = cmd.Flags().GetString("user"); err != nil {
			printAndExit("Invalid user option")
		}
		if useHostname, err = cmd.Flags().GetBool("use-hostname"); err != nil {
			printAndExit("Invalid use-hostname option")
		}
		if useDevPreview, err = cmd.Flags().GetBool("developer-preview"); err != nil {
			printAndExit("Invalid developer-preview option")
		}

		var reqData daemon.CreateClusterSetupJSON
		for i := 0; i < len(nodes); i++ {
			reqData.Services = append(reqData.Services, (nodes)[i])
		}
		reqData.RamQuota = ramQuota
		reqData.StorageMode = storageMode
		reqData.UseHostname = useHostname
		reqData.Bucket = parseBucketOption(bucketOption)
		reqData.User = parseUserOption(userOption)
		reqData.UseDeveloperPreview = useDevPreview

		var respData daemon.ClusterJSON
		err = serverRestCall("POST", path, reqData, &respData, false)

		if err != nil {
			fmt.Printf("Failed to setup a cluster: %s\n", err)
			os.Exit(1)
		}
		fmt.Printf("%s\n", respData.EntryPoint)
	},
}

func parseUserOption(opt string) *helper.UserOption {
	parsed := strings.Split(opt, ":")

	var userName, userPassword string
	var roles []string
	if len(parsed) > 0 {
		userName = parsed[0]
	}
	if len(parsed) > 1 {
		userPassword = parsed[1]
	}
	if len(parsed) > 2 {
		roles = strings.Split(parsed[2], ",")
	}

	return &helper.UserOption{
		Name:     userName,
		Password: userPassword,
		Roles:    &roles,
	}

}

func parseBucketOption(opt string) *helper.BucketOption {
	parsed := strings.Split(opt, ":")
	var bucketName, bucketType, bucketPassword string
	if len(parsed) > 0 {
		bucketName = parsed[0]
	}
	if len(parsed) > 1 {
		bucketType = parsed[1]
	}
	if len(parsed) > 2 {
		bucketPassword = parsed[2]
	}

	return &helper.BucketOption{
		Name:     bucketName,
		Type:     bucketType,
		Password: bucketPassword,
	}
}

func printAndExit(msg string) {
	fmt.Printf(msg)
	os.Exit(1)
}

func init() {
	rootCmd.AddCommand(setupCmd)

	setupCmd.Flags().StringArray("node", nil, "Comma separated services.")
	setupCmd.Flags().String("storage-mode", "", "set storage mode")
	setupCmd.Flags().Int("ram-quota", 600, "ram quota")
	setupCmd.Flags().String("bucket", "", "Create a bucket <bucket-name>[:<bucket-type, memcached|couchbase|ephemeral[:<bucket password>]]. if only bucket name is given, couchbase bucket will be created. if server is equal or after 5.0, bucket password will be ignored")
	setupCmd.Flags().String("user", "", "Create a user <user-name>:<user-password>[:<user-role>]. creates a user. default role is admin")
	setupCmd.Flags().Bool("use-hostname", false, "Set true to setup a cluster using hostname. default is false")
	setupCmd.Flags().Bool("developer-preview", false, "Set true to enable developer preview. default is false")
}
