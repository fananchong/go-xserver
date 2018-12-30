package common

import (
	"github.com/gogo/protobuf/proto"
	uuid "github.com/satori/go.uuid"
)

// XNODE : 本节点对象
var XNODE INode

// NodeID : 节点ID类型
type NodeID uuid.UUID

// INode : 节点类接口（其实现，封装自动接入服务器组、服务发现、服务消息传递等细节）
type INode interface {
	GetID() NodeID                                                                 // 获取节点ID
	GetType() NodeType                                                             // 获取节点类型
	GetIP(i IPType) string                                                         // 获取IP
	GetPort(i int) int32                                                           // 获取端口
	GetOverload(i int) uint32                                                      // 获取负载
	GetVersion() string                                                            // 获取版本号
	GetNodeOne(nodeType NodeType) INode                                            // 根据节点类型，随机选择 1 节点
	GetNodeList(nodeType NodeType) []INode                                         // 获取某类型节点列表
	GetNodeAll() []INode                                                           // 获取所有节点
	SendOne(nodeType NodeType, cmd uint64, msg proto.Message) bool                 // 根据节点类型，随机选择 1 节点，发送数据
	SendByType(nodeType NodeType, cmd uint64, msg proto.Message, excludeSelf bool) // 对某类型节点，广播数据
	SendByID(nodeID NodeID, cmd uint64, msg proto.Message) bool                    // 往指定节点，发送数据
	SendMsg(cmd uint64, msg proto.Message) bool                                    // 往该节点，发送数据
	SendAll(cmd uint64, msg proto.Message, excludeSelf bool)                       // 对服务器组，广播数据
}
