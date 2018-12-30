package nodemgr

import (
	"github.com/fananchong/go-xserver/common"
	"github.com/fananchong/go-xserver/internal/utility"
	"github.com/gogo/protobuf/proto"
)

// defaultNodeInterfaceImpl : 缺省的节点接口实现
type defaultNodeInterfaceImpl struct {
	nid common.NodeID
}

func (impl *defaultNodeInterfaceImpl) SetID(nid common.NodeID) {
	impl.nid = nid
}

// GetID : 获取本节点信息，节点ID
func (impl *defaultNodeInterfaceImpl) GetID() common.NodeID {
	return impl.nid
}

// SelectOne : 根据节点类型，随机选择 1 节点
func (impl *defaultNodeInterfaceImpl) GetNodeOne(nodeType common.NodeType) common.INode {
	return xsessionmgr.selectOne(nodeType)
}

// GetNodeList : 获取某类型节点列表
func (impl *defaultNodeInterfaceImpl) GetNodeList(nodeType common.NodeType) []common.INode {
	var ret []common.INode
	xsessionmgr.forByType(nodeType, func(sess *Session) {
		ret = append(ret, sess)
	})
	return ret
}

// GetAllNode : 获取所有节点
func (impl *defaultNodeInterfaceImpl) GetNodeAll() []common.INode {
	var ret []common.INode
	xsessionmgr.forAll(func(sess *Session) {
		ret = append(ret, sess)
	})
	return ret
}

// SendOne : 根据节点类型，随机选择 1 节点，发送数据
func (impl *defaultNodeInterfaceImpl) SendOne(nodeType common.NodeType, cmd uint64, msg proto.Message) bool {
	if sess := xsessionmgr.selectOne(nodeType); sess != nil {
		return sess.SendMsg(cmd, msg)
	}
	return false
}

// SendByType : 对某类型节点，广播数据
func (impl *defaultNodeInterfaceImpl) SendByType(nodeType common.NodeType, cmd uint64, msg proto.Message, excludeSelf bool) {
	xsessionmgr.forByType(nodeType, func(sess *Session) {
		if excludeSelf && utility.EqualNID(sess.GetID(), impl.GetID()) {
			return
		}
		sess.SendMsg(cmd, msg)
	})
}

// SendByID : 往指定节点，发送数据
func (impl *defaultNodeInterfaceImpl) SendByID(nodeID common.NodeID, cmd uint64, msg proto.Message) bool {
	if sess := xsessionmgr.getByID(nodeID); sess != nil {
		return sess.SendMsg(cmd, msg)
	}
	return false
}

// SendAll : 对服务器组，广播数据
func (impl *defaultNodeInterfaceImpl) SendAll(cmd uint64, msg proto.Message, excludeSelf bool) {
	xsessionmgr.forAll(func(sess *Session) {
		if excludeSelf && utility.EqualNID(sess.GetID(), impl.GetID()) {
			return
		}
		sess.SendMsg(cmd, msg)
	})
}
