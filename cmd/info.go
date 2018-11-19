package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var Version = "0.0.0"

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Displays cbdyncluster and daemon info",
	Run: func(cmd *cobra.Command, args []string) {
		daemonVersion, _ := getDaemonVersion()
		fmt.Printf("cbdyncluster:%s\n", Version)
		fmt.Printf("cbdynclusterd:%s\n", daemonVersion)
		fmt.Printf("cbdynclusterd host:%s\n", serverFlag)
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}

