package cmd

import (
	"fmt"
	"os"

	"github.com/couchbaselabs/cbdynclusterd/daemon"
	"github.com/couchbaselabs/cbdynclusterd/store"
	"github.com/spf13/cobra"
)

var createColumnarCmd = &cobra.Command{
	Use:   "create-columnar",
	Short: "Creates a new cloud columnar",
	Run: func(cmd *cobra.Command, args []string) {
		checkConfigInitialized()

		flags := cmd.Flags()
		var url string
		var accessKey string
		var secretKey string
		var username string
		var password string
		var tenant string
		var project string
		var region string
		var provider string
		var singleAZ *bool
		var env string

		var nodes int

		var err error

		var reqData daemon.CreateColumnarJSON

		if url, err = flags.GetString("url"); err != nil {
			printAndExit("Invalid url")
		}
		if accessKey, err = flags.GetString("access-key"); err != nil {
			printAndExit("Invalid access key")
		}
		if secretKey, err = flags.GetString("secret-key"); err != nil {
			printAndExit("Invalid secret key")
		}
		if username, err = flags.GetString("username"); err != nil {
			printAndExit("Invalid username")
		}
		if password, err = flags.GetString("password"); err != nil {
			printAndExit("Invalid password")
		}
		if tenant, err = flags.GetString("tenant"); err != nil {
			printAndExit("Invalid tenant")
		}
		if project, err = flags.GetString("project"); err != nil {
			printAndExit("Invalid project")
		}
		if region, err = flags.GetString("region"); err != nil {
			printAndExit("Invalid region")
		}
		if provider, err = flags.GetString("provider"); err != nil {
			printAndExit("Invalid provider")
		}
		if flags.Changed("single-az") {
			if *singleAZ, err = flags.GetBool("single-az"); err != nil {
				printAndExit("Invalid single-az")
			}
		}
		if env, err = flags.GetString("env"); err != nil {
			printAndExit("Invalid env")
		}

		if nodes, err = flags.GetInt("nodes"); err != nil {
			printAndExit("Invalid region")
		}

		if url != "" {
			reqData.Environment = &store.CloudEnvironment{
				TenantID:  tenant,
				ProjectID: project,
				URL:       url,
				AccessKey: accessKey,
				SecretKey: secretKey,
				Username:  username,
				Password:  password,
			}
		}

		if region != "" {
			reqData.Region = region
		}

		if provider != "" {
			reqData.Provider = provider
		}

		if env != "" {
			reqData.EnvName = env
		}

		reqData.Nodes = nodes

		var respData daemon.ClusterJSON
		err = serverRestCall("POST", "/create-columnar", reqData, &respData, false)
		if err != nil {
			fmt.Printf("Failed to allocate cluster: %s\n", err)
			os.Exit(1)
		}

		fmt.Printf("%s\n", respData.ID)
	},
}

func init() {
	rootCmd.AddCommand(createColumnarCmd)

	createColumnarCmd.Flags().Int("nodes", 1, "Number of compute nodes to use.")
	createColumnarCmd.Flags().String("url", "", "Environemnt url e.g. cloud.couchbase.com")
	createColumnarCmd.Flags().String("access-key", "", "Access key")
	createColumnarCmd.Flags().String("secret-key", "", "Secret key")
	createColumnarCmd.Flags().String("username", "", "Login username")
	createColumnarCmd.Flags().String("password", "", "Login password")
	createColumnarCmd.Flags().String("tenant", "", "Tenant ID")
	createColumnarCmd.Flags().String("project", "", "Project ID")
	createColumnarCmd.Flags().String("region", "", "Region e.g. us-east-1")
	createColumnarCmd.Flags().String("provider", "", "Provider e.g. aws")
	createColumnarCmd.Flags().Bool("single-az", true, "Only deploy in one availability zone")
	createColumnarCmd.Flags().String("image", "", "Custom image e.g. couchbase-cloud-server-7.2.0-1409-qe")
	createColumnarCmd.Flags().String("override-token", "", "Override token to use non default deployment options")
	createColumnarCmd.Flags().String("env", "", "Predefined environment (e.g. prod, stage, dev)")
	createColumnarCmd.Flags().String("server", "", "Server version to use e.g. 7.1.3, Using an image will override this.")
}
