package cdn

import (
	"fmt"

	// "github.com/coreos/etcd/tools/benchmark/cmd"
	// "github.com/coreos/etcd/tools/benchmark/cmd"

	"github.com/fleacloud/app"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// var backend
func initViper() error {
	if apiServerCfgFile != "" {

		apiConf.SetConfigFile(apiServerCfgFile)
	} else {

		apiConf.SetConfigName("apiserver")
		apiConf.AddConfigPath("/etc/cdn/")
		apiConf.AddConfigPath("$HOME/.cdn")
		apiConf.AddConfigPath("./conf/")
		apiConf.AddConfigPath(".")
	}

	if err := apiConf.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", apiConf.ConfigFileUsed())
		return err
	}
	return nil

}

// //全局配置
var apiConf = viper.New()
var apiServerCfgFile string

func apiserverCmdRunE(cmd *cobra.Command, args []string) error {
	initViper()

	if err := httpServer(); err != nil {
		return err
	}
	return nil
}

var cdnApiRootCmd = &cobra.Command{
	Use:     "cdnapi",
	Aliases: []string{"api", "a", "ca"},
	Short:   "start a cdn api server.",
	Long:    "start a cdn api server.",
	RunE:    apiserverCmdRunE,
}

func init() {
	cdnApiRootCmd.PersistentFlags().StringVar(&apiServerCfgFile, "config", "", "config file")
	app.AddCommond(cdnApiRootCmd)
}
