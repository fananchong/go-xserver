package nodenormal

import (
	"github.com/fananchong/go-xserver/common"
	nodecommon "github.com/fananchong/go-xserver/internal/node/common"
	"github.com/fananchong/go-xserver/internal/protocol"
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
