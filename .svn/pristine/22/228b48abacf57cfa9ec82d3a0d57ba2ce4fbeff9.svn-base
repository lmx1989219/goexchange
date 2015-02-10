package exchange

import ()

type MatchPool struct {
	matchQueueMap map[string]map[string]QueueService
}

func NewMatchPool() *MatchPool {
	return &MatchPool{matchQueueMap: make(map[string]map[string]QueueService)}
}

func (this MatchPool) BindMarketMatchQueue(market string) {
	qs := make(map[string]QueueService)
	this.matchQueueMap[market] = qs
}

func (this MatchPool) GetMatchQueueByMarket(market string) map[string]QueueService {
	return this.matchQueueMap[market]
}

func (this MatchPool) GetAllMatchQueue() map[string]map[string]QueueService {
	return this.matchQueueMap
}
