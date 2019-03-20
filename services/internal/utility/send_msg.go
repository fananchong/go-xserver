package utility

import (
	"fmt"

	"github.com/fananchong/go-xserver/common"
	"github.com/fananchong/gotcp"
	"github.com/gogo/protobuf/proto"
)

// SendMsgToClient : 发送数据
func SendMsgToClient(ctx *common.Context, account string, cmd uint64, msg proto.Message) (bool, error) {
	data, flag, err := gotcp.Encode(cmd, msg)
	if err != nil {
		return false, err
	}
	data = append(data, flag)
	if ctx.Node.SendMsgToClient(account, cmd, data) {
		return true, nil
	}
	return false, fmt.Errorf("Sending message failed, account: %s, cmd:%d", account, cmd)
}

// BroadcastMsgToClient :
func BroadcastMsgToClient(ctx *common.Context, cmd uint64, msg proto.Message) (bool, error) {
	data, flag, err := gotcp.Encode(cmd, msg)
	if err != nil {
		return false, err
	}
	data = append(data, flag)
	if ctx.Node.BroadcastMsgToClient(cmd, data) {
		return true, nil
	}
	return false, fmt.Errorf("Broadcast message failed, cmd:%d", cmd)
}
