package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"reflect"

	"github.com/anaskhan96/soup"
)

func ReplaceSpace(s string) string {
	var result []rune
	const space = ' '
	for _, r := range s {
		if r == space {
			result = append(result, '\u3000')
			continue
		}
		result = append(result, r)
	}
	return string(result)
}

func GetStockKey(stockNum int) (stockKey string) {
	type MsgArray struct {
		Ex  string `json:"ex"`
		D   string `json:"d"`
		It  string `json:"it"`
		N   string `json:"n"`
		I   string `json:"i"`
		Ip  string `json:"ip"`
		W   string `json:"w"`
		U   string `json:"u"`
		T   string `json:"t"`
		P   string `json:"p"`
		Ch  string `json:"ch"`
		Key string `json:"key"`
		Y   string `json:"y"`
	}

	type QueryTime struct {
		StockDetail    int `json:"stockDetail"`
		TotalMicroTime int `json:"totalMicroTime"`
	}

	type Stock struct {
		MsgArrays []*MsgArray `json:"msgArray"`
		Rtmessage string      `json:"rtmessage"`
		Rtcode    string      `json:"rtcode"`
		QueryTime `json:"queryTime"`
	}

	url := fmt.Sprintf("http://mis.tse.com.tw/stock/api/getStock.jsp?ch=%d.tw", stockNum)
	resp, err := soup.Get(url)
	if err != nil {
		fmt.Println("http transport error is:", err)
		stockKey = ""
		return
	}
	stock := Stock{}
	json.Unmarshal([]byte(resp), &stock)
	stockKey = stock.MsgArrays[0].Key

	return
}

func GetStockInfo(stockKey string) {
	type MsgArray struct {
		C  string `json:"c"`
		D  string `json:"d"`
		Nf string `json:"nf"`
		O  string `json:"o"`
		H  string `json:"h"`
		L  string `json:"l"`
	}

	type StockInfo struct {
		MsgArrays []*MsgArray `json:"msgArray"`
	}

	getField := func(obj *MsgArray, field string) string {
		r := reflect.ValueOf(obj)
		f := reflect.Indirect(r).FieldByName(field)
		return string(f.String())
	}

	dict := map[string]string{
		"股票號碼": "C",
		"日期":   "D",
		"公司名稱": "Nf",
		"開盤":   "O",
		"最高":   "H",
		"最低":   "L"}

	url := fmt.Sprintf("http://mis.tse.com.tw/stock/api/getStockInfo.jsp?ex_ch=%s", stockKey)
	resp, err := soup.Get(url)
	if err != nil {
		fmt.Println("http transport error is:", err)
		return
	}
	stockInfo := StockInfo{}
	json.Unmarshal([]byte(resp), &stockInfo)
	for key, value := range dict {
		s := fmt.Sprintf("%-6s: %s", key, getField(stockInfo.MsgArrays[0], value))
		fmt.Println(ReplaceSpace(s))
	}
}

func main() {
	number := flag.Int("number", 0, "stock number")
	flag.Parse()

	stockKey := GetStockKey(*number)
	GetStockInfo(stockKey)
}
