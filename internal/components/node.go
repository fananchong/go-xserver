package components

import (
	"github.com/fananchong/go-xserver/common"
	nodemgr "github.com/fananchong/go-xserver/internal/node/mgr"
	nodenormal "github.com/fananchong/go-xserver/internal/node/normal"
)

// Node : 节点组件
type Node struct {
	node0 *nodemgr.Node
	node1 *nodenormal.Node
}

// Start : 实例化组件
func (node *Node) Start() bool {
	loadPlugin()
	switch pluginType {
	case common.Mgr:
		node.node0 = nodemgr.NewNode()
		if node.node0.Init() {
			common.XNODE = node.node0
			return node.node0.Start()
		}
	default:
		node.node1 = nodenormal.NewNode(pluginType)
		if node.node1.Init() {
			common.XNODE = node.node1
			return node.node1.Start()
		}
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
