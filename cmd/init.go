package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pelletier/go-toml"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Configures the client",
	Long: `Configues the client with basic information such as the
server to communicate with and the users email for ownership
tracking purposes.`,
	Run: func(cmd *cobra.Command, args []string) {
		configFile := viper.ConfigFileUsed()
		if configFile == "" {
			fmt.Printf("Unable to determine appropriate config file to use")
			os.Exit(1)
		}

		fmt.Println("Using config file:", configFile)

		tmap, err := toml.LoadFile(configFile)
		if err != nil {
			tmap, _ = toml.TreeFromMap(nil)
		}

		if serverFlag != "" {
			fmt.Printf("Server (%s): ", serverFlag)
		} else {
			fmt.Printf("Server: ")
		}
		var serverStr string
		fmt.Scanf("%s", &serverStr)
		if serverStr != "" {
			serverFlag = serverStr
		}
		tmap.Set("server", serverFlag)

		if userFlag != "" {
			fmt.Printf("User Email (%s): ", userFlag)
		} else {
			fmt.Printf("User Email: ")
		}
		var userStr string
		fmt.Scanf("%s", &userStr)
		if userStr != "" {
			userFlag = userStr
		}
		tmap.Set("userEmail", userFlag)

		s := tmap.String()
		b := []byte(s)

		err = ioutil.WriteFile(configFile, b, 0644)
		if err != nil {
			fmt.Printf("Failed to write config file `%s`: %s", configFile, err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
