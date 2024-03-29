package cmd

import (
	"fmt"
	"os"

	"github.com/couchbaselabs/cbdynclusterd/daemon"
	"github.com/couchbaselabs/cbdynclusterd/store"
	"github.com/spf13/cobra"
)

var createCloudCmd = &cobra.Command{
	Use:   "create-cloud",
	Short: "Creates a new cloud cluster",
	Run: func(cmd *cobra.Command, args []string) {
		checkConfigInitialized()

		flags := cmd.Flags()
		var nodes []string
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
		var image string
		var overrideToken string
		var env string
		var err error
		var server string

		if nodes, err = flags.GetStringArray("node"); err != nil {
			printAndExit("Invalid node")
		}
		var reqData daemon.CreateCloudClusterJSON
		for i := 0; i < len(nodes); i++ {
			reqData.Services = append(reqData.Services, (nodes)[i])
		}

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
		if image, err = flags.GetString("image"); err != nil {
			printAndExit("Invalid image")
		}
		if overrideToken, err = flags.GetString("override-token"); err != nil {
			printAndExit("Invalid override token")
		}
		if env, err = flags.GetString("env"); err != nil {
			printAndExit("Invalid env")
		}

		if server, err = flags.GetString("server"); err != nil {
			printAndExit("Invalid server")
		}

		if url != "" {
			reqData.Environment = &store.CloudEnvironment{
				TenantID:      tenant,
				ProjectID:     project,
				URL:           url,
				AccessKey:     accessKey,
				SecretKey:     secretKey,
				Username:      username,
				Password:      password,
				OverrideToken: overrideToken,
			}
		}

		if region != "" {
			reqData.Region = region
		}

		if provider != "" {
			reqData.Provider = provider
		}

		if image != "" {
			reqData.Image = image
		}

		if env != "" {
			reqData.EnvName = env
		}

		reqData.SingleAZ = singleAZ
		reqData.Server = server

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
	createCloudCmd.Flags().String("url", "", "Environemnt url e.g. cloud.couchbase.com")
	createCloudCmd.Flags().String("access-key", "", "Access key")
	createCloudCmd.Flags().String("secret-key", "", "Secret key")
	createCloudCmd.Flags().String("username", "", "Login username")
	createCloudCmd.Flags().String("password", "", "Login password")
	createCloudCmd.Flags().String("tenant", "", "Tenant ID")
	createCloudCmd.Flags().String("project", "", "Project ID")
	createCloudCmd.Flags().String("region", "", "Region e.g. us-east-1")
	createCloudCmd.Flags().String("provider", "", "Provider e.g. aws")
	createCloudCmd.Flags().Bool("single-az", true, "Only deploy in one availability zone")
	createCloudCmd.Flags().String("image", "", "Custom image e.g. couchbase-cloud-server-7.2.0-1409-qe")
	createCloudCmd.Flags().String("override-token", "", "Override token to use non default deployment options")
	createCloudCmd.Flags().String("env", "", "Predefined environment (e.g. prod, stage, dev)")
	createCloudCmd.Flags().String("server", "", "Server version to use e.g. 7.1.3, Using an image will override this.")
}
