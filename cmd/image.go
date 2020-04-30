package cmd

import (
	"fmt"
	"os"

	"github.com/couchbaselabs/cbdynclusterd/daemon"
	"github.com/spf13/cobra"
)

var buildImageCmd = &cobra.Command{
	Use:   "build-image",
	Short: "Builds an image",
	//Long:  `Allocates a new cluster`,
	Run: func(cmd *cobra.Command, args []string) {
		checkConfigInitialized()

		serverVersion, _ := cmd.Flags().GetString("server-version")

		reqData := daemon.BuildImageJSON{
			ServerVersion: serverVersion,
		}

		var respData daemon.BuildImageResponseJSON
		err := serverRestCall("POST", "/images", reqData, &respData, false)
		if err != nil {
			fmt.Printf("Failed to allocate cluster: %s\n", err)
			os.Exit(1)
		}

		fmt.Printf("%s\n", respData.ImageName)
	},
}

func init() {
	rootCmd.AddCommand(buildImageCmd)

	buildImageCmd.Flags().String("server-version", "5.5.0", "The server version to use when allocating the nodes.")
}
