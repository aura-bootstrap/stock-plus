package stock

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type NightMarketCommissionFunction struct {
	input  <-chan string
	output chan<- string
	tm     time.Time
	stock  Stock
}

type Stock struct {
	ID     string
	Price  float64
	Amount int
}

func (f *NightMarketCommissionFunction) Main(input <-chan string, output chan<- string) {
	f.input = input
	f.output = output
	f.tm = f.getTime()
	f.stock = f.getStock()

	f.waitUntilPrepare()

	f.prepare()

	f.waitUntilProcess()

	f.process()
}

func (f *NightMarketCommissionFunction) getTime() time.Time {
	for {
		f.output <- "请输入目标时间 (格式: 15:04:05.000)："
		text := <-f.input
		tm, err := time.Parse("15:04:05.000", text)
		if err != nil {
			f.output <- "时间格式不正确, 请重新输入"
			continue
		}
		return tm
	}
}

func (f *NightMarketCommissionFunction) getStock() Stock {
	for {
		f.output <- "请输入股票信息 (格式: 股票代码,价格,数量)："
		text := <-f.input
		parts := strings.Split(text, ",")
		if len(parts) != 3 {
			f.output <- "股票信息格式不正确, 请重新输入"
			continue
		}
		price, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			f.output <- "价格格式不正确, 请重新输入"
			continue
		}
		amount, err := strconv.Atoi(parts[2])
		if err != nil {
			f.output <- "数量格式不正确, 请重新输入"
			continue
		}
		return Stock{ID: parts[0], Price: price, Amount: amount}
	}
}

func (f *NightMarketCommissionFunction) waitUntilPrepare() {
	timeBefore10Seconds := f.tm.Add(-10 * time.Second)
	f.printf("等待准备时间到达: %v", timeBefore10Seconds)
	timeUntil10SecondsBefore := time.Until(timeBefore10Seconds)
	timer := time.NewTimer(timeUntil10SecondsBefore)
	<-timer.C
}

func (f *NightMarketCommissionFunction) prepare() {
	url := fmt.Sprintf("http://localhost:8888/buy?stock_code=%s&price=%.2f&amount=%d", f.stock.ID, f.stock.Price, f.stock.Amount)
	f.printf("发送连通性测试请求: %s", url)
	resp, err := http.Get(url)
	if err != nil {
		f.printf("连通性测试请求失败: %v", err)
	} else {
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			f.printf("读取连通性测试响应失败: %v", err)
		} else {
			f.printf("连通性测试响应内容: %s", string(body))
		}
	}
}

func (f *NightMarketCommissionFunction) waitUntilProcess() {
	f.printf("等待处理时间到达: %v", f.tm)
	timeUntilTarget := time.Until(f.tm)
	timer := time.NewTimer(timeUntilTarget)
	<-timer.C
}

func (f *NightMarketCommissionFunction) process() {
	url := fmt.Sprintf("http://localhost:8888/buy?stock_code=%s&price=%.2f&amount=%d", f.stock.ID, f.stock.Price, f.stock.Amount)
	maxRetries := 30
	for i := 0; i < maxRetries; i++ {
		f.printf("发送请求: %s", url)
		resp, err := http.Get(url)
		if err != nil {
			f.printf("请求失败: %v", err)
			continue
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			f.printf("读取响应失败: %v", err)
			continue
		}

		f.printf("响应内容: %s", string(body))

		var result map[string]interface{}
		if err := json.Unmarshal(body, &result); err != nil {
			f.printf("解析JSON失败: %v", err)
			continue
		}

		message, ok := result["message"].(string)
		if !ok {
			f.printf("响应中没有找到message字段")
			continue
		}

		decodedMessage, err := strconv.Unquote(`"` + message + `"`)
		if err != nil {
			f.printf("解码message失败: %v", err)
			continue
		}

		if strings.Contains(decodedMessage, "合同") {
			f.printf("请求成功，message: %s", decodedMessage)
			break
		} else {
			f.printf("请求失败，message: %s", decodedMessage)
		}
	}
}

func (f *NightMarketCommissionFunction) printf(format string, args ...any) {
	text := fmt.Sprintf("%s %s\n", time.Now().Format("2006-01-02 15:04:05.000"), fmt.Sprintf(format, args...))
	f.output <- text
}
