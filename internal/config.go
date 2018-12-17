package internal

import (
	"fmt"
	"os"

	"github.com/fananchong/go-xserver/common"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	configPath string
	app        string
)

func loadConfig() error {
	rootCmd := &cobra.Command{
		Use: "app",
		Run: func(c *cobra.Command, args []string) {
		},
	}
	flags := rootCmd.PersistentFlags()
	flags.StringVarP(&configPath, "config", "c", "../config", "config path name")
	flags.StringVarP(&app, "app", "a", "", "application name")
	viper.BindPFlags(rootCmd.PersistentFlags())
	cobra.OnInitialize(func() {
		viper.SetConfigFile(configPath + "/config.toml")
		viper.AutomaticEnv()
		if err := viper.ReadInConfig(); err != nil {
			fmt.Println("viper.ReadInConfig fail, err =", err)
			os.Exit(1)
		}
		if err := viper.Unmarshal(&common.XCONFIG); err != nil {
			fmt.Println("viper.Unmarshal fail, err =", err)
			os.Exit(1)
		}
		viper.WatchConfig()
		viper.OnConfigChange(func(e fsnotify.Event) {
			c := common.Config{}
			if err := viper.Unmarshal(&c); err != nil {
				common.XLOG.Errorln("viper.Unmarshal fail, err =", err)
			} else {
				if c.Common.Version != "" {
					common.XCONFIG = c
					common.XLOG.Printf("config: %#v\n", common.XCONFIG)
				}
			}
		})
	})

	return rootCmd.Execute()
}
