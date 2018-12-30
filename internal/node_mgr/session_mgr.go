package nodemgr

import (
	"sync"

	"github.com/fananchong/go-xserver/common"
	"github.com/fananchong/go-xserver/internal/utility"
)

var (
	xsessionmgr = newSessionMgr()
)

// SessionMgr : 网络会话对象管理类
type SessionMgr struct {
	m  sync.RWMutex
	ss map[common.NodeType][]*Session
}

func newSessionMgr() *SessionMgr {
	return &SessionMgr{
		ss: make(map[common.NodeType][]*Session),
	}
}

func (sessmgr *SessionMgr) register(sess *Session) {
	t := sess.GetType()
	sessmgr.m.Lock()
	defer sessmgr.m.Unlock()
	sessmgr.ss[t] = append(sessmgr.ss[t], sess)
}

func (sessmgr *SessionMgr) lose(sess *Session) {
	sessmgr.m.Lock()
	defer sessmgr.m.Unlock()
	if lst, ok := sessmgr.ss[sess.GetType()]; ok {
		findindex := -1
		for i, v := range lst {
			if v == sess {
				findindex = i
				break
			}
		}
		if findindex >= 0 {
			lst = append(lst[:findindex], lst[findindex+1:]...)
			common.XLOG.Infoln("lose node, type:", sess.GetType(), "id:", utility.ServerID2UUID(sess.GetSID()).String())
		}
	}
}

func (sessmgr *SessionMgr) forByType(t common.NodeType, f func(*Session)) {
	var templst []*Session
	sessmgr.m.RLock()
	if lst, ok := sessmgr.ss[t]; ok {
		templst = append(templst, lst...)
	}
	sessmgr.m.RUnlock()
	for _, v := range templst {
		f(v)
	}
}

func (sessmgr *SessionMgr) forAll(f func(*Session)) {
	var templst []*Session
	sessmgr.m.RLock()
	for _, lst := range sessmgr.ss {
		templst = append(templst, lst...)
	}
	sessmgr.m.RUnlock()
	for _, v := range templst {
		f(v)
	}
}
