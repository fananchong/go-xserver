package main

import (
	"fmt"

	"github.com/fananchong/go-xserver/common"
)

// PluginObj : 代表一个插件对象
var PluginObj common.Plugin

// PluginType : 插件类型
var PluginType common.NodeType

var login *Login

func init() {
	fmt.Println("LOAD PLUGIN: LOGIN")
	PluginObj = &Plugin{}
	PluginType = common.Login
	login = NewLogin()
	login.Init()
}

// Plugin : 插件类
type Plugin struct {
}

// Start : 插件类实现启动
func (plugin *Plugin) Start() bool {
	common.XLOG.Infoln("Plugin Start")
	login.Start()
	return true
}

// Close : 插件类实现关闭
func (plugin *Plugin) Close() {
	common.XLOG.Infoln("Plugin Close")
	login.Close()
}

// main : 作为插件包，该函数可以不存在。添加之，是避免 go-lint 烦人的错误提示
func main() {}
