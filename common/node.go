package common

import uuid "github.com/satori/go.uuid"

// XNODE : 本节点对象
var XNODE INode

// NodeID : 节点ID类型
type NodeID uuid.UUID

// INode : 节点类接口（其实现，封装自动接入服务器组、服务发现、服务消息传递等细节）
type INode interface {
	// 获取本节点信息
	GetID() NodeID          // 节点ID
	GetType() NodeType      // 节点类型
	GetIP(i IPType) string  // IP
	GetPort(i int) int      // 端口
	GetOverload(i int) uint // 负载
	GetVersion() string     // 版本号

	// 根据节点类型，随机选择 1 节点
	SelectOne(nodeType NodeType) INode

	// 发送消息
	SendOne(nodeType NodeType, data []byte)                      // 根据节点类型，随机选择 1 节点，发送数据
	SendByType(nodeType NodeType, data []byte, exclude []NodeID) // 对某类型节点，广播数据
	SendByID(nodeID NodeID, data []byte)                         // 往指定节点，发送数据
	Send(data []byte)                                            // 往该节点，发送数据
	SendAll(data []byte, exclude []NodeID)                       // 对服务器组，广播数据

	// 注册自定义服务器组内网络事件（TODO: 参数待定）
	RegisterOnConnect()
	RegisterOnRecv()
	RegisterOnDisconnect()
}
