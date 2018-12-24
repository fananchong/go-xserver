package nodenormal

import (
	"time"

	"github.com/fananchong/go-xserver/common"
	"github.com/fananchong/go-xserver/internal/db"
)

func getMgrInfoByBlock() (string, uint16) {
	data := db.NewMgrServer(common.XCONFIG.DbMgr.Name, 0)
	for {
		if err := data.Load(); err == nil {
			break
		} else {
			time.Sleep(1 * time.Second)
		}
	}
	return data.GetAddr(), data.GetPort()
}
