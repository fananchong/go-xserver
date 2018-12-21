package main

import (
	"fmt"

	"github.com/fananchong/go-xserver/common"
)

// PluginObj : 代表一个插件对象
var PluginObj common.Plugin

func init() {
	fmt.Println("LOAD PLUGIN: MGR")
	PluginObj = &Plugin{}
}

// Plugin : 插件类
type Plugin struct {
}

// Init : 插件类实现初始化
func (plugin *Plugin) Init() (nodeType common.NodeType, ok bool) {
	common.XLOG.Infoln("Plugin Init")
	nodeType = common.Mgr
	ok = true
	return
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
