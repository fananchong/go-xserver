package components

import (
	"os"
	"path/filepath"
	"time"

	"github.com/fananchong/glog"
	"github.com/fananchong/go-xserver/common"
	"github.com/fananchong/gotcp"
	"github.com/spf13/viper"
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
	return true
}

func (log *Log) init() {
	tmpLog := glog.GetLogger()
	tmpLog.SetAppName(filepath.Base(os.Args[0]) + "_" + viper.GetString("app") + viper.GetString("suffix"))
	logDir := log.ctx.Config.Common.LogDir
	if logDir != "" {
		os.MkdirAll(logDir, os.ModePerm)
	}
	tmpLog.SetLogDir(logDir)
	tmpLog.SetLogLevel(log.ctx.Config.Common.LogLevel - 1)
	tmpLog.SetFlushInterval(time.Duration(log.ctx.Config.Common.LogFlushInterval) * time.Millisecond)
	log.ctx.Log = tmpLog

	// TODO : gotcp 需要支持非全局LOG类实例
	gotcp.SetLogger(log.ctx.Log)
}

// Close : 关闭组件
func (log *Log) Close() {
	if log.ctx != nil && log.ctx.Log != nil {
		log.ctx.Log.Flush()
	}
}
