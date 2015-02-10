package main

import (
	"code.google.com/p/goprotobuf/proto"
	"fmt"
	//ex "gw.com.cn/dzhyun/dzhexchange"
	"gw.com.cn/dzhyun/dzhyun.git"
	"math/rand"
	"runtime"
	"strconv"
	"time"
)

import . "gw.com.cn/dzhyun/sdk.bus.git"
import "log"

//系统总线 推送 10.15.107.246:10201  10.15.107.246:10202
//应用总线 请求应答 10.15.107.246:10100

func main() {
	TestRouter()
}

//ua连的app总线端口
const address string = "tcp://10.15.107.246:10100"

type NullOutput struct {
}

func (this *NullOutput) Write(p []byte) (int, error) {
	return 0, nil
}

var rd = rand.New(rand.NewSource(time.Now().UnixNano() / (1000 * 1000)))
var codesarr = rand.New(rand.NewSource(time.Now().Unix()))

//单台ua-app，1w并发5s左右，每隔1s发送一次 发送5次， 平均每次压力结束耗时8s,压力不是线性增长
func TestRouter() {
	runtime.GOMAXPROCS(8)
	codes := LoadData()
	quit := make(chan int, 1000)
	const path string = "/root/app/dzhexchange"

	log.Println("ua go ...")
	count := 0
	tc := 10
	startTime := time.Now()

	bus := NewBusHelper()
	defer func() {
		log.Println("ua is exiting ...")
		bus.Shutdown()
	}()
	if bus.AddRouter(address) == nil {
		for i := 0; i < tc; i++ {
			go func() {
				oid := rd.Int63() / 1000 / 10000
				wtprice := rd.Float64()*3 + 2
				fmt.Println(oid)
				request := bus.CreateRequest()
				ranIdx := codesarr.Intn(2000)
				//发送json格式数据包
				jsonObj := `{"ProtoNO": "9990",
				"OrderId":` + fmt.Sprintf("%d", oid) + `,"ProductCode":"` + codes[ranIdx] + `","MarketCode":"SZ",
				"CustomerNo":"test","WtPrice":` + strconv.FormatFloat(float64(wtprice), 'f', -1, 32) + `,"WtCount":100,"WtType":1,"BsFlag":1,"OcFlag":1,"GameId":"FAFP",
				"WtTime":"`+time.Now().String()+`"}`

				response, e := request.Send(path+"?"+jsonObj, time.Second*60)
				if e != nil {
					fmt.Println(e.Error())
				} else {
					obj := &dzhyun.BusResponse{}
					proto.Unmarshal(response.Data, obj)
					fmt.Println("撮合应答：", obj.GetResult())
				}
				quit <- 1

			}()
		}
	} else {
		fmt.Println("ua: AddRouter failed")
	}

	for {
		<-quit
		count++
		if count == tc {
			break
		}
	}
	endTime := time.Now()
	log.Println("**** test", tc, "requests:", endTime.Sub(startTime), "****")
	close(quit)
}
