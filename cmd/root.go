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

	"github.com/couchbaselabs/cbdynclusterd/cluster"
	"github.com/couchbaselabs/cbdynclusterd/daemon"
	"github.com/golang/glog"
	"github.com/mitchellh/go-homedir"
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
	rootCmd.PersistentFlags().StringVar(&userFlag, "auth", "", "the user to operate as (as a @couchbase email)")
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

	serverFlag = getArg("server")
	userFlag = getArg("auth")
}

func getArg(arg string) string {
	var val string
	if rootCmd.PersistentFlags().Changed(arg) {
		// read from commandline option
		val, _ = rootCmd.PersistentFlags().GetString(arg)
	} else {
		// read from config
		val = viper.GetString(arg)
	}

	return val
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
		glog.Errorf("Rest call error:%s", err)
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		if resp.StatusCode == 404 {
			// We're not going to be able to decode this
			return fmt.Errorf("endpoint %s does not exist", path)
		}
		var jsonError daemon.ErrorJSON
		jsonDec := json.NewDecoder(resp.Body)
		err = jsonDec.Decode(&jsonError)
		if err != nil {
			glog.Errorf("Error on json decode of Rest call response:%s", err)
			return err
		}

		glog.Errorf("Error occurred executing command on cbdynclusterd: %s", jsonError.Error.Message)
		return errors.New(jsonError.Error.Message)
	}

	if dataOut != nil {
		jsonDec := json.NewDecoder(resp.Body)
		err = jsonDec.Decode(dataOut)
		if err != nil {
			glog.Errorf("Error on json decode of Rest call response body:%s", err)
			return err
		}
	}

	return nil
}

func getAllClusters() ([]*cluster.Cluster, error) {
	var respData daemon.GetClustersJSON
	err := serverRestCall("GET", "/clusters", nil, &respData, psAllFlag)
	if err != nil {
		return nil, err
	}

	var clusters []*cluster.Cluster
	for _, jsonCluster := range respData {
		cluster, err := daemon.UnjsonifyCluster(&jsonCluster)
		if err != nil {
			return nil, err
		}

		clusters = append(clusters, cluster)
	}

	return clusters, nil
}

func getCluster(clusterID string) (*cluster.Cluster, error) {
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

func getDockerHost() (string, error) {
	var respData daemon.DockerHostJSON
	path := fmt.Sprintf("/docker-host")
	err := serverRestCall("GET", path, nil, &respData, false)
	if err != nil {
		return "", err
	}
	dockerHost, err := daemon.UnjsonifyDockerHost(&respData)
	if err != nil {
		return "", err
	}

	return dockerHost, err
}

func getDaemonVersion() (string, error) {
	var respData daemon.VersionJSON
	path := fmt.Sprintf("/version")
	err := serverRestCall("GET", path, nil, &respData, false)
	if err != nil {
		return "", err
	}
	version, err := daemon.UnjsonifyVersion(&respData)
	if err != nil {
		return "", err
	}
	return version, err
}

func getConnString(clusterID string, useSSL, useSrv bool) (string, error) {
	var respData daemon.ConnStringResponseJSON
	path := fmt.Sprintf("/cluster/%s/connstr", clusterID)
	reqData := daemon.ConnStringJSON{UseSSL: useSSL, UseSrv: useSrv}
	err := serverRestCall("GET", path, reqData, &respData, false)
	if err != nil {
		return "", err
	}
	return respData.ConnStr, err
}

func writeBytes(path string, b []byte) error {
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
