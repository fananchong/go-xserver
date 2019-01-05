package components

import (
	"os"
	"plugin"
	"sync"

	"github.com/fananchong/go-xserver/common"
	"github.com/spf13/viper"
)

var pluginObj common.Plugin
var pluginType common.NodeType
var once0 sync.Once
var once1 sync.Once

// Plugin : 插件组件
type Plugin struct {
}

// Start : 实例化组件
func (*Plugin) Start() bool {
	loadPlugin()
	if pluginObj != nil {
		return pluginObj.Start()
	}
	return false
}

// Close : 关闭组件
func (*Plugin) Close() {
	once1.Do(func() {
		if pluginObj != nil {
			pluginObj.Close()
			pluginObj = nil
		}
	})
}

func loadPlugin() {
	once0.Do(func() {
		appName := viper.GetString("app")
		if appName == "" {
			printUsage()
			os.Exit(1)
		}
		so, err := plugin.Open(appName + ".so")
		if err != nil {
			common.XLOG.Errorln(err)
			os.Exit(1)
		}
		obj, err := so.Lookup("PluginObj")
		if err != nil {
			common.XLOG.Errorln(err)
			os.Exit(1)
		}
		t, err := so.Lookup("PluginType")
		if err != nil {
			common.XLOG.Errorln(err)
			os.Exit(1)
		}
		pluginObj = *obj.(*common.Plugin)
		pluginType = *t.(*common.NodeType)
	})
}
