package internal

import (
	"github.com/fananchong/glog"
	"github.com/fananchong/go-xserver/common"
)

// NewGLogger : Constructor function of class glog
func NewGLogger() common.ILogger {
	return glog.GetLogger()
}
