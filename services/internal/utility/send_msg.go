package utility

import (
	"fmt"

	"github.com/fananchong/go-xserver/common"
	"github.com/fananchong/go-xserver/common/config"
	"github.com/fananchong/go-xserver/common/context"
	"github.com/fananchong/gotcp"
	"github.com/gogo/protobuf/proto"
)

// 框架层使用 []byte 数据块作为参数
// 逻辑层可以自定义协议格式，这里用的是 protobuf
// 封装下接口，方便代码调用

// SendMsgToClient : 发送数据给客户端
func SendMsgToClient(ctx *common.Context, account string, cmd uint64, msg proto.Message) (bool, error) {
	data, flag, err := gotcp.Encode(cmd, msg)
	if err != nil {
		return false, err
	}
	if ctx.SendMsgToClient(account, cmd, data, flag) {
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
	if ctx.BroadcastMsgToClient(cmd, data, flag) {
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
	if ctx.SendMsgToServer(t, cmd, data, flag) {
		return true, nil
	}
	return false, fmt.Errorf("SendMsgToServer message failed, cmd:%d", cmd)
}

// ReplyMsgToServer : 回发消息给请求服务器
func ReplyMsgToServer(ctx *common.Context, targetID context.NodeID, cmd uint64, msg proto.Message) (bool, error) {
	data, flag, err := gotcp.EncodeCmd(cmd, msg)
	if err != nil {
		return false, err
	}
	if ctx.ReplyMsgToServer(targetID, cmd, data, flag) {
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
	if ctx.BroadcastMsgToServer(t, cmd, data, flag) {
		return true, nil
	}
	return false, fmt.Errorf("BroadcastMsgToServer message failed, cmd:%d", cmd)
}
