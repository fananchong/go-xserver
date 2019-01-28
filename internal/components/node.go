package components

import (
	"github.com/fananchong/go-xserver/common"
	nodemgr "github.com/fananchong/go-xserver/internal/node/mgr"
	nodenormal "github.com/fananchong/go-xserver/internal/node/normal"
)

// Node : 节点组件
type Node struct {
	ctx   *common.Context
	node0 *nodemgr.Node
	node1 *nodenormal.Node
}

// NewNode : 实例化
func NewNode(ctx *common.Context) *Node {
	node := &Node{ctx: ctx}
	switch getPluginType(node.ctx) {
	case common.Mgr:
		node.node0 = nodemgr.NewNode(node.ctx)
		if node.node0.Init() {
			node.ctx.Node = node.node0
		}
	default:
		node.node1 = nodenormal.NewNode(node.ctx, pluginType)
		if node.node1.Init() {
			node.ctx.Node = node.node1
		}
	}
	return node
}

// Start : 实例化组件
func (node *Node) Start() bool {
	if node.node0 != nil {
		return node.node0.Start()
	}
	if node.node1 != nil {
		return node.node1.Start()
	}
	return false
}

// Close : 关闭组件
func (node *Node) Close() {
	if node.node0 != nil {
		node.node0.Close()
		node.node0 = nil
	}
	if node.node1 != nil {
		node.node1.Close()
		node.node1 = nil
	}
}
