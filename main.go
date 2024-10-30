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

type Stock struct {
	ID     string
	Price  float64
	Amount int
}

func main() {
	targetTime := getTargetTimeFromUser()
	stock := getStockFromUser()
	maxRetries := getMaxRetriesFromUser()

	waitUntilTargetTime(targetTime, stock)

	processStock(stock, maxRetries)

	fmt.Println("按回车键退出...")
	bufio.NewReader(os.Stdin).ReadString('\n')
}

func getMaxRetriesFromUser() int {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("请输入最大重试次数：")
		retriesStr, _ := reader.ReadString('\n')
		retries, err := strconv.Atoi(strings.TrimSpace(retriesStr))
		if err != nil {
			fmt.Println("重试次数格式不正确，请重新输入。")
			continue
		}
		return retries
	}
}

func getStockFromUser() Stock {
	for {
		reader := bufio.NewReader(os.Stdin)

		fmt.Println("请输入股票代码：")
		stockID, _ := reader.ReadString('\n')
		stockID = strings.TrimSpace(stockID)

		fmt.Println("请输入股票价格：")
		priceStr, _ := reader.ReadString('\n')
		price, _ := strconv.ParseFloat(strings.TrimSpace(priceStr), 64)

		fmt.Println("请输入股票数量：")
		quantityStr, _ := reader.ReadString('\n')
		quantity, _ := strconv.Atoi(strings.TrimSpace(quantityStr))

		fmt.Println("您输入的股票信息如下：")
		fmt.Printf("股票代码：%s\n", stockID)
		fmt.Printf("股票价格：%.2f\n", price)
		fmt.Printf("股票数量：%d\n", quantity)
		fmt.Println("是否确认？(y/n)：")
		confirmation, _ := reader.ReadString('\n')
		confirmation = strings.TrimSpace(confirmation)
		if strings.ToLower(confirmation) != "y" {
			fmt.Println("请重新输入股票信息。")
			continue
		}

		return Stock{ID: stockID, Price: price, Amount: quantity}
	}
}

func getTargetTimeFromUser() time.Time {
	reader := bufio.NewReader(os.Stdin)
	var year, month, day, hour, minute, second, microsecond int
	var err error

	for {
		fmt.Println("请输入年份（格式：YYYY）：")
		yearStr, _ := reader.ReadString('\n')
		year, err = strconv.Atoi(strings.TrimSpace(yearStr))
		if err != nil {
			fmt.Println("年份格式不正确，请重新输入。")
			continue
		}

		fmt.Println("请输入月份（格式：MM）：")
		monthStr, _ := reader.ReadString('\n')
		month, err = strconv.Atoi(strings.TrimSpace(monthStr))
		if err != nil {
			fmt.Println("月份格式不正确，请重新输入。")
			continue
		}

		fmt.Println("请输入日期（格式：DD）：")
		dayStr, _ := reader.ReadString('\n')
		day, err = strconv.Atoi(strings.TrimSpace(dayStr))
		if err != nil {
			fmt.Println("日期格式不正确，请重新输入。")
			continue
		}

		fmt.Println("请输入小时（格式：HH）：")
		hourStr, _ := reader.ReadString('\n')
		hour, err = strconv.Atoi(strings.TrimSpace(hourStr))
		if err != nil {
			fmt.Println("小时格式不正确，请重新输入。")
			continue
		}

		fmt.Println("请输入分钟（格式：MM）：")
		minuteStr, _ := reader.ReadString('\n')
		minute, err = strconv.Atoi(strings.TrimSpace(minuteStr))
		if err != nil {
			fmt.Println("分钟格式不正确，请重新输入。")
			continue
		}

		fmt.Println("请输入秒（格式：SS）：")
		secondStr, _ := reader.ReadString('\n')
		second, err = strconv.Atoi(strings.TrimSpace(secondStr))
		if err != nil {
			fmt.Println("秒格式不正确，请重新输入。")
			continue
		}

		fmt.Println("请输入微秒（格式：UUUUUU）：")
		microsecondStr, _ := reader.ReadString('\n')
		microsecond, err = strconv.Atoi(strings.TrimSpace(microsecondStr))
		if err != nil {
			fmt.Println("微秒格式不正确，请重新输入。")
			continue
		}

		targetTime := time.Date(year, time.Month(month), day, hour, minute, second, microsecond*1e3, time.Local)

		fmt.Printf("您输入的目标时间是：%s，是否确认？(y/n)：", targetTime.Format("2006-01-02 15:04:05.000000"))
		confirmation, _ := reader.ReadString('\n')
		confirmation = strings.TrimSpace(confirmation)
		if strings.ToLower(confirmation) == "y" {
			return targetTime
		} else {
			fmt.Println("请重新输入目标时间。")
		}
	}
}

func waitUntilTargetTime(targetTime time.Time, stock Stock) {
	logWithTimestamp(fmt.Sprintf("等待时间到达: %v", targetTime))
	timeUntilTarget := time.Until(targetTime)

	// 计算目标时间前10秒的时间点
	timeBefore10Seconds := targetTime.Add(-10 * time.Second)
	timeUntil10SecondsBefore := time.Until(timeBefore10Seconds)

	// 等待直到目标时间前10秒
	timer := time.NewTimer(timeUntil10SecondsBefore)
	<-timer.C

	// 发送连通性测试请求
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

	// 等待剩余时间直到目标时间
	timeUntilTarget = time.Until(targetTime)
	timer = time.NewTimer(timeUntilTarget)
	<-timer.C
}

func processStock(stock Stock, maxRetries int) {
	url := fmt.Sprintf("http://localhost:8888/buy?stock_code=%s&price=%.2f&amount=%d", stock.ID, stock.Price, stock.Amount)
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
	fmt.Printf("%s %s\n", time.Now().Format("2006-01-02 15:04:05.000000"), message)
}
