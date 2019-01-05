package components

import (
	"os"

	"github.com/fananchong/glog"
	"github.com/fananchong/go-xserver/common"
	"github.com/fananchong/gotcp"
)

// Log : 日志组件
type Log struct {
}

// Start : 实例化组件
func (*Log) Start() bool {
	common.XLOG = glog.GetLogger()
	logDir := common.XCONFIG.Common.LogDir
	if logDir != "" {
		os.MkdirAll(logDir, os.ModePerm)
	}
	common.XLOG.SetLogDir(logDir)
	common.XLOG.SetLogLevel(common.XCONFIG.Common.LogLevel - 1)
	gotcp.SetLogger(common.XLOG)
	return true
}

// Close : 关闭组件
func (*Log) Close() {
	common.XLOG.Flush()
}
