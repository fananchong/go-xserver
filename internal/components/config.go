package components

import (
	"fmt"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/fananchong/go-xserver/common"
	"github.com/fananchong/go-xserver/common/config"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Config : 配置组件
type Config struct {
	ctx        *common.Context
	configPath string
	app        string
	suffix     string
	rootCmd    *cobra.Command
	viperObj   *viper.Viper
	configObj  *config.FrameworkConfig
}

// NewConfig : 实例化
func NewConfig(ctx *common.Context) *Config {
	cfg := &Config{
		ctx:       ctx,
		viperObj:  viper.New(),
		configObj: &config.FrameworkConfig{},
	}
	if err := cfg.init(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	ctx.IConfig = cfg
	return cfg
}

// Start : 实例化组件
func (*Config) Start() bool {
	return true
}

// Close : 关闭组件
func (*Config) Close() {
	// No need to do anything
}

const constEnvPrefix string = "GOXSERVER"

func (cfg *Config) init() error {
	cfg.rootCmd = &cobra.Command{
		Use: "go-xserver",
		Run: func(c *cobra.Command, args []string) {
			if appName := cfg.viperObj.GetString("app"); appName == "" {
				cfg.PrintUsage()
				os.Exit(1)
			}
		},
	}
	cfg.rootCmd.SetHelpFunc(func(c *cobra.Command, args []string) {
		cfg.PrintUsage()
		os.Exit(1)
	})
	flags := cfg.rootCmd.PersistentFlags()
	flags.StringVarP(&cfg.configPath, "config", "c", "../config", "配置目录路径")
	flags.StringVarP(&cfg.app, "app", "a", "", "应用名（插件，必填）")
	flags.StringVarP(&cfg.suffix, "suffix", "s", "", "Log 文件名后缀，多开时可以通过它，避免多个进程共用 1 个 Log 文件")
	cfg.viperObj.BindPFlags(cfg.rootCmd.PersistentFlags())
	cfg.bindConfig(cfg.rootCmd, *cfg.configObj)
	cobra.OnInitialize(func() {
		cfg.viperObj.SetConfigFile(cfg.configPath + "/framework.toml")
		cfg.viperObj.AutomaticEnv()
		if err := cfg.viperObj.ReadInConfig(); err != nil {
			fmt.Println("Failed to read configuration file, err =", err)
			os.Exit(1)
		}
		if err := cfg.viperObj.Unmarshal(cfg.configObj); err != nil {
			fmt.Println("Parsing the configuration file failed, err =", err)
			os.Exit(1)
		}
		cfg.viperObj.WatchConfig()
		cfg.viperObj.OnConfigChange(func(e fsnotify.Event) {
			c := &config.FrameworkConfig{}
			if err := cfg.viperObj.Unmarshal(c); err != nil {
				cfg.ctx.Errorln("Parsing the configuration file failed, err =", err)
			} else {
				if c.Common.Version != "" {
					cfg.configObj = c
					cfg.ctx.Printf("Configuration information is: %#v\n", cfg.configObj)
				}
			}
		})
	})
	return cfg.rootCmd.Execute()
}

func (cfg *Config) bindConfig(c *cobra.Command, s interface{}) {
	flags := c.Flags()
	cfg.parseStruct(flags, reflect.TypeOf(s), "")
}

func (cfg *Config) parseStruct(flags *pflag.FlagSet, rt reflect.Type, prefix string) {
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
			cfg.parseStruct(flags, sf.Type, rawName)
			continue
		case reflect.Bool:
			v, _ := strconv.ParseBool(defaultValue)
			flags.Bool(name, v, desc)
		case reflect.String:
			flags.String(name, defaultValue, desc)
		case reflect.Int, reflect.Int32:
			v, _ := strconv.ParseInt(defaultValue, 10, 32)
			flags.Int(name, int(v), desc)
		case reflect.Uint, reflect.Uint32:
			v, _ := strconv.ParseUint(defaultValue, 10, 32)
			flags.Uint(name, uint(v), desc)
		case reflect.Int64:
			v, _ := strconv.ParseInt(defaultValue, 10, 64)
			flags.Int64(name, int64(v), desc)
		case reflect.Uint64:
			v, _ := strconv.ParseUint(defaultValue, 10, 64)
			flags.Uint64(name, uint64(v), desc)
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
		cfg.viperObj.BindPFlag(rawName, flags.Lookup(name))
		cfg.viperObj.BindEnv(rawName, strings.Replace(fmt.Sprintf("%s_%s", constEnvPrefix, strings.ToUpper(name)), "-", "_", -1))
	}
}

// LoadConfig : 逻辑层可以用该接口加载配置文件
func (cfg *Config) LoadConfig(cfgfile string, cfgobj interface{}) bool {
	// TODO:
	return false
}

// Config : 获取框架层配置
func (cfg *Config) Config() *config.FrameworkConfig {
	return cfg.configObj
}

// PrintUsage : 打印命令行参数
func (cfg *Config) PrintUsage() {
	cfg.rootCmd.Usage()
	fmt.Println("")
	fmt.Println("Environment variables:")
	keys := cfg.viperObj.AllKeys()
	sort.Sort(sort.StringSlice(keys))
	for _, k := range keys {
		env := strings.ToUpper(strings.Replace(constEnvPrefix+"_"+k, ".", "_", -1))
		fmt.Printf("   %s\n", env)
	}
}

// GetViperObj : 获取 viper 对象
func (cfg *Config) GetViperObj() *viper.Viper {
	return cfg.viperObj
}
