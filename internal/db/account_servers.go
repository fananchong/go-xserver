package db

import (
	"encoding/json"
	"fmt"

	"github.com/fananchong/go-xserver/common"
	"github.com/fananchong/go-xserver/internal/protocol"
)

// 账号对应分配的服务资源
// Login Server 会负责分配
// 可以分配多个服务器资源

// AccountServer : 分配的服务资源信息（单个）
type AccountServer struct {
	ServerID *protocol.SERVER_ID
	Address  string
	Port     int32
	Type     common.NodeType
}

// Marshal : 序列化
func (accountserver *AccountServer) Marshal() (string, error) {
	data, err := json.Marshal(accountserver)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Unmarshal : 反序列化
func (accountserver *AccountServer) Unmarshal(data string) error {
	return json.Unmarshal([]byte(data), accountserver)
}

// GetKeyAllocServer : 账号对应的服务器资源的 KEY
func GetKeyAllocServer(nodeType uint32, account string) string {
	return fmt.Sprintf("srv%d_%s", nodeType, account)
}
