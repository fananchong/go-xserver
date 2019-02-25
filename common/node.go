package common

import (
	"github.com/gogo/protobuf/proto"
	uuid "github.com/satori/go.uuid"
)

// NodeID : 节点ID类型
type NodeID uuid.UUID

// FuncTypeOnRelayMsg : 处理中继消息的函数声明
type FuncTypeOnRelayMsg func(source NodeType, sess INode, account string, cmd uint64, data []byte)

// INode : 节点类接口（其实现，封装自动接入服务器组、服务发现、服务消息传递等细节）
type INode interface {
	GetID() NodeID                                                                 // 【1】获取节点ID
	GetType() NodeType                                                             // 【1】获取节点类型
	GetIP(i IPType) string                                                         // 【1】获取IP
	GetPort(i int) int32                                                           // 【1】获取端口
	GetOverload(i int) uint32                                                      // 【1】获取负载
	GetVersion() string                                                            // 【1】获取版本号
	GetNodeOne(nodeType NodeType) INode                                            // 【2】根据节点类型，随机选择 1 节点
	GetNodeList(nodeType NodeType) []INode                                         // 【2】获取某类型节点列表
	GetNodeAll() []INode                                                           // 【2】获取所有节点
	HaveNode(nodeID NodeID) bool                                                   // 【2】是否有节点
	SendOne(nodeType NodeType, cmd uint64, msg proto.Message) bool                 // 【3】根据节点类型，随机选择 1 节点，发送数据
	SendByType(nodeType NodeType, cmd uint64, msg proto.Message, excludeSelf bool) // 【3】对某类型节点，广播数据
	SendByID(nodeID NodeID, cmd uint64, msg proto.Message) bool                    // 【3】往指定节点，发送数据
	SendMsg(cmd uint64, msg proto.Message) bool                                    // 【3】往该节点，发送数据
	SendAll(cmd uint64, msg proto.Message, excludeSelf bool)                       // 【3】对服务器组，广播数据
	EnableMessageRelay(v bool)                                                     // 【4】开启消息中继功能。开启该功能的节点，会连接 Gateway 。 C -> Gateway -> Node ; Node1 -> Gateway -> Node2(s)
	RegisterFuncOnRelayMsg(f FuncTypeOnRelayMsg)                                   // 【4】注册自定义处理Gateway中继过来的消息
	SendClientMsgByRelay(account string, cmd uint64, data []byte) bool             // 【4】发送消息给客户端，通过 Gateway 中继
}
