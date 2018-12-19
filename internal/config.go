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
	"github.com/spf13/viper"
)

var (
	configPath string
	app        string
	ENV_PREFIX string = "GOXSERVER"
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
		fmt.Println("ENV:")
		keys := viper.AllKeys()
		sort.Sort(sort.StringSlice(keys))
		for _, k := range keys {
			env := strings.ToUpper(strings.Replace(ENV_PREFIX+"_"+k, ".", "_", -1))
			fmt.Printf("    %s\n", env)
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

	rt := reflect.TypeOf(s)
	for i := 0; i < rt.NumField(); i++ {
		sf := rt.Field(i)
		name := strings.ToLower(sf.Name)
		if sf.Type.Kind() != reflect.Struct {
			fmt.Println("bindConfig fail")
			os.Exit(1)
		}
		srt := sf.Type
		for j := 0; j < srt.NumField(); j++ {
			ssf := srt.Field(j)
			sname := fmt.Sprintf("%s-%s", name, strings.ToLower(ssf.Name))
			switch ssf.Type.Kind() {
			case reflect.Bool:
				flags.Bool(sname, false, "")
			case reflect.String:
				flags.String(sname, "", "")
			case reflect.Int:
				flags.Int(sname, 0, "")
			case reflect.Uint:
				flags.Uint(sname, 0, "")
			case reflect.Slice:
				switch sssf := ssf.Type.Elem(); sssf.Kind() {
				case reflect.Bool:
					flags.BoolSlice(sname, []bool{}, "")
				case reflect.String:
					flags.StringSlice(sname, []string{}, "")
				case reflect.Int:
					flags.IntSlice(sname, []int{}, "")
				case reflect.Uint:
					flags.UintSlice(sname, []uint{}, "")
				default:
					fmt.Println("bindConfig fail, err = unsupported field: ", sname)
					os.Exit(1)
				}
			default:
				fmt.Println("bindConfig fail, err = unsupported field: ", sname)
				os.Exit(1)
			}
			viper.BindPFlag(fmt.Sprintf("%s.%s", name, ssf.Name), flags.Lookup(sname))
			viper.BindEnv(fmt.Sprintf("%s.%s", name, ssf.Name), fmt.Sprintf("%s_%s_%s", ENV_PREFIX, strings.ToUpper(name), strings.ToUpper(ssf.Name)))
		}
	}
}
