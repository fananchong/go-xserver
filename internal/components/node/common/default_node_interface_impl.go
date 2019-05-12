package nodecommon

import (
	"sync"

	"github.com/fananchong/go-xserver/common/config"
	"github.com/fananchong/go-xserver/common/context"
	"github.com/fananchong/go-xserver/internal/protocol"
	"github.com/fananchong/go-xserver/internal/utils"
	"github.com/gogo/protobuf/proto"
)

// type INode interface {
// 	GetID() context.NodeID                                                         // 【1】获取节点ID
// 	GetType() NodeType                                                             // 【1】获取节点类型
// 	GetIP(i IPType) string                                                         // 【1】获取IP
// 	GetPort(i int) int32                                                           // 【1】获取端口
// 	GetOverload(i int) uint32                                                      // 【1】获取负载
// 	GetVersion() string                                                            // 【1】获取版本号
// 	GetNodeOne(nodeType NodeType) INode                                            // 【2】根据节点类型，随机选择 1 节点
// 	GetNodeList(nodeType NodeType) []INode                                         // 【2】获取某类型节点列表
// 	GetNodeAll() []INode                                                           // 【2】获取所有节点
// 	GetNode(nodeID context.NodeID) INode                                           // 【2】获取节点
// 	HaveNode(nodeID context.NodeID) bool                                           // 【2】是否有节点
// 	PrintNodeInfo(log ILogger, nodeType NodeType)                                  // 【2】打印节点信息
// 	PrintAllNodeInfo(log ILogger)                                                  // 【2】打印节点信息
// 	SendOne(nodeType NodeType, cmd uint64, msg proto.Message) bool                 // 【3】根据节点类型，随机选择 1 节点，发送数据
// 	SendByType(nodeType NodeType, cmd uint64, msg proto.Message, excludeSelf bool) // 【3】对某类型节点，广播数据
// 	SendByID(nodeID context.NodeID, cmd uint64, msg proto.Message) bool            // 【3】往指定节点，发送数据
// 	SendMsg(cmd uint64, msg proto.Message) bool                                    // 【3】往该节点，发送数据
// 	SendAll(cmd uint64, msg proto.Message, excludeSelf bool)                       // 【3】对服务器组，广播数据
//  EnableMessageRelay(v bool)                                                     // 【4】开启消息中继功能。开启该功能的节点，会连接 Gateway 。 C -> Gateway -> Node ; Node1 -> Gateway -> Node2(s)
//  RegisterFuncOnRelayMsg(f FuncTypeOnRelayMsg)                                   // 【4】注册自定义处理 Gateway 中继过来的消息
//  SendMsgToClient(account string, cmd uint64, data []byte) bool                  // 【5】发送消息给客户端，通过 Gateway 中继
//  BroadcastMsgToClient(cmd uint64, data []byte) bool                             // 【5】广播消息给客户端，通过 Gateway 中继
//	SendMsgToServer(t config.NodeType, cmd uint64, data []byte) bool               // 【5】发送消息给某类型服务（随机一个）
//  ReplyMsgToServer(cmd uint64, data []byte) bool                                 // 【5】回发消息给请求服务器
//	BroadcastMsgToServer(t config.NodeType, cmd uint64, data []byte) bool          // 【5】广播消息给某类型服务
//  RegisterFuncOnLoseAccount(f FuncTypeOnLoseAccount)                             // 【6】注册自定义处理`丢失账号`
// }

// DefaultNodeInterfaceImpl : 缺省的节点接口实现
type DefaultNodeInterfaceImpl struct {
	Info               *protocol.SERVER_INFO
	nid                context.NodeID
	once               sync.Once
	enableMessageReply bool
	SessMgr            *SessionMgr
	funcOnRelayMsg     context.FuncTypeOnRelayMsg
	funcOnLoseAccount  context.FuncTypeOnLoseAccount
}

// GetID : 获取本节点信息，节点ID
func (impl *DefaultNodeInterfaceImpl) GetID() context.NodeID {
	if impl.Info != nil {
		impl.once.Do(func() {
			impl.nid = ServerID2NodeID(impl.Info.GetId())
		})
		return impl.nid
	}
	return context.NodeID(0)
}

// GetType : 获取节点类型
func (impl *DefaultNodeInterfaceImpl) GetType() config.NodeType {
	if impl.Info != nil {
		return config.NodeType(impl.Info.GetType())
	}
	return config.Unknow
}

// GetIP : 获取本节点信息，IP
func (impl *DefaultNodeInterfaceImpl) GetIP(i utils.IPType) string {
	if impl.Info != nil {
		return impl.Info.GetAddrs()[i]
	}
	return ""
}

// GetPort : 获取本节点信息，端口
func (impl *DefaultNodeInterfaceImpl) GetPort(i int) int32 {
	if impl.Info != nil {
		return impl.Info.GetPorts()[i]
	}
	return 0
}

// GetOverload : 获取本节点信息，负载
func (impl *DefaultNodeInterfaceImpl) GetOverload(i int) uint32 {
	if impl.Info != nil {
		return impl.Info.GetOverload()[i]
	}
	return 0
}

// GetVersion : 获取本节点信息，版本号
func (impl *DefaultNodeInterfaceImpl) GetVersion() string {
	if impl.Info != nil {
		return impl.Info.GetVersion()
	}
	return ""
}

// GetSID : 获取 SID
func (impl *DefaultNodeInterfaceImpl) GetSID() *protocol.SERVER_ID {
	if impl.Info != nil {
		return impl.Info.GetId()
	}
	return &protocol.SERVER_ID{}
}

// GetNodeOne : 根据节点类型，随机选择 1 节点
func (impl *DefaultNodeInterfaceImpl) GetNodeOne(nodeType config.NodeType) *SessionBase {
	node := impl.SessMgr.SelectOne(nodeType)
	if node != nil {
		return node
	}
	return nil
}

// GetNodeList : 获取某类型节点列表
func (impl *DefaultNodeInterfaceImpl) GetNodeList(nodeType config.NodeType) []*SessionBase {
	var ret []*SessionBase
	impl.SessMgr.ForByType(nodeType, func(sess *SessionBase) {
		ret = append(ret, sess)
	})
	return ret
}

// GetNodeAll : 获取所有节点
func (impl *DefaultNodeInterfaceImpl) GetNodeAll() []*SessionBase {
	var ret []*SessionBase
	impl.SessMgr.ForAll(func(sess *SessionBase) {
		ret = append(ret, sess)
	})
	return ret
}

// HaveNode : 是否有节点
func (impl *DefaultNodeInterfaceImpl) HaveNode(nodeID context.NodeID) bool {
	node := impl.SessMgr.GetByID(nodeID)
	return node != nil
}

// GetNode : 获取节点
func (impl *DefaultNodeInterfaceImpl) GetNode(nodeID context.NodeID) *SessionBase {
	node := impl.SessMgr.GetByID(nodeID)
	if node != nil {
		return node
	}
	return nil
}

// PrintNodeInfo : 打印节点信息
func (impl *DefaultNodeInterfaceImpl) PrintNodeInfo(log context.ILogger, nodeType config.NodeType) {
	log.Infoln("==========================================================================================================")
	for _, v := range impl.GetNodeList(nodeType) {
		log.Infoln("id:", v.GetID(), "type:",
			v.GetType(), "port0:", v.GetPort(0), "port1:", v.GetPort(1), "ip0:",
			v.GetIP(utils.IPINNER), "ip1:", v.GetIP(utils.IPOUTER))
	}
	log.Infoln("----------------------------------------------------------------------------------------------------------")
}

// PrintAllNodeInfo : 打印节点信息
func (impl *DefaultNodeInterfaceImpl) PrintAllNodeInfo(log context.ILogger) {
	log.Infoln("==========================================================================================================")
	for _, v := range impl.GetNodeAll() {
		log.Infoln("id:", v.GetID(), "type:",
			v.GetType(), "port0:", v.GetPort(0), "port1:", v.GetPort(1), "ip0:",
			v.GetIP(utils.IPINNER), "ip1:", v.GetIP(utils.IPOUTER))
	}
	log.Infoln("----------------------------------------------------------------------------------------------------------")
}

// SendOne : 根据节点类型，随机选择 1 节点，发送数据
func (impl *DefaultNodeInterfaceImpl) SendOne(nodeType config.NodeType, cmd uint64, msg proto.Message) bool {
	if sess := impl.SessMgr.SelectOne(nodeType); sess != nil {
		return sess.SendMsg(cmd, msg)
	}
	return false
}

// SendByType : 对某类型节点，广播数据
func (impl *DefaultNodeInterfaceImpl) SendByType(nodeType config.NodeType, cmd uint64, msg proto.Message, excludeSelf bool) {
	impl.SessMgr.ForByType(nodeType, func(sess *SessionBase) {
		if excludeSelf && EqualNID(sess.GetID(), impl.GetID()) {
			return
		}
		sess.SendMsg(cmd, msg)
	})
}

// SendByID : 往指定节点，发送数据
func (impl *DefaultNodeInterfaceImpl) SendByID(nodeID context.NodeID, cmd uint64, msg proto.Message) bool {
	if sess := impl.SessMgr.GetByID(nodeID); sess != nil {
		return sess.SendMsg(cmd, msg)
	}
	return false
}

// SendAll : 对服务器组，广播数据
func (impl *DefaultNodeInterfaceImpl) SendAll(cmd uint64, msg proto.Message, excludeSelf bool) {
	impl.SessMgr.ForAll(func(sess *SessionBase) {
		if excludeSelf && EqualNID(sess.GetID(), impl.GetID()) {
			return
		}
		sess.SendMsg(cmd, msg)
	})
}

// EnableMessageRelay : 开启消息中继功能。开启该功能的节点，会连接 Gateway 。 C -> Gateway -> Node ; Node1 -> Gateway -> Node2(s)
func (impl *DefaultNodeInterfaceImpl) EnableMessageRelay(v bool) {
	impl.enableMessageReply = v
}

// IsEnableMessageRelay : 是否开启了消息中继功能
func (impl *DefaultNodeInterfaceImpl) IsEnableMessageRelay() bool {
	return impl.enableMessageReply
}

// RegisterFuncOnRelayMsg : 注册自定义处理Gateway中继过来的消息
func (impl *DefaultNodeInterfaceImpl) RegisterFuncOnRelayMsg(f context.FuncTypeOnRelayMsg) {
	impl.funcOnRelayMsg = f
}

// FuncOnRelayMsg : 获取`自定义处理Gateway中继过来的消息`的函数句柄
func (impl *DefaultNodeInterfaceImpl) FuncOnRelayMsg() context.FuncTypeOnRelayMsg {
	return impl.funcOnRelayMsg
}

// RegisterFuncOnLoseAccount : 注册自定义处理`丢失账号`
func (impl *DefaultNodeInterfaceImpl) RegisterFuncOnLoseAccount(f context.FuncTypeOnLoseAccount) {
	impl.funcOnLoseAccount = f
}

// FuncOnLoseAccount : 获取`自定义处理丢失账号`的函数句柄
func (impl *DefaultNodeInterfaceImpl) FuncOnLoseAccount() context.FuncTypeOnLoseAccount {
	return impl.funcOnLoseAccount
}
