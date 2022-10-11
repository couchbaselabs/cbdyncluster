package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var setupTrustedCertCmd = &cobra.Command{
	Use:   "setup-trusted-cert <cluster ID>",
	Short: "Setup trusted cert",
	Long:  "Setup trusted cert using cluster ID",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		checkConfigInitialized()

		path := "/cluster/" + args[0] + "/setup-trusted-cert"

		err := serverRestCall("POST", path, nil, nil, false)

		if err != nil {
			fmt.Printf("Failed to setup cluster encryption: %s\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(setupTrustedCertCmd)
}
