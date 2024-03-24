package acrnm

import (
	"net/http"
	"time"

	"github.com/Drelf2018/asyncio"
	"github.com/Drelf2018/event"
	"github.com/Drelf2018/request"
	"github.com/Drelf2020/utils"
)

var log = utils.SetOutputFile("acrnm.log")

func init() {
	utils.SetTimestampFormat(time.DateTime)
}

const (
	url string = "http://fastapi.web-framework-2ml5.1990019364850918.cn-hangzhou.fc.devsapp.net"
	New        = "new"
)

var (
	session, _ = request.NewTypeSession[[]*Product](http.MethodGet, url)
	m          = NewMap(session.Must()...)
)

func GetProducts() []*Product {
	return session.Must()
}

func Run[T Sender](s T) {
	log.Info("开始运行")
	event.Heartbeat(1000, 7500, func(b *event.Beat) {
		if b.Count%1440 == 0 {
			log.Infof("运行中#%d，当前容量：%d", b.Count, m.Len())
		}
		newVals, _, _ := m.Updates(session.Must())
		for _, n := range newVals {
			log.Info(n)
			ok := asyncio.Retry(3, 1, func() bool { return s.Send(n) == nil })
			if !ok {
				log.Errorf("发送失败：%v", n)
			}
		}
	})
}
