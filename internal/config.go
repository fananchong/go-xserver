package internal

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fananchong/go-xserver/common"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	_ = flag.String("assets", "", "path of assets")
	_ = flag.String("app", "", "app name")
)

func loadConfig() error {
	viper.SetConfigName("config")
	viper.AddConfigPath(getAssetsPath())
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
	viper.AutomaticEnv()
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		common.XLOG.Infoln("Config file changed:", e.Name)
	})
	return nil
}

func getAssetsPath() string {
	var path string
	for i, v := range os.Args {
		if v == "-assets" || v == "--assets" {
			path = os.Args[i+1] + "/"
			break
		}
		if strings.Contains(v, "-assets=") {
			path = os.Args[i]
			path = strings.Replace(path, "--assets=", "", -1)
			path = strings.Replace(path, "-assets=", "", -1)
			path = path + "/"
			break
		}
	}
	if path == "" {
		path = "./"
	}
	dir, err := filepath.Abs(filepath.Dir(path))
	if err != nil {
		fmt.Println("no find assets path, path: " + path)
		return path
	}
	fmt.Println("Assets Path:", dir)
	path = dir + "/"
	return path
}
