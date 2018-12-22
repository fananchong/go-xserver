package utility

import (
	"errors"
	"net"
	"os"
	"regexp"

	"github.com/fananchong/go-xserver/common"
)

// GetIPInner : 获取内网 IP
func GetIPInner() string {
	switch common.XCONFIG.Network.IPType {
	case 0:
		ip, err := networkCard2IP(common.XCONFIG.Network.IPInner)
		if err != nil {
			common.XLOG.Errorln(err)
			os.Exit(1)
		}
		return ip
	default:
		return common.XCONFIG.Network.IPInner
	}
}

// GetIPOuter : 获取外网 IP
func GetIPOuter() string {
	switch common.XCONFIG.Network.IPType {
	case 0:
		ip, err := networkCard2IP(common.XCONFIG.Network.IPOuter)
		if err != nil {
			common.XLOG.Errorln(err)
			os.Exit(1)
		}
		return ip
	default:
		return common.XCONFIG.Network.IPOuter
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
