package internal

import (
	"fmt"
	"os"
	"reflect"
	"sort"
	"strconv"
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
	rootCmd    *cobra.Command
)

const (
	constEnvPrefix string = "GOXSERVER"
)

func loadConfig() error {
	rootCmd = &cobra.Command{
		Use: "go-xserver",
		Run: func(c *cobra.Command, args []string) {
		},
	}
	rootCmd.SetHelpFunc(func(c *cobra.Command, args []string) {
		printUsage()
		os.Exit(1)
	})
	flags := rootCmd.PersistentFlags()
	flags.StringVarP(&configPath, "config", "c", "../config", "配置目录路径")
	flags.StringVarP(&app, "app", "a", "", "应用名（插件，必填）")
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
		desc := sf.Tag.Get("desc")
		defaultValue := sf.Tag.Get("default")
		switch sf.Type.Kind() {
		case reflect.Struct:
			parseStruct(flags, sf.Type, rawName)
			continue
		case reflect.Bool:
			v, _ := strconv.ParseBool(defaultValue)
			flags.Bool(name, v, desc)
		case reflect.String:
			flags.String(name, defaultValue, desc)
		case reflect.Int:
			v, _ := strconv.ParseInt(defaultValue, 10, 32)
			flags.Int(name, int(v), desc)
		case reflect.Uint:
			v, _ := strconv.ParseUint(defaultValue, 10, 32)
			flags.Uint(name, uint(v), desc)
		case reflect.Slice:
			switch ssf := sf.Type.Elem(); ssf.Kind() {
			case reflect.Bool:
				var v []bool
				if defaultValue != "" {
					for _, tmp := range strings.Split(strings.Trim(defaultValue, "[]"), ",") {
						tmp2, _ := strconv.ParseBool(strings.Trim(tmp, " "))
						v = append(v, tmp2)
					}
				}
				flags.BoolSlice(name, v, desc)
			case reflect.String:
				var v []string
				if defaultValue != "" {
					for _, tmp := range strings.Split(strings.Trim(defaultValue, "[]"), ",") {
						v = append(v, strings.Trim(tmp, " "))
					}
				}
				flags.StringSlice(name, v, desc)
			case reflect.Int, reflect.Int32:
				var v []int
				if defaultValue != "" {
					for _, tmp := range strings.Split(strings.Trim(defaultValue, "[]"), ",") {
						tmp2, _ := strconv.ParseInt(strings.Trim(tmp, " "), 10, 32)
						v = append(v, int(tmp2))
					}
				}
				flags.IntSlice(name, v, desc)
			case reflect.Uint, reflect.Uint32:
				var v []uint
				if defaultValue != "" {
					for _, tmp := range strings.Split(strings.Trim(defaultValue, "[]"), ",") {
						tmp2, _ := strconv.ParseUint(strings.Trim(tmp, " "), 10, 32)
						v = append(v, uint(tmp2))
					}
				}
				flags.UintSlice(name, v, desc)
			default:
				fmt.Println("bindConfig fail, err = unsupported field: ", rawName)
				os.Exit(1)
			}
		default:
			fmt.Println("bindConfig fail, err = unsupported field: ", rawName)
			os.Exit(1)
		}
		viper.BindPFlag(rawName, flags.Lookup(name))
		viper.BindEnv(rawName, strings.Replace(fmt.Sprintf("%s_%s", constEnvPrefix, strings.ToUpper(name)), "-", "_", -1))
	}
}

func printUsage() {
	rootCmd.Usage()
	fmt.Println("")
	fmt.Println("Environment variables:")
	keys := viper.AllKeys()
	sort.Sort(sort.StringSlice(keys))
	for _, k := range keys {
		env := strings.ToUpper(strings.Replace(constEnvPrefix+"_"+k, ".", "_", -1))
		fmt.Printf("   %s\n", env)
	}
}
