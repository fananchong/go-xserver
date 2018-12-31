package nodemgr

import (
	"time"

	"github.com/fananchong/go-xserver/common"
	"github.com/fananchong/go-xserver/internal/db"
	nodecommon "github.com/fananchong/go-xserver/internal/node/common"
	"github.com/fananchong/go-xserver/internal/protocol"
	"github.com/fananchong/go-xserver/internal/utility"
	"github.com/fananchong/gotcp"
	"github.com/gogo/protobuf/proto"
)

// Node : 管理节点
type Node struct {
	nodecommon.DefaultNodeInterfaceImpl
	components []utility.IComponent
}

// NewNode : 管理节点实现类的构造函数
func NewNode() *Node {
	node := &Node{}
	node.Info = &protocol.SERVER_INFO{}
	node.Info.Id = utility.NodeID2ServerID(utility.NewNID())
	node.Info.Type = uint32(common.Mgr)
	node.Info.Addrs = []string{utility.GetIPInner(), utility.GetIPOuter()}
	node.Info.Ports = common.XCONFIG.Network.Port
	// TODO: 后续支持
	// node.Info.Overload
	// node.Info.Version
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
	nodecommon.XSESSIONMGR.ForAll(func(sess *nodecommon.SessionBase) {
		msg := &protocol.MSG_MGR_PING{}
		sess.SendMsg(uint64(protocol.CMD_MGR_PING), msg)
	})
}

// SendMsg : 往该节点，发送数据
func (node *Node) SendMsg(cmd uint64, msg proto.Message) bool {
	panic("")
}
