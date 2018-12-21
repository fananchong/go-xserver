package internal

import (
	"github.com/fananchong/go-xserver/common"
	"github.com/fananchong/go-xserver/internal/mgr"
)

// Node : 节点实现类
type Node struct {
	nodeType common.NodeType
	*mgr.NodeBase
}

// NewNode : 节点实现类的构造函数
func NewNode(nodeType common.NodeType) *Node {
	return &Node{
		nodeType: nodeType,
	}
}

// Init : 初始化节点
func (node *Node) Init() bool {
	if node.nodeType == common.Mgr {
		impl := mgr.NewNodeMgr()
		node.NodeBase = &impl.NodeBase
		return impl.Init()
	} else {
		impl := mgr.NewNodeNormal()
		node.NodeBase = &impl.NodeBase
		return impl.Init()
	}
}
