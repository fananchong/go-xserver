package db

import "encoding/json"

// AccountServers : 账号对应分配的服务资源
type AccountServers struct {
	AddressList []string
	PortList    []int32
	TypeList    []int32
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
	return json.Unmarshal([]byte(data), accountservers)
}
