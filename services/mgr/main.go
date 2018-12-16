package main

import (
	"fmt"

	"github.com/fananchong/go-xserver/common"
)

// PluginObj : Exported to internal.App
var PluginObj common.Plugin

func init() {
	fmt.Println("LOAD PLUGIN: MGR")
	PluginObj = &Plugin{}
}

// Plugin : Plugin class implementation
type Plugin struct {
}

// Init : Plugin class function Init implementation
func (plugin *Plugin) Init() bool {
	common.XLOG.Infoln("Plugin Init")
	return true
}

// Start : Plugin class function Start implementation
func (plugin *Plugin) Start() bool {
	common.XLOG.Infoln("Plugin Start")
	return true
}

// Close : Plugin class function Close implementation
func (plugin *Plugin) Close() {
	common.XLOG.Infoln("Plugin Close")
}
