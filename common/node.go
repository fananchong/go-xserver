package common

import (
	uuid "github.com/satori/go.uuid"
)

// NodeID : 节点ID类型
type NodeID uuid.UUID

// FuncTypeOnRelayMsg : 处理中继消息的函数声明
type FuncTypeOnRelayMsg func(source NodeType, sess INode, account string, cmd uint64, data []byte)

// FuncTypeOnLoseAccount : 处理丢失账号的函数声明
type FuncTypeOnLoseAccount func(account string)

// INode : 节点类接口（其实现，封装自动接入服务器组、服务发现、服务消息传递等细节）
type INode interface {
	EnableMessageRelay(v bool)                                         // 【4】开启消息中继功能。开启该功能的节点，会连接 Gateway 。 C -> Gateway -> Node ; Node1 -> Gateway -> Node2(s)
	RegisterFuncOnRelayMsg(f FuncTypeOnRelayMsg)                       // 【4】注册自定义处理 Gateway 中继过来的消息
	SendClientMsgByRelay(account string, cmd uint64, data []byte) bool // 【4】发送消息给客户端，通过 Gateway 中继
	RegisterFuncOnLoseAccount(f FuncTypeOnLoseAccount)                 // 【5】注册自定义处理`丢失账号`
}
