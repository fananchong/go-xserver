package db

import (
	"encoding/json"

	"github.com/fananchong/go-xserver/common"
)

// AccountServers : 账号对应分配的服务资源
type AccountServers struct {
	UIDList     []common.NodeID
	AddressList []string
	PortList    []int32
	TypeList    []common.NodeType
}

// Marshal : 序列化
func (accountservers *AccountServers) Marshal() (string, error) {
	data, err := json.Marshal(accountservers)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Unmarshal : 反序列化
func (accountservers *AccountServers) Unmarshal(data string) error {
	accountservers.reset()
	return json.Unmarshal([]byte(data), accountservers)
}

func (accountservers *AccountServers) reset() {
	accountservers.UIDList = accountservers.UIDList[:0]
	accountservers.AddressList = accountservers.AddressList[:0]
	accountservers.PortList = accountservers.PortList[:0]
	accountservers.TypeList = accountservers.TypeList[:0]
}
