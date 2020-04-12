package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/couchbaselabs/cbdynclusterd/daemon"
	"github.com/spf13/cobra"
)

var certAuthCmd = &cobra.Command{
	Use:   "setup-cert-auth <cluster ID>",
	Short: "Enable client cert auth",
	Long:  "Enable and client cert auth for the username provided, writing cert files to the directory provided",
	Run: func(cmd *cobra.Command, args []string) {
		checkConfigInitialized()

		if len(args) == 0 {
			printAndExit("No cluster-id provided")
		}

		path := "/cluster/" + args[0] + "/setup-cert-auth"

		flags := cmd.Flags()
		userOption, err := flags.GetString("user")
		if userOption == "" {
			printAndExit("No user provided")
		}
		emailOption, err := flags.GetString("email")
		if emailOption == "" {
			emailOption = userOption + "@couchbase.com"
		}
		outOption, err := flags.GetString("out-dir")
		if outOption == "" {
			printAndExit("No out-dir provided")
		}

		var reqData daemon.SetupClientCertAuthJSON
		reqData.UserName = userOption
		reqData.UserEmail = emailOption

		var respData daemon.CertAuthResultJSON

		err = serverRestCall("POST", path, reqData, &respData, false)
		if err != nil {
			fmt.Printf("Failed to setup client cert auth: %s\n", err)
			os.Exit(1)
		}

		if err := writeCert(outOption+"/ca.pem", respData.CACert); err != nil {
			log.Fatalf("Failed to write CA cert file: %v", err)
		}

		if err := writeCert(outOption+"/client.pem", respData.ClientCert); err != nil {
			log.Fatalf("Failed to write client cert file: %v", err)
		}

		if err := writeCert(outOption+"/client.key", respData.ClientKey); err != nil {
			log.Fatalf("Failed to write client key file: %v", err)
		}
	},
}

func writeCert(path string, b []byte) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.Write(b); err != nil {
		return err
	}

	return nil
}

func init() {
	rootCmd.AddCommand(certAuthCmd)

	certAuthCmd.Flags().String("user", "", "The username to create a cert for")
	certAuthCmd.Flags().String("email", "", "The user email to create a cert for")
	certAuthCmd.Flags().String("out-dir", ".", "The directory to write certs to")
}
