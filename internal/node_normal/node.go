package nodenormal

import (
	"github.com/fananchong/go-xserver/common"
	"github.com/gogo/protobuf/proto"
)

// Node : 普通节点
type Node struct {
	nodeType common.NodeType
	sess     *Session
}

// NewNode : 普通节点实现类的构造函数
func NewNode(nodeType common.NodeType) *Node {
	return &Node{
		nodeType: nodeType,
		sess:     NewSession(),
	}
}

// Init : 初始化节点
func (node *Node) Init() bool {
	return true
}

// Start : 节点开始工作
func (node *Node) Start() bool {
	node.sess.Start()
	return true
}

// Close : 关闭节点
func (node *Node) Close() {
	node.sess.Close()
}

// GetID : 获取本节点信息，节点ID
func (node *Node) GetID() common.NodeID {
	return node.sess.NID
}

// GetType : 获取本节点信息，节点类型
func (node *Node) GetType() common.NodeType {
	return node.nodeType
}

// GetIP : 获取本节点信息，IP
func (node *Node) GetIP(i common.IPType) string {
	return ""
}

// GetPort : 获取本节点信息，端口
func (node *Node) GetPort(i int) int32 {
	return 0
}

// GetOverload : 获取本节点信息，负载
func (node *Node) GetOverload(i int) uint32 {
	return 0
}

// GetVersion : 获取本节点信息，版本号
func (node *Node) GetVersion() string {
	return common.XCONFIG.Common.Version
}

// SelectOne : 根据节点类型，随机选择 1 节点
func (node *Node) GetNodeOne(nodeType common.NodeType) common.INode {
	return node
}

// GetNodeList : 获取某类型节点列表
func (node *Node) GetNodeList(nodeType common.NodeType) []common.INode {
	return nil
}

// GetAllNode : 获取所有节点
func (node *Node) GetNodeAll() []common.INode {
	return nil
}

// SendOne : 根据节点类型，随机选择 1 节点，发送数据
func (node *Node) SendOne(nodeType common.NodeType, cmd uint64, msg proto.Message) bool {
	return false
}

// SendByType : 对某类型节点，广播数据
func (node *Node) SendByType(nodeType common.NodeType, cmd uint64, msg proto.Message, excludeSelf bool) {

}

// SendByID : 往指定节点，发送数据
func (node *Node) SendByID(nodeID common.NodeID, cmd uint64, msg proto.Message) bool {
	return false
}

// Send : 往该节点，发送数据
func (node *Node) SendMsg(cmd uint64, msg proto.Message) bool {
	return false
}

// SendAll : 对服务器组，广播数据
func (node *Node) SendAll(cmd uint64, msg proto.Message, excludeSelf bool) {

}

func (node *Node) RegisterOnConnect()    {}
func (node *Node) RegisterOnRecv()       {}
func (node *Node) RegisterOnDisconnect() {}
