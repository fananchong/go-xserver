package components

import (
	"os"

	"github.com/fananchong/go-xserver/common"
	nodenormal "github.com/fananchong/go-xserver/internal/components/node/normal"
)

// 服务节点通过该组件加入服务器组

// Node : 节点组件
type Node struct {
	ctx  *common.Context
	node *nodenormal.Node
}

// NewNode : 实例化
func NewNode(ctx *common.Context) *Node {
	node := &Node{ctx: ctx}
	switch getPluginType(node.ctx) {
	case common.Mgr:
		node.node = nil
	default:
		node.node = nodenormal.NewNode(node.ctx, pluginType)
		if node.node.Init() {
			node.ctx.Node = node.node
		}
	}
	return node
}

// Start : 实例化组件
func (node *Node) Start() bool {
	go func() {
		WaitComponent(node.ctx.Ctx)
		node.ctx.Log.Infoln("Service node start ...")
		if node.node != nil && node.node.Start() == false {
			node.ctx.Log.Errorln("Service node failed to start")
			os.Exit(1)
		}
	}()
	return true
}

// Close : 关闭组件
func (node *Node) Close() {
	if node.node != nil {
		node.node.Close()
		node.node = nil
	}
}
