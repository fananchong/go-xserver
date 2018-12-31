package nodenormal

import (
	"github.com/fananchong/go-xserver/common"
)

// Node : 普通节点
type Node struct {
	nodeType common.NodeType
	*Session
}

// NewNode : 普通节点实现类的构造函数
func NewNode(nodeType common.NodeType) *Node {
	return &Node{
		nodeType: nodeType,
		Session:  NewSession(),
	}
}

// Init : 初始化节点
func (node *Node) Init() bool {
	return true
}

// Start : 节点开始工作
func (node *Node) Start() bool {
	node.Session.Start()
	return true
}

// Close : 关闭节点
func (node *Node) Close() {
	node.Session.Close()
}

// GetType : 获取本节点信息，节点类型
func (node *Node) GetType() common.NodeType {
	return node.nodeType
}
