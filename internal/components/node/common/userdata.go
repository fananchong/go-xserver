package nodecommon

import (
	"github.com/fananchong/go-xserver/common"
)

// UserData : 传递给 Session 的数据
type UserData struct {
	Ctx     *common.Context
	SessMgr *SessionMgr
}
