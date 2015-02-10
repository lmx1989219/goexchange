package exchange

import (
	l4g "github.com/alecthomas/log4go"
	. "gw.com.cn/dzhyun/app.frame.git"
)

type RequestHandler struct {
	mh      *MatchHandler
	posMaps map[int64]int64
}

func NewReqHandler() *RequestHandler {
	return &RequestHandler{posMaps: make(map[int64]int64)}
}

func (self *RequestHandler) SetPos(co *CommonOrder) {
	self.posMaps[co.Oid] = co.Position
}

func (self *RequestHandler) SetMatchHandler(mhandler *MatchHandler) {
	self.mh = mhandler
}

//委托下单
func (self *RequestHandler) DoXd(co *CommonOrder) bool {
	l4g.Debug("下单: %s ", co.Code)
	App.(*AppMain).AddCustomWork(&OrderCopyEvent{Obj: co})
	self.mh.GetMatchQueue(co.Code).InsertOrder(co)
	return true
}

//委托撤单
func (self *RequestHandler) DoCd(co *CommonOrder) bool {
	l4g.Debug("撤单: %s ", co.Code)
	App.(*AppMain).AddCustomWork(&OrderCopyEvent{Obj: co})
	return self.mh.GetMatchQueue(co.Code).RemoveOrderByObj(co)
}