package internal

import (
	"github.com/fananchong/go-xserver/common"
	nodemgr "github.com/fananchong/go-xserver/internal/node_mgr"
	nodenormal "github.com/fananchong/go-xserver/internal/node_normal"
)

// NewNode : 节点实现类的构造函数
func NewNode(nodeType common.NodeType) common.INode {
	switch nodeType {
	case common.Mgr:
		node := nodemgr.NewNode()
		if node.Init() {
			return node
		}
	default:
		node := nodenormal.NewNode()
		if node.Init() {
			return node
		}
	}
	return nil
}
