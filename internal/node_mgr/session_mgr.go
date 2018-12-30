package nodemgr

import (
	"sync"

	"github.com/fananchong/go-xserver/common"
	"github.com/fananchong/go-xserver/internal/protocol"
	"github.com/fananchong/go-xserver/internal/utility"
)

var (
	xsessionmgr = newSessionMgr()
)

// SessionMgr : 网络会话对象管理类
type SessionMgr struct {
	m  sync.RWMutex
	ss map[common.NodeType][]*protocol.SERVER_INFO
}

func newSessionMgr() *SessionMgr {
	return &SessionMgr{
		ss: make(map[common.NodeType][]*protocol.SERVER_INFO),
	}
}

func (sessmgr *SessionMgr) register(info *protocol.SERVER_INFO) {
	t := common.NodeType(info.GetType())
	sessmgr.m.Lock()
	defer sessmgr.m.Unlock()
	sessmgr.ss[t] = append(sessmgr.ss[t], info)
}

func (sessmgr *SessionMgr) lose(nid common.NodeID, t common.NodeType) {
	sid := utility.NodeID2ServerID(nid)
	sessmgr.m.Lock()
	defer sessmgr.m.Unlock()
	if lst, ok := sessmgr.ss[t]; ok {
		findindex := -1
		for i, v := range lst {
			if v.GetId().GetLow() == sid.GetLow() && v.GetId().GetHigh() == sid.GetHigh() {
				findindex = i
				break
			}
		}
		if findindex >= 0 {
			lst = append(lst[:findindex], lst[findindex+1:]...)
			common.XLOG.Infoln("lose node, type:", t, "id:", utility.NodeID2UUID(nid).String())
		}
	}
}

func (sessmgr *SessionMgr) forByType(t common.NodeType, f func(info *protocol.SERVER_INFO)) {
	var templst []*protocol.SERVER_INFO
	sessmgr.m.RLock()
	if lst, ok := sessmgr.ss[t]; ok {
		templst = append(templst, lst...)
	}
	sessmgr.m.RUnlock()
	for _, v := range templst {
		f(v)
	}
}

func (sessmgr *SessionMgr) forAll(f func(info *protocol.SERVER_INFO)) {
	var templst []*protocol.SERVER_INFO
	sessmgr.m.RLock()
	for _, lst := range sessmgr.ss {
		templst = append(templst, lst...)
	}
	sessmgr.m.RUnlock()
	for _, v := range templst {
		f(v)
	}
}
