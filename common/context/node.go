package context

import "github.com/fananchong/go-xserver/common/config"

// NodeID : 服务节点ID类型
type NodeID uint32

// FuncTypeOnRelayMsg : 处理中继消息的函数声明
type FuncTypeOnRelayMsg func(source config.NodeType, account string, cmd uint64, data []byte)

// FuncTypeOnLoseAccount : 处理丢失账号的函数声明
type FuncTypeOnLoseAccount func(account string)

// INode : 节点类接口（其实现，封装自动接入服务器组、服务发现、服务消息传递等细节）
type INode interface {
	EnableMessageRelay(v bool)                                            // 【4】开启消息中继功能。开启该功能的节点，会连接 Gateway 。 C -> Gateway -> Node ; Node1 -> Gateway -> Node2(s)
	RegisterFuncOnRelayMsg(f FuncTypeOnRelayMsg)                          // 【4】注册自定义处理 Gateway 中继过来的消息
	SendMsgToClient(account string, cmd uint64, data []byte) bool         // 【5】发送消息给客户端，通过 Gateway 中继
	BroadcastMsgToClient(cmd uint64, data []byte) bool                    // 【5】广播消息给客户端，通过 Gateway 中继
	SendMsgToServer(t config.NodeType, cmd uint64, data []byte) bool      // 【5】发送消息给某类型服务（随机一个）
	ReplyMsgToServer(cmd uint64, data []byte) bool                        // 【5】回发消息给请求服务器
	BroadcastMsgToServer(t config.NodeType, cmd uint64, data []byte) bool // 【5】广播消息给某类型服务
	RegisterFuncOnLoseAccount(f FuncTypeOnLoseAccount)                    // 【6】注册自定义处理`丢失账号`
}
