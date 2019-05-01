package utils

import (
	"errors"
	"net"
	"os"
	"regexp"
	"sync"

	"github.com/fananchong/go-xserver/common"
)

var (
	ipinner     string
	ipouter     string
	onceipinner sync.Once
	onceipouter sync.Once
)

// IPType : IP 类型
type IPType int

const (
	// IPINNER : 类型 0 ，内网 IP
	IPINNER IPType = iota

	// IPOUTER : 类型 1 ，外网 IP
	IPOUTER
)

// PortType : PORT 类型
type PortType int

const (
	// PORTFORCLIENT : 端口类型（对客户端）
	PORTFORCLIENT PortType = iota

	// PORTFORINTRANET : 端口类型（对内网）
	PORTFORINTRANET
)

// GetIPInner : 获取内网 IP
func GetIPInner(ctx *common.Context) string {
	onceipinner.Do(func() {
		switch ctx.Config.Network.IPType {
		case 0:
			ip, err := networkCard2IP(ctx.Config.Network.IPInner)
			if err != nil {
				ctx.Errorln(err)
				os.Exit(1)
			}
			ipinner = ip
		default:
			ipinner = ctx.Config.Network.IPInner
		}
	})
	return ipinner
}

// GetIPOuter : 获取外网 IP
func GetIPOuter(ctx *common.Context) string {
	onceipouter.Do(func() {
		switch ctx.Config.Network.IPType {
		case 0:
			ip, err := networkCard2IP(ctx.Config.Network.IPOuter)
			if err != nil {
				ctx.Errorln(err)
				os.Exit(1)
			}
			ipouter = ip
		default:
			ipouter = ctx.Config.Network.IPOuter
		}
	})
	return ipouter
}

// GetIP : 根据类型获取IP
func GetIP(ctx *common.Context, t IPType) string {
	switch t {
	case IPINNER:
		return GetIPInner(ctx)
	case IPOUTER:
		return GetIPOuter(ctx)
	default:
		panic("unknow ip type.")
	}
}

func networkCard2IP(name string) (string, error) {
	nic, err := net.InterfaceByName(name)
	if err != nil {
		return "", err
	}
	addresses, err := nic.Addrs()
	if err != nil {
		return "", err
	}
	r, _ := regexp.Compile(`((25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))\.){3}(25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))`)
	for _, addr := range addresses {
		s := r.FindAllString(addr.String(), -1)
		if len(s) != 0 {
			return s[0], nil
		}
	}
	return "", errors.New("no find address. nic: " + name)
}

// GetIntranetListenPort : 获取服务器组内监听端口
func GetIntranetListenPort(ctx *common.Context) int32 {
	return ctx.Config.Network.Port[PORTFORINTRANET]
}

// GetDefaultServicePort : 获取缺省的对外端口
func GetDefaultServicePort(ctx *common.Context) int32 {
	return ctx.Config.Network.Port[PORTFORCLIENT]
}
