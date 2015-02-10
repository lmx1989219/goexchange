package exchange

import (
	logger "github.com/alecthomas/log4go"
	frame "gw.com.cn/dzhyun/app.frame.git"
	. "gw.com.cn/dzhyun/sdk.storage.git"
)

const (
	poskey = "ORDER_POS_KEY" //存储订单pos，key+pcode
	appSelfHash = "EXCHANGE_DATA1_"
)

type PosObj struct{
	key string
	value int64
}
 
type StorageSvc struct {
	storePool *Pool
}

func NewStorageSvc() *StorageSvc {
	ea := &StorageSvc{}
	return ea
}

func (self *StorageSvc) GetStorePool() *Pool {
	return self.storePool

}

func (self *StorageSvc) InitStore(address, u, p string) {
	self.storePool = NewPool(address, u, p)
	logger.Debug("ExchangeApp::存储: %v", self.storePool)
}

func (self *StorageSvc) GetGegudongtaiObj(key string) {
	store := self.storePool.GetStore("")
	defer self.storePool.Release(store)
//	logger.Debug("存储 store:%v", store)
	if nil != store {
		oplist := store.NewListOperator(key)
//		logger.Debug("存储 oplist:%v", oplist)
		//-1 取最新的行情
		data, e := oplist.Get(-1)
		if nil != e {
			logger.Error(e)
		}
		dq := ParsePbQuote(store, data, key)
		logger.Debug("rcv dync quote %s %s %f", dq.Code, dq.MarketCode, dq.NewPrice)
		App.(*frame.AppMain).AddCustomWork(&RcvQuoteEvent{Obj: dq})
	}
}

func (self *StorageSvc) GetGegudongtaiMaimai5(store *Store, key string) (b, s []float64) {
	oplist := store.NewListOperator(key)
	data, e := oplist.Get(-1)
	if nil != e {
		logger.Error(e)
	}
	return ParseMainmaiQuote(data, key)
}

////保存订单用于恢复
//func (self *StorageSvc) SaveOrder(order *CommonOrder) {
//	store := self.storePool.GetStore("")
//	defer self.storePool.Release(store)
//	data, _ := Commonorder2Pb(order)
//	oplist := store.NewListOperator(order.Code)
//	erro := oplist.Add(data)
//	if nil != erro {
//		logger.Error(erro)
//	}
//}

///*
//* order需要pb序列化
//* @pos 当前位置，如果有重复则作为index更新,下单需要保存pos，撤单不用
//*/
//func (self *StorageSvc) SaveOrder(order *CommonOrder) int64{
//	store := self.storePool.GetStore("")
//	defer self.storePool.Release(store)
//	oplist := store.NewListOperator(order.Code)
//	obj, _ := oplist.Get(order.Position)
//	if nil == obj {
//		data, e := Commonorder2Pb(order)
//		if nil != e {
//			logger.Error(e)
//		}
//		oplist.Add(data)
//		curpos,_:=oplist.Count()
//		return curpos
//	}else{
//		co, _ := Pb2Commonorder(obj)
//		if co.Oid == order.Oid {
//			updateobj, _ := Commonorder2Pb(order)
//			erro := oplist.Replace(order.Position, updateobj)
//			if nil != erro {
//				logger.Error(erro)
//			}
//		}
//	}
//	return -10
//}
//
////order需要pb序列化
//func (self *StorageSvc) RemoveOrder(order *CommonOrder) {
//	store := self.storePool.GetStore("")
//	defer self.storePool.Release(store)
//	oplist := store.NewListOperator(order.Code)
//	lists, _ := oplist.GetRange(0, -1)
//	var index int64
//	for idx, _ := range lists {
//		//如果当前obj=order 返回索引值
//		curobj := lists[idx]
//		co, _ := Pb2Commonorder(curobj)
//		if co.Oid == order.Oid {
//			index = int64(idx)
//			break
//		}
//	}
//	erro := oplist.Remove(index)
//	if nil != erro {
//		logger.Error(erro)
//	}
//}
//
////查询已报订单
//func (self *StorageSvc) GetOrder(key string) []*CommonOrder {
//	store := self.storePool.GetStore("")
//	defer self.storePool.Release(store)
//	cos := []*CommonOrder{}
//	oplist := store.NewListOperator(key)
//	lists, _ := oplist.GetRange(0, -1)
//	leng := 0
//	for idx, _ := range lists {
//		objbytes := lists[idx]
//		co, _ := Pb2Commonorder(objbytes)
//		//pb反序列化
//		cos = append(cos, co)
//		leng++
//	}
//	return cos[:leng]
//}
//
////删除全部订单
//func (self *StorageSvc) RemoveAllOrders(key string){
//	store := self.storePool.GetStore("")
//	defer self.storePool.Release(store)
//	oplist := store.NewListOperator(key)
//	erro := oplist.RemoveAll()
//	if nil != erro {
//		logger.Error(erro)
//	}
//}

/*
* order需要pb序列化
*/
func (self *StorageSvc) SaveOrder(order *CommonOrder) {
	store := self.storePool.GetStore("")
	defer self.storePool.Release(store)
	ophash := store.NewHashOperator(appSelfHash)
	data, e := Commonorder2Pb(order)
	if nil != e {
		logger.Error(e)
	}
	err := ophash.Set(order.Oid,data)
	if nil != err {
		logger.Error(e)
	}
}


//查询已报订单
func (self *StorageSvc) GetOrder() []*CommonOrder {
	store := self.storePool.GetStore("")
	defer self.storePool.Release(store)
	cos := []*CommonOrder{}
	ophash := store.NewHashOperator(appSelfHash)
	lists, _ := ophash.Keys()
	leng := 0
	for idx, _ := range lists {
		objbytes ,_ := ophash.Get(string(lists[idx]))
		//pb反序列化
		co, _ := Pb2Commonorder(objbytes)
		if co.IsCancle == true {
			continue
		}
		cos = append(cos, co)
		leng++
	}
	return cos[:leng]
}

//删除全部订单
func (self *StorageSvc) RemoveAllOrders(){
	store := self.storePool.GetStore("")
	defer self.storePool.Release(store)
	ophash := store.NewHashOperator(appSelfHash)
	erro := ophash.RemoveAll()
	if nil != erro {
		logger.Error(erro)
	}
}