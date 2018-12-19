package internal

import (
	"fmt"
	"os"
	"reflect"
	"sort"
	"strings"

	"github.com/fananchong/go-xserver/common"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	configPath string
	app        string
)

const (
	// ConstEnvPrefix : you can tell Viper to use a prefix while reading from the environment variables
	ConstEnvPrefix string = "GOXSERVER"
)

func loadConfig() error {
	rootCmd := &cobra.Command{
		Use: "go-xserver",
		Run: func(c *cobra.Command, args []string) {
		},
	}
	rootCmd.SetHelpFunc(func(c *cobra.Command, args []string) {
		c.Usage()
		fmt.Println("")
		fmt.Println("Environment variables:")
		keys := viper.AllKeys()
		sort.Sort(sort.StringSlice(keys))
		for _, k := range keys {
			env := strings.ToUpper(strings.Replace(ConstEnvPrefix+"_"+k, ".", "_", -1))
			fmt.Printf("   %s\n", env)
		}
		os.Exit(1)
	})
	flags := rootCmd.PersistentFlags()
	flags.StringVarP(&configPath, "config", "c", "../config", "config path name")
	flags.StringVarP(&app, "app", "a", "", "application name")
	viper.BindPFlags(rootCmd.PersistentFlags())
	bindConfig(rootCmd, common.Config{})
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

func bindConfig(c *cobra.Command, s interface{}) {
	flags := c.Flags()
	parseStruct(flags, reflect.TypeOf(s), "")
}

func parseStruct(flags *pflag.FlagSet, rt reflect.Type, prefix string) {
	for i := 0; i < rt.NumField(); i++ {
		sf := rt.Field(i)
		rawName := strings.ToLower(sf.Name)
		if prefix != "" {
			rawName = prefix + "." + rawName
		}
		name := strings.Replace(rawName, ".", "-", -1)
		switch sf.Type.Kind() {
		case reflect.Struct:
			parseStruct(flags, sf.Type, rawName)
			continue
		case reflect.Bool:
			flags.Bool(name, false, "")
		case reflect.String:
			flags.String(name, "", "")
		case reflect.Int:
			flags.Int(name, 0, "")
		case reflect.Uint:
			flags.Uint(name, 0, "")
		case reflect.Slice:
			switch ssf := sf.Type.Elem(); ssf.Kind() {
			case reflect.Bool:
				flags.BoolSlice(name, []bool{}, "")
			case reflect.String:
				flags.StringSlice(name, []string{}, "")
			case reflect.Int:
				flags.IntSlice(name, []int{}, "")
			case reflect.Uint:
				flags.UintSlice(name, []uint{}, "")
			default:
				fmt.Println("bindConfig fail, err = unsupported field: ", rawName)
				os.Exit(1)
			}
		default:
			fmt.Println("bindConfig fail, err = unsupported field: ", rawName)
			os.Exit(1)
		}
		viper.BindPFlag(rawName, flags.Lookup(name))
		viper.BindEnv(rawName, strings.Replace(fmt.Sprintf("%s_%s", ConstEnvPrefix, strings.ToUpper(name)), "-", "_", -1))
	}
}
