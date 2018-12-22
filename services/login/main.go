package main

import (
	"fmt"

	"github.com/fananchong/go-xserver/common"
)

// PluginObj : 代表一个插件对象
var PluginObj common.Plugin

func init() {
	fmt.Println("LOAD PLUGIN: LOGIN")
	PluginObj = &Plugin{}
}

// Plugin : 插件类
type Plugin struct {
}

// Init : 插件类实现初始化
func (plugin *Plugin) Init() (nodeType common.NodeType, ok bool) {
	common.XLOG.Infoln("Plugin Init")
	nodeType = common.Login
	ok = true
	return
}

// RegisterCallBack : 插件类实现启动前处理，注册自定义回调
func (plugin *Plugin) RegisterCallBack() {
	common.XLOG.Infoln("Plugin RegisterCallBack")
}

// Start : 插件类实现启动
func (plugin *Plugin) Start() bool {
	common.XLOG.Infoln("Plugin Start")
	return true
}

// Close : 插件类实现关闭
func (plugin *Plugin) Close() {
	common.XLOG.Infoln("Plugin Close")
}

// main : 作为插件包，该函数可以不存在。添加之，是避免 go-lint 烦人的错误提示
func main() {}
