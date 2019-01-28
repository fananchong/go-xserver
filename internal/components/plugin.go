package components

import (
	"os"
	"plugin"

	"github.com/fananchong/go-xserver/common"
	"github.com/spf13/viper"
)

var pluginObj common.Plugin
var pluginType common.NodeType

// Plugin : 插件组件
type Plugin struct {
	ctx *common.Context
}

// NewPlugin : 实例化
func NewPlugin(ctx *common.Context) *Plugin {
	p := &Plugin{ctx: ctx}
	loadPlugin(ctx)
	return p
}

// Start : 实例化组件
func (p *Plugin) Start() bool {
	if pluginObj != nil {
		return pluginObj.Start()
	}
	return false
}

// Close : 关闭组件
func (*Plugin) Close() {
	if pluginObj != nil {
		pluginObj.Close()
		pluginObj = nil
	}
}

func loadPlugin(ctx *common.Context) {
	appName := viper.GetString("app")
	if appName == "" {
		printUsage()
		os.Exit(1)
	}
	so, err := plugin.Open(appName + ".so")
	if err != nil {
		ctx.Log.Errorln(err)
		os.Exit(1)
	}
	obj, err := so.Lookup("PluginObj")
	if err != nil {
		ctx.Log.Errorln(err)
		os.Exit(1)
	}
	t, err := so.Lookup("PluginType")
	if err != nil {
		ctx.Log.Errorln(err)
		os.Exit(1)
	}
	c, err := so.Lookup("Ctx")
	if err != nil {
		ctx.Log.Errorln(err)
		os.Exit(1)
	}
	pluginObj = *obj.(*common.Plugin)
	pluginType = *t.(*common.NodeType)
	*c.(**common.Context) = ctx
}

func getPluginType(ctx *common.Context) common.NodeType {
	return pluginType
}
