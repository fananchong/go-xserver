package nodemgr

import (
	"time"

	"github.com/fananchong/go-xserver/common"
	"github.com/fananchong/go-xserver/internal/db"
	"github.com/fananchong/go-xserver/internal/protocol"
	"github.com/fananchong/go-xserver/internal/utility"
	"github.com/fananchong/gotcp"
	"github.com/gogo/protobuf/proto"
)

// Node : 管理节点
type Node struct {
	defaultNodeInterfaceImpl
	components []utility.IComponent
}

// NewNode : 管理节点实现类的构造函数
func NewNode() *Node {
	node := &Node{}
	node.SetID(utility.NewNID())
	return node
}

// Init : 初始化节点
func (node *Node) Init() bool {
	// tcp server
	server := &gotcp.Server{}
	server.RegisterSessType(Session{})
	server.SetAddress(utility.GetIPInner(), utility.GetIntranetListenPort())
	server.SetUnfixedPort(true)

	// register ticker
	registerTicker := utility.NewTickerHelper(1*time.Second, node.register)

	// ping ticker
	pingTicker := utility.NewTickerHelper(5*time.Second, node.ping)

	// bind components
	node.components = []utility.IComponent{
		server,
		registerTicker,
		pingTicker,
	}
	return true
}

// Start : 节点开始工作
func (node *Node) Start() bool {
	for _, v := range node.components {
		if v != nil && v.Start() == false {
			panic("")
		}
	}
	return true
}

// Close : 关闭节点
func (node *Node) Close() {
	for _, v := range node.components {
		v.Close()
	}
}

func (node *Node) register() {
	data := db.NewMgrServer(common.XCONFIG.DbMgr.Name, 0)
	data.SetAddr(utility.GetIPInner())
	data.SetPort(utility.GetIntranetListenPort())
	if err := data.Save(); err != nil {
		common.XLOG.Errorln(err)
	}
}

func (node *Node) ping() {
	xsessionmgr.forAll(func(sess *Session) {
		msg := &protocol.MSG_MGR_PING{}
		sess.SendMsg(uint64(protocol.CMD_MGR_PING), msg)
	})
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
