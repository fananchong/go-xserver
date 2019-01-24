package db

import (
	"encoding/json"

	"github.com/fananchong/go-xserver/common"
)

// AccountServer : 账号对应分配的服务资源
type AccountServer struct {
	NodeID  common.NodeID
	Address string
	Port    int32
	Type    common.NodeType
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
