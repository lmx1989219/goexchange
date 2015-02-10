package exchange

import (
	"database/sql"
	l4g "github.com/alecthomas/log4go"
	_ "mysql"
)

type ContractsBean struct {
	MarketCode string
	Code       string
}

func LoadData() []string {
	codeSlice := []string{}
	db, err := sql.Open("mysql", "root:601519@tcp(10.15.107.105:3306)/mobile_counter_stock?charset=utf8")
	CheckErr(err)
	defer db.Close()
	rows, err := db.Query("select * from T_CONCTRACTS")
	CheckErr(err)
	i := 0
	for rows.Next() {
		//value number must be same as db table's fields number
		var ID int
		var CODE string
		var NAME string
		var MARKET string
		_ = rows.Scan(&ID, &CODE, &NAME, &MARKET)
		cb := &ContractsBean{MarketCode: MARKET, Code: CODE}
		MarketQuoteMap[CODE] = cb
		codeSlice = append(codeSlice, CODE)
		i++
	}
	l4g.Info("合约加载完成")
	return codeSlice
}

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}
