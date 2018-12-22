package internal

import (
	"github.com/fananchong/glog"
	"github.com/fananchong/go-xserver/common"
)

// NewGLogger : glog 构造函数
func NewGLogger() common.ILogger {
	return glog.GetLogger()
}
