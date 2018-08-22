package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/couchbaselabs/cbdynclusterd/daemon"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var serverFlag string
var userFlag string
var forceFlag bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cbdyncluster",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cbdyncluster.toml)")
	rootCmd.PersistentFlags().StringVar(&serverFlag, "server", "dyncluster.hq.couchbase.com", "the server to operate on")
	rootCmd.PersistentFlags().StringVar(&userFlag, "user", "", "the user to operate as (as a @couchbase email)")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		configFile := path.Join(home, ".cbdyncluster.toml")
		viper.SetConfigFile(configFile)
	}

	viper.AutomaticEnv()
	viper.ReadInConfig()

	cfgServer := viper.GetString("server")
	if cfgServer != "" {
		serverFlag = cfgServer
	}

	cfgUser := viper.GetString("userEmail")
	if cfgUser != "" {
		userFlag = cfgUser
	}
}

func checkConfigInitialized() {
	if serverFlag == "" || userFlag == "" {
		fmt.Println("You must configure cbdyncluster before you can use it.")
		fmt.Println("  Please run `cbdyncluster init` before any other commands")
		os.Exit(1)
	}
}

var httpClient http.Client

func serverRestCall(method, path string, data interface{}, dataOut interface{}, asAdmin bool) error {
	url := fmt.Sprintf("http://%s%s", serverFlag, path)

	var body io.Reader
	if data != nil {
		jsonBytes, err := json.Marshal(data)
		if err != nil {
			return err
		}

		body = bytes.NewReader(jsonBytes)
	}

	req, err := http.NewRequest(method, url, body)
	req.Header.Set("cbdn-user", userFlag)
	if asAdmin {
		req.Header.Set("cbdn-admin", "true")
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var jsonError daemon.ErrorJSON
		jsonDec := json.NewDecoder(resp.Body)
		err = jsonDec.Decode(&jsonError)
		if err != nil {
			return err
		}

		return errors.New(jsonError.Error.Message)
	}

	if dataOut != nil {
		jsonDec := json.NewDecoder(resp.Body)
		err = jsonDec.Decode(dataOut)
		if err != nil {
			return err
		}
	}

	return nil
}

func getAllClusters() ([]*daemon.Cluster, error) {
	var respData daemon.GetClustersJSON
	err := serverRestCall("GET", "/clusters", nil, &respData, psAllFlag)
	if err != nil {
		return nil, err
	}

	var clusters []*daemon.Cluster
	for _, jsonCluster := range respData {
		cluster, err := daemon.UnjsonifyCluster(&jsonCluster)
		if err != nil {
			return nil, err
		}

		clusters = append(clusters, cluster)
	}

	return clusters, nil
}

func getCluster(clusterID string) (*daemon.Cluster, error) {
	var respData daemon.ClusterJSON
	path := fmt.Sprintf("/cluster/%s", clusterID)
	err := serverRestCall("GET", path, nil, &respData, psAllFlag)
	if err != nil {
		return nil, err
	}

	cluster, err := daemon.UnjsonifyCluster(&respData)
	if err != nil {
		return nil, err
	}

	return cluster, nil
}
