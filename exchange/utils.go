package exchange

import (
	"code.google.com/p/goprotobuf/proto"
	l4g "github.com/alecthomas/log4go"
	dzhyun "gw.com.cn/dzhyun/dzhyun.git"
	. "gw.com.cn/dzhyun/sdk.bus.git"
	storage "gw.com.cn/dzhyun/sdk.storage.git"
	//	"strconv"
	"time"
)

func Commonorder2Pb(order *CommonOrder) ([]byte, error) {
	orderPb := &CommonOrderPb{
		Oid:       proto.Int64(order.Oid),
		Price:     proto.Float64(order.Price),
		Code:      proto.String(order.Code),
		BsFlag:    proto.Int32(order.BsFlag),
		Count:     proto.Int32(order.Count),
		OrderTime: proto.String(order.OrderTime),
		IsCancle:  proto.Bool(order.IsCancle),
		Position:  proto.Int64(order.Position),
	}
	return proto.Marshal(orderPb)
}

func Pb2Commonorder(data []byte) (*CommonOrder, error) {
	orderpb := &CommonOrderPb{}
	e := proto.Unmarshal(data, orderpb)
	t, _ := time.Parse(time.RFC3339, orderpb.GetOrderTime())
	order := &CommonOrder{
		Oid:       orderpb.GetOid(),
		Price:     orderpb.GetPrice(),
		Code:      orderpb.GetCode(),
		BsFlag:    orderpb.GetBsFlag(),
		Count:     orderpb.GetCount(),
		OrderTime: t.String(),
		IsCancle:  orderpb.GetIsCancle(),
		Position:  orderpb.GetPosition(),
	}
	return order, e
}

func CreateRequestData(json *JsonNode) *C2M_JSON {
	pcode := json.GetNodeByName("ProductCode").ValueString
	marketcode := json.GetNodeByName("MarketCode").ValueString
	oid := int64(json.GetNodeByName("OrderId").ValueNumber)
	wtprice := float64(json.GetNodeByName("WtPrice").ValueNumber)
	bf := int32(json.GetNodeByName("BsFlag").ValueNumber)
	wtcount := int32(json.GetNodeByName("WtCount").ValueNumber)
	wttime := json.GetNodeByName("WtTime").ValueString
	ntime, _ := time.Parse(time.RFC3339, wttime)
	return &C2M_JSON{
		ProductCode: pcode,
		OrderId:     oid,
		BsFlag:      bf,
		WtPrice:     wtprice,
		MarketCode:  marketcode,
		WtCount:     wtcount,
		WTtime:      ntime.String(),
	}
}

const (
	Mmkey = "MaimaipanDangri"
)

//redis行情obj
func ParsePbQuote(store *storage.Store, data []byte, topic string) *DynaQuote {
	quoteobj := &dzhyun.GeguDongtaiDangri{}
	proto.Unmarshal(data, quoteobj)
	qtime := quoteobj.GetTime()
	np := float64(quoteobj.GetPrice()) / 100
	buy, sell := AppStore.(*StorageSvc).GetGegudongtaiMaimai5(store, Mmkey+Substr(topic, 17, 8))
	dq := &DynaQuote{
		Code:       Substr(topic, 19, 6),
		MarketCode: Substr(topic, 17, 2),
		NewPrice:   np,
		Buy:        buy,
		Sell:       sell,
		QuoteTime:  time.Unix(int64(qtime), 0).Format(time.RFC3339), //"2006-01-02 15:04:05"),
	}
	l4g.Debug("quote:%v", dq)
	return dq
}

//买卖5档
func ParseMainmaiQuote(data []byte, topic string) (buy, sell []float64) {
	mm := &dzhyun.MaimaipanDangri{}
	proto.Unmarshal(data, mm)
	buy = make([]float64, 5)

	buy[0] = float64(mm.GetWeituobuyjia1()) / 100
	buy[1] = float64(mm.GetWeituobuyjia2()) / 100
	buy[2] = float64(mm.GetWeituobuyjia3()) / 100
	buy[3] = float64(mm.GetWeituobuyjia4()) / 100
	buy[4] = float64(mm.GetWeituobuyjia5()) / 100

	sell = make([]float64, 5)
	sell[0] = float64(mm.GetWeituoselljia1()) / 100
	sell[1] = float64(mm.GetWeituoselljia2()) / 100
	sell[2] = float64(mm.GetWeituoselljia3()) / 100
	sell[3] = float64(mm.GetWeituoselljia4()) / 100
	sell[4] = float64(mm.GetWeituoselljia5()) / 100

	l4g.Debug("buy 5: %v %s:%v", buy, "sell 5", sell)
	return buy, sell
}

func ParseCommonOrder(orderPbObj *C2M_JSON) *CommonOrder {
	commonOrder := new(CommonOrder)
	commonOrder.Oid = orderPbObj.OrderId
	commonOrder.Code = orderPbObj.ProductCode
	commonOrder.Price = orderPbObj.WtPrice
	commonOrder.BsFlag = orderPbObj.BsFlag
	commonOrder.Count = orderPbObj.WtCount
	commonOrder.OrderTime = orderPbObj.WTtime
	commonOrder.MarketCode = orderPbObj.MarketCode
	commonOrder.IsCancle = false
	commonOrder.Position = -10
	return commonOrder
}

func Reply(socket *BusResponse, respCode uint32) {
	data := new(dzhyun.BusResponse)
	data.Result = proto.Uint32(respCode)
	data.GetDataResult = make([]byte, 1)
	socket.Send(data)
}

func Substr(str string, start, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0
	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length
	if start > end {
		start, end = end, start
	}
	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}
	return string(rs[start:end])
}

func Scheduling() {
	for {
		now := time.Now()
		if now.Hour() == 16 && now.Minute() == 1 {
			AppStore.(*StorageSvc).RemoveAllOrders()
			l4g.Info("redis备份订单已经清空")
			for _,v := range MarketQuoteMap{
				Mh.GetMatchQueue(v.Code).RemoveAll()
			}
			time.Sleep(time.Second * 60)
		}
		time.Sleep(time.Second)
	}
}
