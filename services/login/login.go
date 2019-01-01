package main

import (
	"fmt"
	"net/http"

	"github.com/fananchong/go-xserver/common"

	"github.com/fananchong/go-xserver/common/utils"
)

var xsvr = NewLogin()

// Login : 登陆服务器
type Login struct {
	web  *http.Server
	addr string
}

// NewLogin : 构造函数
func NewLogin() *Login {
	return &Login{}
}

// Init : 初始化
func (login *Login) Init() {
	login.addr = fmt.Sprintf("%s:%d", utils.GetIPOuter(), utils.GetDefaultServicePort())
	login.web = &http.Server{Addr: login.addr, Handler: nil}
}

// RegisterCallBack : 注册回调
func (login *Login) RegisterCallBack() {
	http.HandleFunc("/login", login.msgLogin)
}

// Start : 启动
func (login *Login) Start() bool {
	go func() {
		common.XLOG.Infoln("LISTEN:", login.addr)
		login.web.ListenAndServe()
	}()
	return true
}

// Close : 关闭
func (login *Login) Close() {
	if login.web != nil {
		login.web.Close()
		login.web = nil
	}
}
