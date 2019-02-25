package utility

import (
	"fmt"

	"github.com/fananchong/go-xserver/common"
	"github.com/fananchong/gotcp"
	"github.com/gogo/protobuf/proto"
)

// SendMsgToClient : 发送数据
func SendMsgToClient(sess common.INode, account string, cmd uint64, msg proto.Message) (bool, error) {
	data, flag, err := gotcp.Encode(cmd, msg)
	if err != nil {
		return false, err
	}
	data = append(data, flag)
	if sess.SendClientMsgByRelay(account, cmd, data) {
		return true, nil
	}
	return false, fmt.Errorf("Sending message failed, account: %s, cmd:%d", account, cmd)
}
