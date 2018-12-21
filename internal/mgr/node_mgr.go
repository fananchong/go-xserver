package mgr

// NodeMgr : 管理节点
type NodeMgr struct {
	NodeBase
}

// NewNodeMgr : 管理节点实现类的构造函数
func NewNodeMgr() *NodeMgr {
	return &NodeMgr{}
}

// Init : 初始化节点
func (node *NodeMgr) Init() bool {
	return true
}
