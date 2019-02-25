package nodecommon

import (
	"sync"
	"time"

	"github.com/fananchong/go-xserver/common"
	"github.com/fananchong/go-xserver/common/utils"
	"github.com/fananchong/go-xserver/internal/protocol"
	"github.com/fananchong/go-xserver/internal/utility"
	"github.com/fananchong/gotcp"
	"github.com/gogo/protobuf/proto"
)

// Node : 管理节点基类
type Node struct {
	Ctx *common.Context
	DefaultNodeInterfaceImpl
	components []utils.IComponent
	mtx        sync.Mutex
}

// NewNode : 管理节点基类实现类的构造函数
func NewNode(ctx *common.Context, nodeType common.NodeType) *Node {
	node := &Node{Ctx: ctx}
	node.SessMgr = NewSessionMgr()
	node.Info = &protocol.SERVER_INFO{}
	node.Info.Id = utility.NodeID2ServerID(utility.NewNID())
	node.Info.Type = uint32(nodeType)
	node.Info.Addrs = []string{utils.GetIPInner(ctx), utils.GetIPOuter(ctx)}
	node.Info.Ports = ctx.Config.Network.Port
	// TODO: 后续支持
	// node.Info.Overload
	// node.Info.Version
	return node
}

// Init : 初始化节点
func (node *Node) Init(sessType interface{}, components []utils.IComponent) bool {
	// tcp server
	server := &gotcp.Server{}
	server.RegisterSessType(sessType)
	server.SetAddress(utils.GetIPInner(node.Ctx), utils.GetIntranetListenPort(node.Ctx))
	server.SetUnfixedPort(true)
	server.SetUserData(&UserData{Ctx: node.Ctx, SessMgr: node.SessMgr})

	// ping ticker
	pingTicker := utils.NewTickerHelper("PING", node.Ctx, 5*time.Second, node.ping)

	// bind components
	node.components = []utils.IComponent{
		server,
		pingTicker,
	}
	node.components = append(node.components, components...)
	return true
}

// Start : 节点开始工作
func (node *Node) Start() bool {
	node.mtx.Lock()
	defer node.mtx.Unlock()
	for _, v := range node.components {
		if v != nil && v.Start() == false {
			panic("")
		}
		switch v.(type) {
		case *gotcp.Server:
			node.Ctx.Config.Network.Port[common.PORTFORINTRANET] = v.(*gotcp.Server).GetRealPort()
		}
	}
	return true
}

// Close : 关闭节点
func (node *Node) Close() {
	node.mtx.Lock()
	defer node.mtx.Unlock()
	for _, v := range node.components {
		v.Close()
	}
}

func (node *Node) ping() {
	node.SessMgr.ForAll(func(sess *SessionBase) {
		msg := &protocol.MSG_MGR_PING{}
		sess.SendMsg(uint64(protocol.CMD_MGR_PING), msg)
	})
}

// SendMsg : 往该节点，发送数据
func (node *Node) SendMsg(cmd uint64, msg proto.Message) bool {
	panic("")
}

// SendClientMsgByRelay : 发送消息给客户端，通过 Gateway 中继
func (node *Node) SendClientMsgByRelay(account string, cmd uint64, data []byte) bool {
	panic("")
}
