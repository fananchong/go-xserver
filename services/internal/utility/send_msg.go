package utility

import (
	"fmt"

	"github.com/fananchong/go-xserver/common"
	"github.com/fananchong/go-xserver/common/config"
	"github.com/fananchong/gotcp"
	"github.com/gogo/protobuf/proto"
)

// SendMsgToClient : 发送数据给客户端
func SendMsgToClient(ctx *common.Context, account string, cmd uint64, msg proto.Message) (bool, error) {
	data, flag, err := gotcp.Encode(cmd, msg)
	if err != nil {
		return false, err
	}
	data = append(data, flag)
	if ctx.SendMsgToClient(account, cmd, data) {
		return true, nil
	}
	return false, fmt.Errorf("Sending message failed, account: %s, cmd:%d", account, cmd)
}

// SendMsgToClientByRoleName : 发送数据给客户端
func SendMsgToClientByRoleName(ctx *common.Context, roleName string, cmd uint64, msg proto.Message) (bool, error) {
	account := ctx.IRole2Account.GetAndActive(roleName)
	if account != "" {
		return SendMsgToClient(ctx, account, cmd, msg)
	}
	return false, fmt.Errorf("Sending message failed, roleName: %s, cmd:%d, account:%s", roleName, cmd, account)
}

// BroadcastMsgToClient : 广播数据给客户端
func BroadcastMsgToClient(ctx *common.Context, cmd uint64, msg proto.Message) (bool, error) {
	data, flag, err := gotcp.Encode(cmd, msg)
	if err != nil {
		return false, err
	}
	data = append(data, flag)
	if ctx.BroadcastMsgToClient(cmd, data) {
		return true, nil
	}
	return false, fmt.Errorf("Broadcast message failed, cmd:%d", cmd)
}

// SendMsgToServer : 发送消息给某类型服务（随机一个）
func SendMsgToServer(ctx *common.Context, t config.NodeType, cmd uint64, msg proto.Message) (bool, error) {
	data, flag, err := gotcp.EncodeCmd(cmd, msg)
	if err != nil {
		return false, err
	}
	data = append(data, flag)
	if ctx.SendMsgToServer(t, cmd, data) {
		return true, nil
	}
	return false, fmt.Errorf("SendMsgToServer message failed, cmd:%d", cmd)
}

// ReplyMsgToServer : 回发消息给请求服务器
func ReplyMsgToServer(ctx *common.Context, cmd uint64, msg proto.Message) (bool, error) {
	data, flag, err := gotcp.EncodeCmd(cmd, msg)
	if err != nil {
		return false, err
	}
	data = append(data, flag)
	if ctx.ReplyMsgToServer(cmd, data) {
		return true, nil
	}
	return false, fmt.Errorf("ReplyMsgToServer message failed, cmd:%d", cmd)
}

// BroadcastMsgToServer : 广播消息给某类型服务
func BroadcastMsgToServer(ctx *common.Context, t config.NodeType, cmd uint64, msg proto.Message) (bool, error) {
	data, flag, err := gotcp.EncodeCmd(cmd, msg)
	if err != nil {
		return false, err
	}
	data = append(data, flag)
	if ctx.BroadcastMsgToServer(t, cmd, data) {
		return true, nil
	}
	return false, fmt.Errorf("BroadcastMsgToServer message failed, cmd:%d", cmd)
}
