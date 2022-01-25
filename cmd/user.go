package cmd

import (
	"fmt"
	"os"

	"github.com/couchbaselabs/cbdynclusterd/daemon"
	"github.com/couchbaselabs/cbdynclusterd/helper"
	"github.com/spf13/cobra"
)

var addUserCmd = &cobra.Command{
	Use:   "add-user <cluster ID>",
	Short: "Adds a user",
	Long:  `Adds a user to a cluster using cluster ID`,
	Run: func(cmd *cobra.Command, args []string) {
		checkConfigInitialized()

		if len(args) == 0 {
			printAndExit("No cluster-id provided")
		}

		path := "/cluster/" + args[0] + "/add-user"

		flags := cmd.Flags()
		var err error
		var name, password string

		if name, err = flags.GetString("name"); err != nil {
			printAndExit("Invalid name")
		}
		if password, err = flags.GetString("password"); err != nil {
			printAndExit("Invalid password")
		}

		var reqData daemon.AddUserJSON

		reqData.User = &helper.UserOption{
			Name:     name,
			Password: password,
			Roles:    &[]string{},
		}

		err = serverRestCall("POST", path, reqData, nil, false)

		if err != nil {
			fmt.Printf("Failed to add user: %s\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(addUserCmd)

	addUserCmd.Flags().String("name", "", "name of the user.")
	addUserCmd.Flags().String("password", "", "password of the user.")

}
