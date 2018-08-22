package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var rmForceFlag bool

var rmCmd = &cobra.Command{
	Use:   "rm [cluster_id]",
	Short: "Deallocates a cluster",
	//Long:  `Deallocates a new cluster`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		checkConfigInitialized()

		path := fmt.Sprintf("/cluster/%s", args[0])
		err := serverRestCall("DELETE", path, nil, nil, rmForceFlag)
		if err != nil {
			log.Printf("Failed to remove cluster: %s\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)

	rmCmd.Flags().BoolVarP(&rmForceFlag, "force", "f", false, "allow removal of non-owned clusters")
}
