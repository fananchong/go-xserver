package components

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

// Config : 配置组件
type Config struct {
	ctx *common.Context
}

// NewConfig : 实例化
func NewConfig(ctx *common.Context) *Config {
	cfg := &Config{ctx: ctx}
	if err := loadConfig(ctx); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return cfg
}

// Start : 实例化组件
func (confing *Config) Start() bool {
	return true
}

// Close : 关闭组件
func (*Config) Close() {
	// No need to do anything
}

var (
	configPath string
	app        string
	suffix     string
	rootCmd    *cobra.Command
)

const (
	constEnvPrefix string = "GOXSERVER"
)

func loadConfig(ctx *common.Context) error {
	rootCmd = &cobra.Command{
		Use: "go-xserver",
		Run: func(c *cobra.Command, args []string) {
			if appName := viper.GetString("app"); appName == "" {
				printUsage()
				os.Exit(1)
			}
		},
	}
	rootCmd.SetHelpFunc(func(c *cobra.Command, args []string) {
		printUsage()
		os.Exit(1)
	})
	flags := rootCmd.PersistentFlags()
	flags.StringVarP(&configPath, "config", "c", "../config", "配置目录路径")
	flags.StringVarP(&app, "app", "a", "", "应用名（插件，必填）")
	flags.StringVarP(&suffix, "suffix", "s", "", "Log 文件名后缀，多开时可以通过它，避免多个进程共用 1 个 Log 文件")
	viper.BindPFlags(rootCmd.PersistentFlags())
	bindConfig(rootCmd, common.Config{})
	cobra.OnInitialize(func() {
		viper.SetConfigFile(configPath + "/config.toml")
		viper.AutomaticEnv()
		if err := viper.ReadInConfig(); err != nil {
			fmt.Println("Failed to read configuration file, err =", err)
			os.Exit(1)
		}
		ctx.Config = &common.Config{}
		if err := viper.Unmarshal(ctx.Config); err != nil {
			fmt.Println("Parsing the configuration file failed, err =", err)
			os.Exit(1)
		}
		viper.WatchConfig()
		viper.OnConfigChange(func(e fsnotify.Event) {
			c := common.Config{}
			if err := viper.Unmarshal(&c); err != nil {
				ctx.Log.Errorln("Parsing the configuration file failed, err =", err)
			} else {
				if c.Common.Version != "" {
					ctx.Config = &c
					ctx.Log.Printf("Configuration information is: %#v\n", ctx.Config)
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
			var v []string
			if defaultValue != "" {
				for _, tmp := range strings.Split(strings.Trim(defaultValue, "[]"), ",") {
					v = append(v, strings.Trim(tmp, " "))
				}
			}
			flags.StringSlice(name, v, desc)
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
