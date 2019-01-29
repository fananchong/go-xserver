package gateway

import (
	"github.com/fananchong/gotcp"
)

// IntranetNode : 登录玩家类
type IntranetNode struct {
	gotcp.Session
}

// OnRecv : 接收到网络数据包，被触发
func (node *IntranetNode) OnRecv(data []byte, flag byte) {

}

// OnClose : 断开连接，被触发
func (node *IntranetNode) OnClose() {

}
