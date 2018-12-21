package mgr

import "github.com/fananchong/go-xserver/common"

// NodeBase : 节点基类
type NodeBase struct {
}

// GetID : 获取本节点信息，节点ID
func (node *NodeBase) GetID() common.NodeID {
	return common.NodeID{}
}

// GetIP : 获取本节点信息，IP
func (node *NodeBase) GetIP(i common.IPType) string {
	return ""
}

// GetPort : 获取本节点信息，端口
func (node *NodeBase) GetPort(i int) int {
	return 0
}

// GetOverload : 获取本节点信息，负载
func (node *NodeBase) GetOverload(i int) uint {
	return 0
}

// GetVersion : 获取本节点信息，版本号
func (node *NodeBase) GetVersion() string {
	return common.XCONFIG.Common.Version
}

// SelectOne : 根据节点类型，随机选择 1 节点
func (node *NodeBase) SelectOne(nodeType common.NodeType) common.INode {
	return node
}

// SendOne : 根据节点类型，随机选择 1 节点，发送数据
func (node *NodeBase) SendOne(nodeType common.NodeType, data []byte) {

}

// SendByType : 对某类型节点，广播数据
func (node *NodeBase) SendByType(nodeType common.NodeType, data []byte, exclude []common.NodeID) {

}

// SendByID : 往指定节点，发送数据
func (node *NodeBase) SendByID(nodeID common.NodeID, data []byte) {

}

// Send : 往该节点，发送数据
func (node *NodeBase) Send(data []byte) {

}

// SendAll : 对服务器组，广播数据
func (node *NodeBase) SendAll(data []byte, exclude []common.NodeID) {

}
