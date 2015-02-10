package exchange

import (
	//	zmq "github.com/pebbe/zmq4"
	. "gw.com.cn/dzhyun/sdk.bus.git"
)

type MatchService interface {
	//撮合逻辑处理
	DealMatch(list QueueService, q *DynaQuote)
}

type QueueService interface {
	//初始化队列
	InitList(pcode string)
	//队列标识
	Identify() string
	//订单入队列
	InsertOrder(o *CommonOrder)
	//删除元素
	DelOrder(bs int32)
	//获取订单
	GetOrder(bs int32) (c *CommonOrder)
	//删除指定元素
	RemoveOrderByObj(c *CommonOrder) bool
	//清空队列
	RemoveAll()
}

const (
	RCV_QUOTE_EVENT    = 1  //定义接收行情事件
	PUSH_DEAL_EVENT    = 2  //定义推送成交信息事件
	RCV_REQ_WTXD_EVENT = 3  //定义接收委托请求下单事件
	RCV_REQ_WTCD_EVENT = 3  //定义接收委托请求撤单事件
	EVENT_DOFAILED     = -1 //事件处理失败
	EVENT_DOSUCCED     = 0  //事件处理成功
	MATCH_ORDER        = 5
	COPY_ORDER        = 6
)

type BaseEvent interface {
	//事件标识
	EventIdentify() int
}

type RcvQuoteEvent struct {
	Obj *DynaQuote
}

func (rqe *RcvQuoteEvent) EventIdentify() int {
	return RCV_QUOTE_EVENT
}

type MatchEvent struct {
	Obj *DynaQuote
}

func (rqe *MatchEvent) EventIdentify() int {
	return RCV_QUOTE_EVENT
}

type PushDealEvent struct {
	Obj *CommonOrder
}

func (rqe *PushDealEvent) EventIdentify() int {
	return PUSH_DEAL_EVENT
}

type RcvWtxdEvent struct {
	Obj    *C2M_JSON
	Socket *BusResponse
}

func (rqe *RcvWtxdEvent) EventIdentify() int {
	return RCV_REQ_WTXD_EVENT
}

type RcvWtcdEvent struct {
	Obj    *C2M_JSON
	Socket *BusResponse
}

func (rqe *RcvWtcdEvent) EventIdentify() int {
	return RCV_REQ_WTCD_EVENT
}

type OrderCopyEvent struct {
	Obj *CommonOrder
}

func (rqe *OrderCopyEvent) EventIdentify() int {
	return COPY_ORDER
}
