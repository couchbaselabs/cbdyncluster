package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var dockerhostCmd = &cobra.Command{
	Use:   "docker-host",
	Short: "Fetches IP of docker host",
	Run: func(cmd *cobra.Command, args []string) {
		checkConfigInitialized()

		dockerHost, err := getDockerHost()
		if err != nil {
			fmt.Printf("Can not get docker host:%s\n", err)
			return
		}
		fmt.Printf("%s\n", dockerHost)
	},
}

func init() {
	rootCmd.AddCommand(dockerhostCmd)
}
