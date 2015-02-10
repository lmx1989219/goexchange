package exchange

import (
	//	"code.google.com/p/goprotobuf/proto"
	l4g "github.com/alecthomas/log4go"
	//	"github.com/pebbe/zmq4"
	//	"time"
	//	. "gw.com.cn/dzhyun/sdk.bus.git"
	. "gw.com.cn/dzhyun/app.frame.git"
	"strconv"
	"sync"

)

//app自定义接口
type CustomPuber struct {
	Dqu *DynaQuoteUtil
	rqh *RequestHandler
	lock    sync.RWMutex
}

func NewCustomPuber() *CustomPuber {
	return &CustomPuber{Dqu: NewDynaQuoteUtil()}
}

func (this *CustomPuber) SetRequestHandler(rq *RequestHandler) {
	this.rqh = rq
}

//基于用户行为做一些逻辑处理，推送成交信息，委托挂单 撤单
func (this *CustomPuber) OnCustom(app *AppMain, data interface{}) {
	switch data.(type) {
	case *RcvWtxdEvent:
		//委托下单处理...
		//l4g.Debug("insert order is begin")
		orderPbObj := data.(*RcvWtxdEvent).Obj
		socket := data.(*RcvWtxdEvent).Socket
		commonOrder := ParseCommonOrder(orderPbObj)
		dynq := this.Dqu.GetDynaQuote(commonOrder.Code)
		//			l4g.Info("oid:%d", commonOrder.Oid)
		if dynq != nil && dynq.Sell[0] != 0  && dynq.Buy[0] != 0 {
			//判断行情是否满足成交，如果满足，推送成交信息,否则加入队列
			if commonOrder.BsFlag == BUY && commonOrder.Price >= dynq.Sell[0] {
				commonOrder.Price = dynq.Sell[0]
				l4g.Info("订单直接成交不经过撮合，订单号：%d,%s:%f", commonOrder.Oid, "价格", dynq.Sell[0])
				Mh.GetMatchQueue(commonOrder.Code).DelOrder(commonOrder.BsFlag)
				App.(*AppMain).AddCustomWork(&PushDealEvent{Obj: commonOrder})
				Reply(socket, 8890)
			} else if commonOrder.BsFlag == SELL && commonOrder.Price <= dynq.Buy[0] {
				commonOrder.Price = dynq.Buy[0]
				l4g.Info("订单直接成交不经过撮合，订单号：%d,%s:%f", commonOrder.Oid, "价格", dynq.Buy[0])
				Mh.GetMatchQueue(commonOrder.Code).DelOrder(commonOrder.BsFlag)
				App.(*AppMain).AddCustomWork(&PushDealEvent{Obj: commonOrder})
				Reply(socket, 8890)
			} else {
				if this.rqh.DoXd(commonOrder) {
					Reply(socket, 8889)
				}
			}
		} else {
			l4g.Info("行情未初始化...%d", commonOrder.Code)
			Reply(socket, 8888)
		}
	case *RcvWtcdEvent:
		//委托撤单处理...
		orderPbObj := data.(*RcvWtcdEvent).Obj
		socket := data.(*RcvWtcdEvent).Socket
		commonOrder := ParseCommonOrder(orderPbObj)
		commonOrder.IsCancle = true
		if this.rqh.DoCd(commonOrder) {
			Reply(socket, 8887)
		}
	case *PushDealEvent:
		this.lock.Lock()
		defer this.lock.Unlock()
		dealobj := data.(*PushDealEvent).Obj
		bus := app.GetBusHelper()
		dealprice := strconv.FormatFloat(float64(dealobj.Price), 'f', -1, 32)
		dealCount := strconv.FormatInt(int64(dealobj.Count), 32)
		jsonResp := `{"Code": "` + dealobj.Code + `",
						"DealPrice":` + dealprice + `,"DealCount":` + dealCount + `,
						"MarketCode":"test","UserName":"lmx01",
						"OrderId":100,"DealType":1,"GameId":"FAFP"}`
		//是否采用推送key的方式？value通过redis取？
		bus.PubTopic(jsonResp)
		dealobj.IsCancle = true
		App.(*AppMain).AddCustomWork(&OrderCopyEvent{Obj: dealobj})
//		l4g.Info("推送成交信息: %s", jsonResp)
	case *RcvQuoteEvent:
		quote := data.(*RcvQuoteEvent).Obj
		this.Dqu.FlushQuoteMap(quote)
	case *MatchEvent:
		quote := data.(*MatchEvent).Obj
		Mh.Match(quote.Code, quote)
	case *OrderCopyEvent:
		co := data.(*OrderCopyEvent).Obj
		AppStore.(*StorageSvc).SaveOrder(co)
	}
}
