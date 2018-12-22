package internal

import (
	"github.com/fananchong/glog"
	"github.com/fananchong/go-xserver/common"
)

func newGLogger() common.ILogger {
	return glog.GetLogger()
}
