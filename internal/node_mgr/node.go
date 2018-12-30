package nodemgr

import (
	"fmt"

	"github.com/fananchong/go-xserver/common"
	"github.com/fananchong/go-xserver/internal/utility"
	"github.com/fananchong/gotcp"
	"github.com/gogo/protobuf/proto"
)

// Node : 管理节点
type Node struct {
	server         gotcp.Server
	registerHelper RegisterMgrHelper
	defaultNodeInterfaceImpl
}

// NewNode : 管理节点实现类的构造函数
func NewNode() *Node {
	node := &Node{}
	node.defaultNodeInterfaceImpl.nid = utility.NewNID()
	return node
}

// Init : 初始化节点
func (node *Node) Init() bool {
	node.server.RegisterSessType(Session{})
	return true
}

// Start : 节点开始工作
func (node *Node) Start() bool {
	node.registerHelper.Start()
	return node.server.Start(fmt.Sprintf("%s:%d", utility.GetIPInner(), utility.GetIntranetListenPort()))
}

// Close : 关闭节点
func (node *Node) Close() {
	node.registerHelper.Close()
	node.server.Close()
}

// GetID : 获取本节点信息，节点ID
func (node *Node) GetID() common.NodeID {
	return node.nid
}

// GetType : 获取本节点信息，节点类型
func (node *Node) GetType() common.NodeType {
	return common.Mgr
}

// GetIP : 获取本节点信息，IP
func (node *Node) GetIP(i common.IPType) string {
	return utility.GetIP(i)
}

// GetPort : 获取本节点信息，端口
func (node *Node) GetPort(i int) int32 {
	return common.XCONFIG.Network.Port[i]
}

// GetOverload : 获取本节点信息，负载
func (node *Node) GetOverload(i int) uint32 {
	// TODO:
	return 0
}

// GetVersion : 获取本节点信息，版本号
func (node *Node) GetVersion() string {
	return common.XCONFIG.Common.Version
}

// SendMsg : 往该节点，发送数据
func (node *Node) SendMsg(cmd uint64, msg proto.Message) bool {
	panic("")
}
