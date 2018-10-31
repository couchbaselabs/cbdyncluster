package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const DOCKER_PORT = "2376"

var dockerhostCmd = &cobra.Command{
	Use:   "docker-host",
	Short: "Fetches IP of docker host",
	Run: func(cmd *cobra.Command, args []string) {
		checkConfigInitialized()
		dockerHost := viper.GetString("server")
		if dockerHost != "" {
			parsed := strings.Split(dockerHost, ":")
			dockerHost = parsed[0]+":"+DOCKER_PORT
		}
		fmt.Printf("%s\n", dockerHost)
	},
}

func init() {
	rootCmd.AddCommand(dockerhostCmd)
}
