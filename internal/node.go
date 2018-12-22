package internal

import (
	"github.com/fananchong/go-xserver/common"
	nodemgr "github.com/fananchong/go-xserver/internal/node_mgr"
	nodenormal "github.com/fananchong/go-xserver/internal/node_normal"
)

func newNode(nodeType common.NodeType) common.INode {
	switch nodeType {
	case common.Mgr:
		node := nodemgr.NewNode()
		if node.Init() {
			return node
		}
	default:
		node := nodenormal.NewNode(nodeType)
		if node.Init() {
			return node
		}
	}
	return nil
}

func startNode(node common.INode) bool {
	nodeType := node.GetType()
	switch nodeType {
	case common.Mgr:
		return node.(*nodemgr.Node).Start()
	default:
		return node.(*nodenormal.Node).Start()
	}
}
