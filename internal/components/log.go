package components

import (
	"os"

	"github.com/fananchong/glog"
	"github.com/fananchong/go-xserver/common"
	"github.com/fananchong/gotcp"
)

// Log : 日志组件
type Log struct {
	ctx *common.Context
}

// NewLog : 实例化
func NewLog(ctx *common.Context) *Log {
	log := &Log{ctx: ctx}
	log.init()
	return log
}

// Start : 实例化组件
func (log *Log) Start() bool {
	OneComponentOK(log.ctx.Ctx)
	return true
}

func (log *Log) init() {
	log.ctx.Log = glog.GetLogger()
	logDir := log.ctx.Config.Common.LogDir
	if logDir != "" {
		os.MkdirAll(logDir, os.ModePerm)
	}
	log.ctx.Log.SetLogDir(logDir)
	log.ctx.Log.SetLogLevel(log.ctx.Config.Common.LogLevel - 1)

	// TODO : gotcp 需要支持非全局LOG类实例
	gotcp.SetLogger(log.ctx.Log)
}

// Close : 关闭组件
func (log *Log) Close() {
	if log.ctx != nil && log.ctx.Log != nil {
		log.ctx.Log.Flush()
	}
}
