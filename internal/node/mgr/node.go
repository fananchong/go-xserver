package nodemgr

import (
	"time"

	"github.com/fananchong/go-xserver/common"
	"github.com/fananchong/go-xserver/common/utils"
	"github.com/fananchong/go-xserver/internal/db"
	nodecommon "github.com/fananchong/go-xserver/internal/node/common"
	"github.com/fananchong/go-xserver/internal/protocol"
	"github.com/fananchong/go-xserver/internal/utility"
	"github.com/fananchong/gotcp"
	"github.com/gogo/protobuf/proto"
)

// Node : 管理节点
type Node struct {
	ctx *common.Context
	nodecommon.DefaultNodeInterfaceImpl
	components []utils.IComponent
}

// NewNode : 管理节点实现类的构造函数
func NewNode(ctx *common.Context) *Node {
	node := &Node{ctx: ctx}
	node.Info = &protocol.SERVER_INFO{}
	node.Info.Id = utility.NodeID2ServerID(utility.NewNID())
	node.Info.Type = uint32(common.Mgr)
	node.Info.Addrs = []string{utils.GetIPInner(ctx), utils.GetIPOuter(ctx)}
	node.Info.Ports = ctx.Config.Network.Port
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
	server.SetAddress(utils.GetIPInner(node.ctx), utils.GetIntranetListenPort(node.ctx))
	server.SetUnfixedPort(true)
	server.SetUserData(node.ctx)

	// register ticker
	registerTicker := utils.NewTickerHelper(node.ctx, 1*time.Second, node.register)

	// ping ticker
	pingTicker := utils.NewTickerHelper(node.ctx, 5*time.Second, node.ping)

	// bind components
	node.components = []utils.IComponent{
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
	data := db.NewMgrServer(node.ctx.Config.DbMgr.Name, 0)
	data.SetAddr(utils.GetIPInner(node.ctx))
	data.SetPort(utils.GetIntranetListenPort(node.ctx))
	if err := data.Save(); err != nil {
		node.ctx.Log.Errorln(err)
	}
}

func (node *Node) ping() {
	nodecommon.GetSessionMgr().ForAll(func(sess *nodecommon.SessionBase) {
		msg := &protocol.MSG_MGR_PING{}
		sess.SendMsg(uint64(protocol.CMD_MGR_PING), msg)
	})
}

// SendMsg : 往该节点，发送数据
func (node *Node) SendMsg(cmd uint64, msg proto.Message) bool {
	panic("")
}
