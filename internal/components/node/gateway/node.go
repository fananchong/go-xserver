package nodegateway

import (
	"github.com/fananchong/go-xserver/common"
	"github.com/fananchong/go-xserver/common/utils"
	nodecommon "github.com/fananchong/go-xserver/internal/components/node/common"
)

// Node : 网关节点
type Node struct {
	*nodecommon.Node
}

// NewNode : 网关节点实现类的构造函数
func NewNode(ctx *common.Context) *Node {
	node := &Node{}
	node.Node = nodecommon.NewNode(ctx, common.Gateway)
	return node
}

// Init : 初始化节点
func (node *Node) Init() bool {
	return node.Node.Init(Session{}, []utils.IComponent{})
}
