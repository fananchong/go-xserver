package components

import (
	"os"

	"github.com/fananchong/go-xserver/common"
	"github.com/fananchong/go-xserver/internal/components/misc"
	nodemgr "github.com/fananchong/go-xserver/internal/components/node/mgr"
)

// Mgr : 管理服务器
type Mgr struct {
	ctx  *common.Context
	node *nodemgr.Node
}

// NewMgr : 实例化
func NewMgr(ctx *common.Context) *Mgr {
	mgr := &Mgr{ctx: ctx}
	pluginType := misc.GetPluginType(mgr.ctx.Ctx)
	if pluginType == common.Mgr {
		mgr.node = nodemgr.NewNode(mgr.ctx)
		if mgr.node.Init() {
			mgr.ctx.Node = mgr.node
		}
	}
	return mgr
}

// Start : 实例化组件
func (mgr *Mgr) Start() bool {
	pluginType := misc.GetPluginType(mgr.ctx.Ctx)
	if pluginType == common.Mgr {
		if mgr.node.Start() == false {
			mgr.ctx.Log.Errorln("Mgr Server node failed to start")
			os.Exit(1)
		}
	}
	return true
}

// Close : 关闭组件
func (mgr *Mgr) Close() {
	if mgr.node != nil {
		mgr.node.Close()
		mgr.node = nil
	}
}
