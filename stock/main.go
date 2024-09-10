package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type NightMarketCommissionFunction struct {
	input  <-chan string
	output chan<- string
}

type Stock struct {
	ID     string
	Price  float64
	Amount int
}

func (f *NightMarketCommissionFunction) Main(input <-chan string, output chan<- string) {
	f.input = input
	f.output = output

	targetTime := f.getTime()
	stocks := f.getStock()

	f.waitUntilTargetTime(targetTime, stock)

	for _, stock := range stocks {
		processStock(stock)
	}

	fmt.Println("按回车键退出...")
	bufio.NewReader(os.Stdin).ReadString('\n')
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

func (f *NightMarketCommissionFunction) waitUntilTargetTime(targetTime time.Time, stocks Stock) {
	logWithTimestamp(fmt.Sprintf("等待时间到达: %v", targetTime))
	timeUntilTarget := time.Until(targetTime)

	// 计算目标时间前10秒的时间点
	timeBefore10Seconds := targetTime.Add(-10 * time.Second)
	timeUntil10SecondsBefore := time.Until(timeBefore10Seconds)

	// 等待直到目标时间前10秒
	timer := time.NewTimer(timeUntil10SecondsBefore)
	<-timer.C

	// 发送连通性测试请求
	for _, stock := range stocks {
		url := fmt.Sprintf("http://localhost:8888/buy?stock_code=%s&price=%.2f&amount=%d", stock.ID, stock.Price, stock.Amount)
		logWithTimestamp(fmt.Sprintf("发送连通性测试请求: %s", url))
		resp, err := http.Get(url)
		if err != nil {
			logWithTimestamp(fmt.Sprintf("连通性测试请求失败: %v", err))
		} else {
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				logWithTimestamp(fmt.Sprintf("读取连通性测试响应失败: %v", err))
			} else {
				logWithTimestamp(fmt.Sprintf("连通性测试响应内容: %s", string(body)))
			}
		}
	}

	// 等待剩余时间直到目标时间
	timeUntilTarget = time.Until(targetTime)
	timer = time.NewTimer(timeUntilTarget)
	<-timer.C
}

func processStock(stock Stock) {
	url := fmt.Sprintf("http://localhost:8888/buy?stock_code=%s&price=%.2f&amount=%d", stock.ID, stock.Price, stock.Amount)
	maxRetries := 30
	for i := 0; i < maxRetries; i++ {
		logWithTimestamp(fmt.Sprintf("发送请求: %s", url))
		resp, err := http.Get(url)
		if err != nil {
			logWithTimestamp(fmt.Sprintf("请求失败: %v", err))
			continue
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			logWithTimestamp(fmt.Sprintf("读取响应失败: %v", err))
			continue
		}

		logWithTimestamp(fmt.Sprintf("响应内容: %s", string(body)))

		var result map[string]interface{}
		if err := json.Unmarshal(body, &result); err != nil {
			logWithTimestamp(fmt.Sprintf("解析JSON失败: %v", err))
			continue
		}

		message, ok := result["message"].(string)
		if !ok {
			logWithTimestamp("响应中没有找到message字段")
			continue
		}

		decodedMessage, err := strconv.Unquote(`"` + message + `"`)
		if err != nil {
			logWithTimestamp(fmt.Sprintf("解码message失败: %v", err))
			continue
		}

		if strings.Contains(decodedMessage, "合同") {
			logWithTimestamp(fmt.Sprintf("请求成功，message: %s", decodedMessage))
			break
		} else {
			logWithTimestamp(fmt.Sprintf("请求失败，message: %s", decodedMessage))
		}
	}
}

func logWithTimestamp(message string) {
	fmt.Printf("%s %s\n", time.Now().Format("2006-01-02 15:04:05.000"), message)
}
