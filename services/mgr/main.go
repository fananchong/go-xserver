package main

import (
	"fmt"

	"github.com/fananchong/go-xserver/common"
)

var PluginObj common.Plugin

func init() {
	fmt.Println("LOAD PLUGIN: MGR")
	PluginObj = &PluginMgr{}
}

// PluginMgr :
type PluginMgr struct {
}

func (plugin *PluginMgr) Init() bool {
	common.XLOG.Infoln("Plugin Init")
	return true
}

func (plugin *PluginMgr) Start() bool {
	common.XLOG.Infoln("Plugin Start")
	return true
}

func (plugin *PluginMgr) Close() {
	common.XLOG.Infoln("Plugin Close")
}
