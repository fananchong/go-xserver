package nodenormal

import (
	"fmt"
	"time"

	"github.com/fananchong/go-xserver/common"
	nodecommon "github.com/fananchong/go-xserver/internal/components/node/common"
	"github.com/fananchong/go-xserver/internal/protocol"
	"github.com/fananchong/go-xserver/internal/utility"
)

// IntranetSession : 网络会话类
type IntranetSession struct {
	*nodecommon.SessionBase
}

// NewIntranetSession : 网络会话类的构造函数
func NewIntranetSession(ctx *common.Context) *IntranetSession {
	sess := &IntranetSession{}
	sess.SessionBase = nodecommon.NewSessionBase(ctx, sess)
	return sess
}

// Start : 启动
func (sess *IntranetSession) Start() {
	go func() {
		for {
			node := nodecommon.GetSessionMgr().GetByID(utility.ServerID2NodeID(sess.Info.GetId()))
			if node == nil {
				// 目标节点已丢失，不用试图去连接啦
				break
			}
			address := fmt.Sprintf("%s:%d", sess.Info.GetAddrs()[common.IPINNER], sess.Info.GetPorts()[common.PORTFORINTRANET])
			sess.Ctx.Log.Infoln("try connect gateway, address:", address, "node:", utility.ServerID2UUID(sess.Info.GetId()).String())
			if sess.Connect(address, sess) == false {
				time.Sleep(1 * time.Second)
				continue
			}
			sess.Verify()
			sess.Ctx.Log.Infoln("connect gateway success, address:", address, "node:", utility.ServerID2UUID(sess.Info.GetId()).String())
			break
		}
	}()
}

// DoRegister : 某节点注册时处理
func (sess *IntranetSession) DoRegister(msg *protocol.MSG_MGR_REGISTER_SERVER, data []byte, flag byte) {
}

// DoVerify : 验证时保存自己的注册消息
func (sess *IntranetSession) DoVerify(msg *protocol.MSG_MGR_REGISTER_SERVER, data []byte, flag byte) {
}

// DoLose : 节点丢失时处理
func (sess *IntranetSession) DoLose(msg *protocol.MSG_MGR_LOSE_SERVER, data []byte, flag byte) {

}

// DoClose : 节点关闭时处理
func (sess *IntranetSession) DoClose(sessbase *nodecommon.SessionBase) {
}
