package exchange

var MarketQuoteMap = make(map[string]*ContractsBean) //市场-行情

var Mh *MatchHandler

var Rh *RequestHandler

var AppSrv interface{}

var App interface{}

var AppStore interface{}

var CustomerAction interface{}
