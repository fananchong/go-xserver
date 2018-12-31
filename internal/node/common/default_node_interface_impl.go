package nodecommon

import (
	"github.com/fananchong/go-xserver/common"
	"github.com/fananchong/go-xserver/internal/protocol"
	"github.com/fananchong/go-xserver/internal/utility"
	"github.com/gogo/protobuf/proto"
)

// DefaultNodeInterfaceImpl : 缺省的节点接口实现
type DefaultNodeInterfaceImpl struct {
	nid  common.NodeID
	Info *protocol.SERVER_INFO
}

// SetID : 设置节点ID
func (impl *DefaultNodeInterfaceImpl) SetID(nid common.NodeID) {
	impl.nid = nid
}

// GetID : 获取本节点信息，节点ID
func (impl *DefaultNodeInterfaceImpl) GetID() common.NodeID {
	return impl.nid
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
	return XSESSIONMGR.SelectOne(nodeType)
}

// GetNodeList : 获取某类型节点列表
func (impl *DefaultNodeInterfaceImpl) GetNodeList(nodeType common.NodeType) []common.INode {
	var ret []common.INode
	XSESSIONMGR.ForByType(nodeType, func(sess *SessionBase) {
		ret = append(ret, sess)
	})
	return ret
}

// GetNodeAll : 获取所有节点
func (impl *DefaultNodeInterfaceImpl) GetNodeAll() []common.INode {
	var ret []common.INode
	XSESSIONMGR.ForAll(func(sess *SessionBase) {
		ret = append(ret, sess)
	})
	return ret
}

// SendOne : 根据节点类型，随机选择 1 节点，发送数据
func (impl *DefaultNodeInterfaceImpl) SendOne(nodeType common.NodeType, cmd uint64, msg proto.Message) bool {
	if sess := XSESSIONMGR.SelectOne(nodeType); sess != nil {
		return sess.SendMsg(cmd, msg)
	}
	return false
}

// SendByType : 对某类型节点，广播数据
func (impl *DefaultNodeInterfaceImpl) SendByType(nodeType common.NodeType, cmd uint64, msg proto.Message, excludeSelf bool) {
	XSESSIONMGR.ForByType(nodeType, func(sess *SessionBase) {
		if excludeSelf && utility.EqualNID(sess.GetID(), impl.GetID()) {
			return
		}
		sess.SendMsg(cmd, msg)
	})
}

// SendByID : 往指定节点，发送数据
func (impl *DefaultNodeInterfaceImpl) SendByID(nodeID common.NodeID, cmd uint64, msg proto.Message) bool {
	if sess := XSESSIONMGR.GetByID(nodeID); sess != nil {
		return sess.SendMsg(cmd, msg)
	}
	return false
}

// SendAll : 对服务器组，广播数据
func (impl *DefaultNodeInterfaceImpl) SendAll(cmd uint64, msg proto.Message, excludeSelf bool) {
	XSESSIONMGR.ForAll(func(sess *SessionBase) {
		if excludeSelf && utility.EqualNID(sess.GetID(), impl.GetID()) {
			return
		}
		sess.SendMsg(cmd, msg)
	})
}
