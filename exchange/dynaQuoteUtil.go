package exchange

import (
	"database/sql"
	l4g "github.com/alecthomas/log4go"
	//	dzhyun "gw.com.cn/dzhyun/dzhyun.git"
	. "gw.com.cn/dzhyun/app.frame.git"
	_ "mysql"
	"strconv"
	"strings"
	"sync"
)

/**
* 刷新动态行情类
 */
type DynaQuoteUtil struct {
	quotes map[string]*DynaQuote //行情表
	lock    sync.RWMutex
}

const (
	quoteBuffersize = 1000
)

func NewDynaQuoteUtil() *DynaQuoteUtil {
	return &DynaQuoteUtil{
		quotes: make(map[string]*DynaQuote),
	}
}

func (dqu *DynaQuoteUtil) FlushQuoteMap(quote *DynaQuote) {
	//刷新内存应该是同步操作
	dqu.lock.Lock()
	defer dqu.lock.Unlock()
	dqu.quotes[quote.Code] = quote
	//发布一个事件:该合约队列执行撮合
	if quote.Sell[0] != 0 && quote.Buy[0] != 0 {
		l4g.Debug("撮合开始执行 市场：%s 合约:%s", quote.MarketCode, quote.Code)
		App.(*AppMain).AddCustomWork(&MatchEvent{Obj: quote})
	}
	//l4g.Debug("当前动态行情表大小：",len(dqu.quotes))
}

func (dqu *DynaQuoteUtil) GetDynaQuote(pcode string) *DynaQuote {
	//	l4g.Debug("当前动态行情表大小：", len(dqu.quotes), pcode, dqu.quotes[pcode])
	return dqu.quotes[pcode]
}

func (dqu *DynaQuoteUtil) RecoveryQuoteFromDB() {
	db, err := sql.Open("mysql", "root:601519@tcp(10.15.107.105:3306)/mobile_counter_fafp?charset=utf8")
	CheckErr(err)
	defer db.Close()
	rows, err := db.Query("select * from T_SLIVER_QUOTE_DETAIL")
	CheckErr(err)
	for rows.Next() {
		var NAME string
		var PRODUCT_CODE string
		var QUOTE_TIME string
		var MARKET_CODE string
		var NEW_PRICE float64
		var BUY_PRICE string
		var SELL_PRICE string
		var ADVSTOP float64
		var DECSTOP float64
		var SETTLEPRICE float64
		var LASTCLOSE float64
		e := rows.Scan(&NAME, &PRODUCT_CODE, &QUOTE_TIME, &NEW_PRICE, &BUY_PRICE, &SELL_PRICE, &MARKET_CODE, &ADVSTOP, &DECSTOP, &SETTLEPRICE, &LASTCLOSE)
		CheckErr(e)
		dq := new(DynaQuote)
		dq.Code = PRODUCT_CODE
		dq.Name = NAME
		dq.NewPrice = NEW_PRICE
		buy5 := strings.Split(BUY_PRICE, ",")
		buy5f := []float64{}
		for _, v := range buy5 {
			vl, _ := strconv.ParseFloat(v, 64)
			buy5f = append(buy5f, vl)
		}
		sell5 := strings.Split(SELL_PRICE, ",")
		sell5f := []float64{}
		for _, v := range sell5 {
			vl, _ := strconv.ParseFloat(v, 64)
			sell5f = append(sell5f, vl)
		}
		dq.Buy = buy5f[:5]
		dq.Sell = sell5f[:5]
		dqu.quotes[PRODUCT_CODE] = dq
	}
	l4g.Info("行情加载完成: %d", len(dqu.quotes))
}
