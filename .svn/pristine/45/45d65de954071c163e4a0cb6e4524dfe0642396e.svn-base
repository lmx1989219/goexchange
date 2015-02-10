package exchange

import (
	l4g "github.com/alecthomas/log4go"
	//	"os"
	//	"time"
	. "gw.com.cn/dzhyun/app.frame.git"
)

type FuMatchService struct {
}

type StkMatchService struct {
}

const (
	BUY  = 1 //买
	SELL = 2 //卖
)

func (stm *StkMatchService) DealMatch(sml QueueService, q *DynaQuote) {
	if nil == sml {
		return
	}
	for {
		order := sml.GetOrder(BUY)
		if order == nil {
			break
		}
		//订单时间大于行情时间不能进行本次撮合,(老行情不能撮合新订单)
		if order.OrderTime >= q.QuoteTime {
			continue
		}
		//成交条件：委托价大于卖一价
		if order.Price > q.Sell[0] {
			l4g.Info("订单已成交，订单号：%d %s:%f %s：%f", order.Oid, "成交价格", q.Sell[0], "委托价格", order.Price)
			sml.DelOrder(order.BsFlag)
			App.(*AppMain).AddCustomWork(&PushDealEvent{Obj: order})
		} else {
			l4g.Info("当前订单不满足成交条件，订单号:%d", order.Oid)
			break
		}
	}
	for {
		order := sml.GetOrder(SELL)
		if order == nil {
			break
		}
		if order.OrderTime > q.QuoteTime {
			continue
		}
		//成交条件：委托价小于买一价
		if order.Price <= q.Buy[0] {
			l4g.Info("订单已成交，订单号：%d %s:%f %s：%f", order.Oid, "成交价格", q.Sell[0], "委托价格", order.Price)
			sml.DelOrder(order.BsFlag)
			App.(*AppMain).AddCustomWork(&PushDealEvent{Obj: order})
		} else {
			l4g.Info("当前订单不满足成交条件，订单号:%d", order.Oid)
			break
		}
	}
}
