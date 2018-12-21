package mgr

// NodeNormal : 普通节点
type NodeNormal struct {
	NodeBase
}

// NewNodeNormal : 普通节点实现类的构造函数
func NewNodeNormal() *NodeNormal {
	return &NodeNormal{}
}

// Init : 初始化节点
func (node *NodeNormal) Init() bool {
	return true
}
