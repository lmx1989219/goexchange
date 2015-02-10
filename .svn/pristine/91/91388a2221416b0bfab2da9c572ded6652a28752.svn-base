package exchange

import (
	logger "github.com/alecthomas/log4go"
	. "gw.com.cn/dzhyun/app.frame.git"
	sdk "gw.com.cn/dzhyun/sdk.bus.git"
	"strings"
)

const (
	SuberKey1       = "GeguDongtaiDangriSH"
	SuberKey2       = "GeguDongtaiDangriSZ"
	GOMAXPROCS      = 64 //保持和存储最大连接数相当
	EVERYGOCHANSIZE = 1024
)

/*
 *	app主框架实现类
 *
 */
type ExchangeApp struct {
}

func NewExchangeApp() *ExchangeApp {
	ea := &ExchangeApp{}
	return ea
}

//初始化APP
func (self *ExchangeApp) OnInit() (param *InitParam) {
	wk := &WorkParam{GOMAXPROCS, EVERYGOCHANSIZE} //工作线程数 缓存区大小
	return &InitParam{Work: *wk, Login: "bus/login?ServiceName=/root/app/dzhexchange"}

}

//初始化APP后，处理数据需求
func (self *ExchangeApp) OnData(app *AppMain) {
	//准备订阅key
	app.SubTopic(SuberKey1)
	app.SubTopic(SuberKey2)

}

//关闭APP，释放chan等
func (self *ExchangeApp) OnClose() {

}

//处理请求:UA请求
func (self *ExchangeApp) OnRequest(request *sdk.BusRequest) {
	if nil == AppStore {
		return
	}
	urlstr := request.GetData()
	logger.Debug("ExchangeApp::OnRequest Data:", urlstr)
	body := strings.Split(urlstr, "?")[1]
	expectedStr := `` + body + ``
	json, err := LoadByString(expectedStr)
	if err != nil {
		logger.Info(err, request)
	}
	reqno := json.GetNodeByName("ProtoNO").ValueString
	c2m := CreateRequestData(json)
	response := request.CreateResponse()
	if reqno == WTXD { //委托下单
		App.(*AppMain).AddCustomWork(&RcvWtxdEvent{Obj: c2m, Socket: response})
	} else if reqno == WTCD { //委托撤单
		App.(*AppMain).AddCustomWork(&RcvWtcdEvent{Obj: c2m, Socket: response})
	}
}

//处理订阅:需要接入行情
func (self *ExchangeApp) OnTopic(topic string) {
	logger.Debug("ExchangeApp::OnTopicData: %s", topic)
	//订阅关注的key
	skipkey := Substr(topic, 19, 6)
	if MarketQuoteMap[skipkey] == nil {
		return
	}
	//去存储取行情obj
	if nil != AppStore {
		AppStore.(*StorageSvc).GetGegudongtaiObj(topic)
	}
}
