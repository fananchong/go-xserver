package main

import (
	"sync"

	"github.com/fananchong/go-xserver/services/internal/protocol"
	"github.com/fananchong/go-xserver/services/internal/utility"
)

// AccountMgr : 账号管理类
type AccountMgr struct {
	accounts map[string]*Account
	mutex    sync.RWMutex
}

// NewAccountMgr : 账号管理类构造函数
func NewAccountMgr() *AccountMgr {
	accountMgr := &AccountMgr{
		accounts: make(map[string]*Account),
	}
	return accountMgr
}

// AddAccount : 加入一个账号
func (accountMgr *AccountMgr) AddAccount(account string) *Account {
	accountMgr.mutex.Lock()
	defer accountMgr.mutex.Unlock()
	if old, ok := accountMgr.accounts[account]; ok {
		return old
	}
	accountObj := NewAccount(account)
	accountMgr.accounts[account] = accountObj
	return accountObj
}

// GetAccount : 获取一个账号
func (accountMgr *AccountMgr) GetAccount(account string) *Account {
	accountMgr.mutex.RLock()
	defer accountMgr.mutex.RUnlock()
	if accountObj, ok := accountMgr.accounts[account]; ok {
		if accountObj.Inited() {
			return accountObj
		}
	}
	return nil
}

// DelAccount : 删除一个账号
func (accountMgr *AccountMgr) DelAccount(account string) {
	accountMgr.mutex.Lock()
	defer accountMgr.mutex.Unlock()
	if accountObj, ok := accountMgr.accounts[account]; ok {
		accountObj.Close()
		delete(accountMgr.accounts, account)
	}
}

// PostMsg : 推送消息
func (accountMgr *AccountMgr) PostMsg(account string, cmd uint64, data []byte, flag uint8) {
	var accountObj *Account
	if protocol.CMD_LOBBY_ENUM(cmd) == protocol.CMD_LOBBY_LOGIN {
		accountObj = accountMgr.onLogin(account)
	} else {
		accountObj = accountMgr.GetAccount(account)
	}
	if accountObj == nil {
		Ctx.Errorln("Account object does not exist. account:", account, "cmd:", cmd)
		return
	}
	if accountObj.IsColse() {
		Ctx.Errorln("The account's work coroutine has been closed. account:", account, "cmd:", cmd)
		return
	}
	accountObj.PostMsg(cmd, data, flag)
}

func (accountMgr *AccountMgr) onLogin(account string) *Account {
	if accountObj := accountMgr.AddAccount(account); accountObj != nil {
		if !accountObj.Inited() {
			var err error
			if err = accountObj.Init(); err == nil {
				accountObj.Start()
				return accountObj
			}
			Ctx.Errorln("[LOGIN LOBBY]", err, "account:", account)
			msg := &protocol.MSG_LOBBY_LOGIN_RESULT{}
			msg.Err = protocol.ENUM_LOBBY_COMMON_ERROR_SYSTEM_ERROR
			utility.SendMsgToClient(Ctx, account, uint64(protocol.CMD_LOBBY_LOGIN), msg)
		} else {
			return accountObj
		}
	}
	Ctx.Errorln("[LOGIN LOBBY] New account fail, account:", account)
	return nil
}

func (accountMgr *AccountMgr) onLogout(account string) {

	// TODO: 账号登出处理，需要被触发

	accountMgr.DelAccount(account)
}
