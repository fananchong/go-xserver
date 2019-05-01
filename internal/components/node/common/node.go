package nodecommon

import (
	"sync"
	"time"

	"github.com/fananchong/go-xserver/common"
	"github.com/fananchong/go-xserver/config"
	"github.com/fananchong/go-xserver/internal/protocol"
	"github.com/fananchong/go-xserver/internal/utils"
	"github.com/fananchong/gotcp"
)

// Node : 管理节点基类
type Node struct {
	Ctx *common.Context
	DefaultNodeInterfaceImpl
	components []utils.IComponent
	mtx        sync.Mutex
}

// NewNode : 管理节点基类实现类的构造函数
func NewNode(ctx *common.Context, nodeType config.NodeType) *Node {
	node := &Node{Ctx: ctx}
	node.SessMgr = NewSessionMgr(ctx)
	node.Info = &protocol.SERVER_INFO{}
	node.Info.Id = NodeID2ServerID(NewNID())
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
			node.Ctx.Config.Network.Port[utils.PORTFORINTRANET] = v.(*gotcp.Server).GetRealPort()
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

// SendMsgToClient : 发送消息给客户端，通过 Gateway 中继
func (node *Node) SendMsgToClient(account string, cmd uint64, data []byte) bool {
	// Gateway 、 MgrServer 调用该接口会 panic
	//    - Gateway 不需要这个接口，没有意义
	//    - MgrServer 如果有需求需要通过 Gateway 给客户端消息，则可以实现之。优先级太低了！
	panic("")
}

// BroadcastMsgToClient : 广播消息给客户端，通过 Gateway 中继
func (node *Node) BroadcastMsgToClient(cmd uint64, data []byte) bool {
	// Gateway 、 MgrServer 调用该接口会 panic
	//    - Gateway 不需要这个接口，没有意义
	//    - MgrServer 如果有需求需要通过 Gateway 给客户端消息，则可以实现之。优先级太低了！
	panic("")
}
