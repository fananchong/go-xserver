package nodemgr

import (
	"time"

	"github.com/fananchong/go-xserver/common"
	"github.com/fananchong/go-xserver/common/utils"
	nodecommon "github.com/fananchong/go-xserver/internal/components/node/common"
	"github.com/fananchong/go-xserver/internal/db"
)

// Node : 管理节点
type Node struct {
	*nodecommon.Node
}

// NewNode : 管理节点实现类的构造函数
func NewNode(ctx *common.Context) *Node {
	node := &Node{}
	node.Node = nodecommon.NewNode(ctx, common.Mgr)
	return node
}

// Init : 初始化节点
func (node *Node) Init() bool {
	// register ticker
	registerTicker := utils.NewTickerHelper("REGISTER", node.Ctx, 1*time.Second, node.register)
	return node.Node.Init(Session{}, []utils.IComponent{registerTicker})
}

func (node *Node) register() {
	data := db.NewMgrServer(node.Ctx.Config.DbMgr.Name, 0)
	data.SetAddr(utils.GetIPInner(node.Ctx))
	data.SetPort(utils.GetIntranetListenPort(node.Ctx))
	if err := data.Save(); err != nil {
		node.Ctx.Log.Errorln(err)
	}
}
