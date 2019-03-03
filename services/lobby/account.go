package main

import (
	"fmt"
	"reflect"
	"runtime/debug"
	"sync/atomic"

	go_redis_orm "github.com/fananchong/go-redis-orm.v2"
	"github.com/fananchong/go-xserver/common"
	"github.com/fananchong/go-xserver/services/internal/db"
)

// LimitRoleNum : 一个账号角色数限制，最大为 1 个。
const LimitRoleNum = 1

// Account : 账号类
type Account struct {
	*db.RoleList
	account   string
	roles     []*Role
	sess      common.INode
	chanMsg   chan ChanMsg
	chanClose chan int
	closeFlag int32
}

// NewAccount : 角色列表类构造函数
func NewAccount(account string) *Account {
	accountObj := &Account{
		account:   account,
		chanMsg:   make(chan ChanMsg, 1024),
		chanClose: make(chan int, 1),
		RoleList:  db.NewRoleList(Ctx.Config.DbAccount.Name, account),
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
	var ids [256]uint64
	ids[0] = accountObj.RoleList.GetSlot0()
	ids[1] = accountObj.RoleList.GetSlot1()
	ids[2] = accountObj.RoleList.GetSlot2()
	ids[3] = accountObj.RoleList.GetSlot3()
	ids[4] = accountObj.RoleList.GetSlot4()
	for i := 0; i < LimitRoleNum; i++ {
		roleID := ids[i]
		if roleID == 0 {
			// 没有角色，则生成角色ID
			id, err := lobby.NewID(db.IDGenTypeRole)
			if err != nil {
				return err
			}
			roleID = id
			// 角色创建时，这里才用下`反射`。因此不要当心`反射`造成的效率问题
			v := reflect.ValueOf(accountObj.RoleList)
			f := v.MethodByName(fmt.Sprintf("SetSlot%d", i))
			f.Call([]reflect.Value{reflect.ValueOf(roleID)})
		}
		role, err := NewRole(roleID, account)
		if err != nil {
			return err
		}
		accountObj.roles[i] = role
	}
	if err := accountObj.Save(); err != nil {
		return err
	}
	return nil
}

// Start : 开始
func (accountObj *Account) Start() {
	Ctx.Log.Infoln("Account work coroutine start, account:", accountObj.account)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				Ctx.Log.Errorln("[except] ", err, "\n", string(debug.Stack()))
			}
		}()
		for {
			select {
			case msg := <-accountObj.chanMsg:
				accountObj.processMsg(msg.Cmd, msg.Data)
			case <-accountObj.chanClose:
				return
			}
		}
	}()
}

// Close : 结束
func (accountObj *Account) Close() {
	Ctx.Log.Infoln("Account work coroutine close#1, account:", accountObj.account)
	atomic.StoreInt32(&accountObj.closeFlag, 1)
	accountObj.chanClose <- 1
	Ctx.Log.Infoln("Account work coroutine close#2, account:", accountObj.account)
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

// SetSession : 设置账号对应的网络会话
func (accountObj *Account) SetSession(sess common.INode) {
	accountObj.sess = sess
}
