package nodecommon

import (
	"sync"

	"github.com/fananchong/go-xserver/common"
	"github.com/fananchong/go-xserver/internal/protocol"
	"github.com/fananchong/go-xserver/internal/utility"
	"github.com/gogo/protobuf/proto"
)

// DefaultNodeInterfaceImpl : 缺省的节点接口实现
type DefaultNodeInterfaceImpl struct {
	Info               *protocol.SERVER_INFO
	nid                common.NodeID
	once               sync.Once
	enableMessageReply bool
	SessMgr            *SessionMgr
	funcOnRelayMsg     common.FuncTypeOnRelayMsg
	funcOnLoseAccount  common.FuncTypeOnLoseAccount
}

// GetID : 获取本节点信息，节点ID
func (impl *DefaultNodeInterfaceImpl) GetID() common.NodeID {
	if impl.Info != nil {
		impl.once.Do(func() {
			impl.nid = utility.ServerID2NodeID(impl.Info.GetId())
		})
		return impl.nid
	}
	return common.NodeID{}
}

// GetType : 获取节点类型
func (impl *DefaultNodeInterfaceImpl) GetType() common.NodeType {
	if impl.Info != nil {
		return common.NodeType(impl.Info.GetType())
	}
	return common.Unknow
}

// GetIP : 获取本节点信息，IP
func (impl *DefaultNodeInterfaceImpl) GetIP(i common.IPType) string {
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
func (impl *DefaultNodeInterfaceImpl) GetNodeOne(nodeType common.NodeType) common.INode {
	node := impl.SessMgr.SelectOne(nodeType)
	if node != nil {
		return node
	}
	return nil
}

// GetNodeList : 获取某类型节点列表
func (impl *DefaultNodeInterfaceImpl) GetNodeList(nodeType common.NodeType) []common.INode {
	var ret []common.INode
	impl.SessMgr.ForByType(nodeType, func(sess *SessionBase) {
		ret = append(ret, sess)
	})
	return ret
}

// GetNodeAll : 获取所有节点
func (impl *DefaultNodeInterfaceImpl) GetNodeAll() []common.INode {
	var ret []common.INode
	impl.SessMgr.ForAll(func(sess *SessionBase) {
		ret = append(ret, sess)
	})
	return ret
}

// HaveNode : 是否有节点
func (impl *DefaultNodeInterfaceImpl) HaveNode(nodeID common.NodeID) bool {
	node := impl.SessMgr.GetByID(nodeID)
	return node != nil
}

// GetNode : 获取节点
func (impl *DefaultNodeInterfaceImpl) GetNode(nodeID common.NodeID) common.INode {
	node := impl.SessMgr.GetByID(nodeID)
	if node != nil {
		return node
	}
	return nil
}

// PrintNodeInfo : 打印节点信息
func (impl *DefaultNodeInterfaceImpl) PrintNodeInfo(log common.ILogger, nodeType common.NodeType) {
	log.Infoln("==========================================================================================================")
	for _, v := range impl.GetNodeList(nodeType) {
		log.Infoln("id:", utility.NodeID2UUID(v.GetID()).String(), "type:",
			v.GetType(), "port0:", v.GetPort(0), "port1:", v.GetPort(1), "ip0:",
			v.GetIP(common.IPINNER), "ip1:", v.GetIP(common.IPOUTER))
	}
	log.Infoln("----------------------------------------------------------------------------------------------------------")
}

// PrintAllNodeInfo : 打印节点信息
func (impl *DefaultNodeInterfaceImpl) PrintAllNodeInfo(log common.ILogger) {
	log.Infoln("==========================================================================================================")
	for _, v := range impl.GetNodeAll() {
		log.Infoln("id:", utility.NodeID2UUID(v.GetID()).String(), "type:",
			v.GetType(), "port0:", v.GetPort(0), "port1:", v.GetPort(1), "ip0:",
			v.GetIP(common.IPINNER), "ip1:", v.GetIP(common.IPOUTER))
	}
	log.Infoln("----------------------------------------------------------------------------------------------------------")
}

// SendOne : 根据节点类型，随机选择 1 节点，发送数据
func (impl *DefaultNodeInterfaceImpl) SendOne(nodeType common.NodeType, cmd uint64, msg proto.Message) bool {
	if sess := impl.SessMgr.SelectOne(nodeType); sess != nil {
		return sess.SendMsg(cmd, msg)
	}
	return false
}

// SendByType : 对某类型节点，广播数据
func (impl *DefaultNodeInterfaceImpl) SendByType(nodeType common.NodeType, cmd uint64, msg proto.Message, excludeSelf bool) {
	impl.SessMgr.ForByType(nodeType, func(sess *SessionBase) {
		if excludeSelf && utility.EqualNID(sess.GetID(), impl.GetID()) {
			return
		}
		sess.SendMsg(cmd, msg)
	})
}

// SendByID : 往指定节点，发送数据
func (impl *DefaultNodeInterfaceImpl) SendByID(nodeID common.NodeID, cmd uint64, msg proto.Message) bool {
	if sess := impl.SessMgr.GetByID(nodeID); sess != nil {
		return sess.SendMsg(cmd, msg)
	}
	return false
}

// SendAll : 对服务器组，广播数据
func (impl *DefaultNodeInterfaceImpl) SendAll(cmd uint64, msg proto.Message, excludeSelf bool) {
	impl.SessMgr.ForAll(func(sess *SessionBase) {
		if excludeSelf && utility.EqualNID(sess.GetID(), impl.GetID()) {
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
func (impl *DefaultNodeInterfaceImpl) RegisterFuncOnRelayMsg(f common.FuncTypeOnRelayMsg) {
	impl.funcOnRelayMsg = f
}

// FuncOnRelayMsg : 获取`自定义处理Gateway中继过来的消息`的函数句柄
func (impl *DefaultNodeInterfaceImpl) FuncOnRelayMsg() common.FuncTypeOnRelayMsg {
	return impl.funcOnRelayMsg
}

// RegisterFuncOnLoseAccount : 注册自定义处理`丢失账号`
func (impl *DefaultNodeInterfaceImpl) RegisterFuncOnLoseAccount(f common.FuncTypeOnLoseAccount) {
	impl.funcOnLoseAccount = f
}

// FuncOnLoseAccount : 获取`自定义处理丢失账号`的函数句柄
func (impl *DefaultNodeInterfaceImpl) FuncOnLoseAccount() common.FuncTypeOnLoseAccount {
	return impl.funcOnLoseAccount
}
