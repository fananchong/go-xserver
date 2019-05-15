package main

import (
	"fmt"
	"runtime/debug"
	"sync/atomic"

	go_redis_orm "github.com/fananchong/go-redis-orm.v2"
	"github.com/fananchong/go-xserver/services/internal/db"
)

// LimitRoleNum : 一个账号角色数限制，最大为 1 个。
const LimitRoleNum = 1

// Account : 账号类
type Account struct {
	*db.RoleList
	account     string       // 账号字符串
	roles       []*Role      // 角色列表
	chanMsg     chan ChanMsg // 用于接收该账号消息
	chanClose   chan int     // 用于接收该账号消息循环
	closeFlag   int32        // 标志该账号是否结束会话
	currentRole *Role        // 该账号当前角色
	inited      bool         // 该账号是否已经初始化完毕
}

// NewAccount : 角色列表类构造函数
func NewAccount(account string) *Account {
	accountObj := &Account{
		account:   account,
		chanMsg:   make(chan ChanMsg, 1024),
		chanClose: make(chan int, 1),
		RoleList:  db.NewRoleList(Ctx.Config().DbAccount.Name, account),
	}
	return accountObj
}

// Init : 账号初始化
func (accountObj *Account) Init() error {
	account := accountObj.account
	if err := accountObj.RoleList.Load(); err != nil {
		if err != go_redis_orm.ERR_ISNOT_EXIST_KEY {
			return err
		}
		if accountObj.FirstInitialization() == false {
			return fmt.Errorf("Account failed to initialize for the first time, account:%s", account)
		}
	}
	accountObj.roles = make([]*Role, LimitRoleNum)
	ids := accountObj.RoleList.GetRoles(false).GetRoleIDs()
	for i := 0; i < LimitRoleNum; i++ {
		if roleID, ok := ids[uint32(i)]; ok {
			role, err := NewRole(roleID, account)
			if err != nil {
				return err
			}
			accountObj.roles[i] = role
		} else {
			accountObj.roles[i] = nil
		}
	}
	accountObj.inited = true
	return nil
}

// Start : 开始
func (accountObj *Account) Start() {
	Ctx.Infoln("Account work coroutine start, account:", accountObj.account)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				Ctx.Errorln("[except] ", err, "\n", string(debug.Stack()))
			}
		}()
		for {
			select {
			case msg := <-accountObj.chanMsg:
				accountObj.processMsg(msg.Cmd, msg.Data, msg.Flag)
			case <-accountObj.chanClose:
				return
			}
		}
	}()
}

// Close : 结束
func (accountObj *Account) Close() {
	Ctx.Infoln("Account work coroutine close#1, account:", accountObj.account)
	atomic.StoreInt32(&accountObj.closeFlag, 1)
	accountObj.chanClose <- 1
	Ctx.Infoln("Account work coroutine close#2, account:", accountObj.account)
}

// FirstInitialization : 账号首次创建初始化
func (accountObj *Account) FirstInitialization() bool {
	return true
}

// IsColse : 是否已经关闭
func (accountObj *Account) IsColse() bool {
	flag := atomic.LoadInt32(&accountObj.closeFlag)
	return flag != 0
}

// GetRoles : 获取账号对应的角色列表
func (accountObj *Account) GetRoles() []*Role {
	return accountObj.roles
}

// Inited : 是否已经初始化
func (accountObj *Account) Inited() bool {
	return accountObj.inited
}

// GetRole : 获取账号当前角色
func (accountObj *Account) GetRole() *Role {
	return accountObj.currentRole
}
