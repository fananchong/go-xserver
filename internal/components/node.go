package components

import (
	"os"

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
	go func() {
		WaitComponent(node.ctx.Ctx)
		var err int
		if node.node0 != nil {
			err |= btoi(node.node0.Start())
		}
		if node.node1 != nil {
			err |= btoi(node.node1.Start())
		}
		if err != 0 {
			node.ctx.Log.Errorln("node start fail")
			os.Exit(1)
		}
	}()
	return true
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

func btoi(b bool) int {
	if b {
		return 0
	}
	return 1
}
