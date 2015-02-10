package exchange

import (
	l4g "github.com/alecthomas/log4go"
)

/**
*	带市场撮合池 key market ,value queueMap(key code,value Queue)
 */
type MatchHandler struct {
	matchPool *MatchPool
	ms              MatchService
}

func (self *MatchHandler) Init(pool *MatchPool,ms_ MatchService) {
	self.matchPool = pool
	self.ms = ms_
	l4g.Debug("撮合池:%v", self)
}

func NewMatchHandler() *MatchHandler {
	//撮合服务
	return &MatchHandler{}
}

func (self *MatchHandler) GetMatchQueue(code string) QueueService {
	cb := MarketQuoteMap[code]
	if nil == cb {
		return nil
	}
	mk := MarketQuoteMap[code].MarketCode
//	l4g.Debug("撮合池:%v", self)
	q := self.matchPool.GetMatchQueueByMarket(mk)
	return q[code]
}

func (self *MatchHandler) Match(code string, q *DynaQuote) {
	l4g.Debug("执行撮合:%v", self)
	queue := self.GetMatchQueue(code)
	self.ms.DealMatch(queue, q)
}
