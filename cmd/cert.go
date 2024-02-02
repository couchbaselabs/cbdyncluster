package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var cert = &cobra.Command{
	Use:   "cert <cluster ID>",
	Short: "Get the trusted cert on the cluster",
	Long:  "Get the trusted cert using cluster ID",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		checkConfigInitialized()

		path := "/cluster/" + args[0] + "/certificate"

		var pem string
		err := serverRestCall("GET", path, nil, &pem, false)

		if err != nil {
			fmt.Printf("Failed to get cluster cert: %s\n", err)
			os.Exit(1)
		}
		fmt.Printf("%s\n", pem)
	},
}

func init() {
	rootCmd.AddCommand(cert)
}
