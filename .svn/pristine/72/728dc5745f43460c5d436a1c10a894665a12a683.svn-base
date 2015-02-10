package exchange

import (
	"container/list"
	l4g "github.com/alecthomas/log4go"
	"sync"
	"time"
)

type StkMatchList struct {
	List MatchList
	l    sync.RWMutex //single write ,more read  slock
	//记录每个节点指针，用于直接删除
	codesMap map[int64]*list.Element
	code string
}

func (stm *StkMatchList) InitList(pcode string) {
	stm.code = pcode
	stm.codesMap = make(map[int64]*list.Element)
}

func (stm *StkMatchList) InsertOrder(o *CommonOrder) {
	stm.l.Lock()
	defer stm.l.Unlock()
	begin := time.Now()
	if o.BsFlag == BUY { //买 价格从低到高 时间由新及旧
		var el *list.Element
		if stm.List.BuyList.Len() == 0 { //list为空
			stm.List.BuyList.PushBack(o)
			stm.codesMap[o.Oid] = stm.List.BuyList.Front()
		} else {
			e := stm.List.BuyList.Front()
			for {
				if e.Value.(*CommonOrder).Price > o.Price {
					stm.List.BuyList.InsertBefore(o, e)
					stm.codesMap[o.Oid] = e
					break
				} else if e.Value.(*CommonOrder).Price == o.Price { //相同价格比较时间
					if e.Value.(*CommonOrder).OrderTime < o.OrderTime {
						stm.List.BuyList.InsertBefore(o, e)
						stm.codesMap[o.Oid] = e
						break
					}
				}
				el = e.Next()
				if el == nil { //结束
					stm.List.BuyList.InsertAfter(o, e)
					stm.codesMap[o.Oid] = e
					break
				}
				e = el
			}
		}

	} else if o.BsFlag == SELL { //卖  价格从高到低 时间由新及旧
		var el *list.Element
		if stm.List.SellList.Len() == 0 { //list为空
			stm.List.SellList.PushBack(o)
			stm.codesMap[o.Oid] = stm.List.BuyList.Front()
		} else {
			e := stm.List.SellList.Front()
			for {
				if e.Value.(*CommonOrder).Price < o.Price {
					stm.List.SellList.InsertBefore(o, e)
					stm.codesMap[o.Oid] = e
					break
				} else if e.Value.(*CommonOrder).Price == o.Price {
					if e.Value.(*CommonOrder).OrderTime < o.OrderTime { //相同价格比较时间
						stm.List.SellList.InsertBefore(o, e)
						stm.codesMap[o.Oid] = e
						break
					}
				}
				el = e.Next()
				if el == nil { //结束
					stm.List.SellList.InsertAfter(o, e)
					stm.codesMap[o.Oid] = e
					break
				}
				e = el
			}
		}
	}
	l4g.Info("当前队列:%s %s:%v", o.Code, "插入耗时", time.Now().Sub(begin))
	l4g.Info("当前买队列:%s %s:%d", o.Code, "大小", stm.List.BuyList.Len())
	l4g.Info("当前卖队列:%s %s:%d", o.Code, "大小", stm.List.SellList.Len())
}

func (stm *StkMatchList) DelOrder(i int32) {
	stm.l.Lock()
	defer stm.l.Unlock()
	if i == BUY { //买
		if stm.List.BuyList.Len() == 0 {
			//			l4g.Info("BuyList is blank-------")
		} else {
			b := stm.List.BuyList.Back()
			stm.List.BuyList.Remove(b)
		}

	} else if i == SELL { //卖
		if stm.List.SellList.Len() == 0 {
			//			l4g.Info("SellList is blank-------")
		} else {
			b := stm.List.SellList.Back()
			stm.List.SellList.Remove(b)
		}
	} else {
		l4g.Info("buy/sell mark is wrong-------")
	}
}

func (stm *StkMatchList) GetOrder(i int32) (c *CommonOrder) {
	if i == BUY { //买
		if stm.List.BuyList.Len() == 0 {
			//			l4g.Info("BuyList is blank-------")
		} else {
			return stm.List.BuyList.Back().Value.(*CommonOrder)
		}

	} else if i == SELL { //卖
		if stm.List.SellList.Len() == 0 {
			//			l4g.Info("SellList is blank-------")
		} else {
			return stm.List.SellList.Back().Value.(*CommonOrder)
		}
	} else {
		l4g.Info("buy/sell mark is wrong-------")
	}
	return nil
}

func (stm *StkMatchList) GetOrderByOid(oid int64) (c *list.Element) {
	//	for e := stm.List.BuyList.Front(); e != nil; e = e.Next() {
	//		if e.Value.(*CommonOrder).Oid == oid {
	//			return e
	//		}
	//	}
	//	for e := stm.List.SellList.Front(); e != nil; e = e.Next() {
	//		if e.Value.(*CommonOrder).Oid == oid {
	//			return e
	//		}
	//	}
	//	return nil
	return stm.codesMap[oid]
}

func (stm *StkMatchList) RemoveOrderByObj(c *CommonOrder) bool {
	stm.l.Lock()
	defer stm.l.Unlock()
	for {
		//根据订单号查询出待删除的订单
		delObj := stm.GetOrderByOid(c.Oid)
		if nil == delObj {
			return true
		}
		if c.BsFlag == BUY { //买
			if stm.List.BuyList.Len() == 0 {
				break
			} else {
				o := stm.List.BuyList.Remove(delObj)
				delete(stm.codesMap, c.Oid)
				l4g.Info("del:%v %s:%d", o, "当前队列大小", stm.List.BuyList.Len())
				return true
			}

		} else if c.BsFlag == SELL { //卖
			if stm.List.SellList.Len() == 0 {
				break
			} else {
				l4g.Info(delObj)
				o := stm.List.SellList.Remove(delObj)
				delete(stm.codesMap, c.Oid)
				l4g.Info("del:%v %s:%d", o, "当前队列大小", stm.List.SellList.Len())
				return true
			}
		} else {
			l4g.Info("buy/sell mark is wrong-------")
			break
		}
	}
	return true
}

func (stm *StkMatchList) RemoveAll() {
	for e := stm.List.BuyList.Front(); e != nil; e = e.Next() {
		stm.List.BuyList.Remove(e)

	}
	for e := stm.List.SellList.Front(); e != nil; e = e.Next() {
		stm.List.SellList.Remove(e)
	}
	for k,_ := range stm.codesMap {
		delete(stm.codesMap,k)
	}
	
	l4g.Info("buy/sell queue "+stm.Identify()+" clean empty-------")
}

func (self *StkMatchList) Identify() string {
	return self.code
}
