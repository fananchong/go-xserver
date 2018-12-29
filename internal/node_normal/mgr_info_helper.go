package nodenormal

import (
	"time"

	"github.com/fananchong/go-xserver/common"
	"github.com/fananchong/go-xserver/internal/db"
)

func getMgrInfoByBlock() (string, int32) {
	common.XLOG.Infoln("Try get mgr server info ...")
	data := db.NewMgrServer(common.XCONFIG.DbMgr.Name, 0)
	for {
		if err := data.Load(); err == nil {
			break
		} else {
			common.XLOG.Debugln(err)
			time.Sleep(1 * time.Second)
		}
	}
	common.XLOG.Infoln("Mgr server address:", data.GetAddr())
	common.XLOG.Infoln("Mgr server port:", data.GetPort())
	return data.GetAddr(), data.GetPort()
}
