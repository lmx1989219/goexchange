package main

import (
<<<<<<< .mine
=======
	"exchange"
>>>>>>> .r8070
	l4g "github.com/alecthomas/log4go"
	. "gw.com.cn/dzhyun/app.frame.git"
	//ex "gw.com.cn/dzhyun/dzhexchange"
	"gw.com.cn/dzhyun/utils.git"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"
)

<<<<<<< .mine
var MatchMarkets []string

=======
var MatchMarkets = [2]string{"SH", "SZ"}

>>>>>>> .r8070
//保证系统仅初始化一次
var once sync.Once

func main() {
	once.Do(initApp)
}

//初始化
func initApp() {
<<<<<<< .mine
	runtime.GOMAXPROCS(GOMAXPROCS)
=======
	//必须声明最大线程数，否则appframe默认为1个线程
	runtime.GOMAXPROCS(exchange.GOMAXPROCS)
>>>>>>> .r8070
	LaxFlag := utils.NewLaxFlagDefault()
	var appcfgfile = LaxFlag.String("lf", "log4go.xml", "app cfgfile")
	LaxFlag.LaxParseDefault()
	l4g.LoadConfiguration(*appcfgfile)

	//业务处理程序
<<<<<<< .mine
	customer := ex.NewCustomPuber()
	matchPool := ex.NewMatchPool()
=======
	customer := exchange.NewCustomPuber()
	time.Sleep(time.Millisecond * 5)
	matchPool := exchange.NewMatchPool()
>>>>>>> .r8070

	for _, v := range MatchMarkets {
		matchPool.BindMarketMatchQueue(v)
	}
	codesArr := ex.LoadData()
	//根据商品代码创建队列
	for i := 0; i < len(codesArr); i++ {
		stkList := new(ex.StkMatchList)
<<<<<<< .mine
		stkList.InitList()
		cb := ex.MarketQuoteMap[codesArr[i]]
		matchPool.GetMatchQueueByMarket(cb.MarketCode)[codesArr[i]] = stkList
=======
		pcode := codesArr[i]
		stkList.InitList(pcode)
		cb := exchange.MarketQuoteMap[pcode]
		matchPool.GetMatchQueueByMarket(cb.MarketCode)[pcode] = stkList
>>>>>>> .r8070
	}
	l4g.Info("撮合池市场大小: %d", len(matchPool.GetAllMatchQueue()))
	ex.Mh = ex.NewMatchHandler()
	//撮合srv
	stms := new(ex.StkMatchService)
	ex.Mh.Init(matchPool, stms)
	ex.Rh = ex.NewReqHandler()
	ex.Rh.SetMatchHandler(ex.Mh)
	customer.SetRequestHandler(ex.Rh)
	customer.Dqu.RecoveryQuoteFromDB()
	ex.CustomerAction = customer
	time.Sleep(time.Millisecond * 5)

	//初始化app
	myapp := ex.NewExchangeApp()
	//初始化应用框架（任务调度，初始化app conf等）
	workmain := NewWorkMain()
	defer func() {
		workmain.Stop()
	}()
	app := NewAppMain(myapp)
	ex.App = app

	app.SetCustom(customer)
	//启动应用框架
	workmain.Start(app)
<<<<<<< .mine
	storeaddress, _ := ex.App.(*AppMain).GetStoreAddr()
=======

	//初始化存储
	storeaddress, _ := exchange.App.(*AppMain).GetStoreAddr()
>>>>>>> .r8070
	storesvc := ex.NewStorageSvc()
	storesvc.InitStore(storeaddress, "", "")
<<<<<<< .mine
	ex.AppStore = storesvc
=======
	var count int64
	//恢复订单到撮合队列
	for _, order := range storesvc.GetOrder() {
		l4g.Info("recover order %v", order)
		cb := exchange.MarketQuoteMap[order.Code]
		matchPool.GetMatchQueueByMarket(cb.MarketCode)[order.Code].InsertOrder(order)
		count++
	}
	l4g.Info("恢复订单 %v %s", count, "笔")
	exchange.AppStore = storesvc
>>>>>>> .r8070

	//定时清理任务
	go exchange.Scheduling()
	l4g.Info("app 启动完成")
	//等待退出
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)
	<-c
	l4g.Info("Receive ctrl-c")
}

/**
 *初始化码表
 *初始化每个市场每个品种的撮合队列
 *初始化市场所属的撮合池
 *加载历史行情 等待行情触发再次撮合
 *通过app框架订阅行情 获取ua请求
 *发布内部事件（下单撤单推送等）
 *监听内部事件 订单持久和恢复/收盘清空队列等暂时未实现
 */
func initCoreSys() {
}
