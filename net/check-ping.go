package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// 执行 ping 命令
func ping(host string) (string, error) {
	cmd := exec.Command("ping", host)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	// 将 GBK 编码转换为 UTF-8 编码
	utf8Output, err := gbkToUtf8(output)
	if err != nil {
		return "", err
	}

	return string(utf8Output), nil
}

// 将 GBK 编码转换为 UTF-8 编码
func gbkToUtf8(data []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(data), simplifiedchinese.GBK.NewDecoder())
	return ioutil.ReadAll(reader)
}

// 解析 ping 输出
func parsePingOutput(output string) (minRtt, maxRtt, avgRtt int, err error) {
	var rttPattern *regexp.Regexp

	// 检测操作系统语言
	lang := os.Getenv("LANG")
	if strings.Contains(lang, "zh") {
		// 中文环境
		rttPattern = regexp.MustCompile(`最短 = (\d+)ms，最长 = (\d+)ms，平均 = (\d+)ms`)
	} else {
		// 英文环境
		rttPattern = regexp.MustCompile(`Minimum = (\d+)ms, Maximum = (\d+)ms, Average = (\d+)ms`)
	}

	match := rttPattern.FindStringSubmatch(output)
	if match != nil {
		minRtt, err = strconv.Atoi(match[1])
		if err != nil {
			return
		}
		maxRtt, err = strconv.Atoi(match[2])
		if err != nil {
			return
		}
		avgRtt, err = strconv.Atoi(match[3])
		if err != nil {
			return
		}
		return minRtt, maxRtt, avgRtt, nil
	}
	return 0, 0, 0, fmt.Errorf("无法解析 ping 输出")
}

// 计算单向延迟
func calculateOneWayDelay(rtt int) float64 {
	return float64(rtt) / 2
}

// func main() {
// 	host := "www.cicc.com"
// 	output, err := ping(host)
// 	if err != nil {
// 		fmt.Printf("执行 ping 命令失败: %v\n", err)
// 		return
// 	}

// 	// 打印 ping 命令的输出
// 	fmt.Println("ping 命令输出:")
// 	fmt.Println(output)

// 	minRtt, maxRtt, avgRtt, err := parsePingOutput(output)
// 	if err != nil {
// 		fmt.Printf("解析 ping 输出失败: %v\n", err)
// 		return
// 	}

// 	fmt.Printf("最小单向延迟: %.2f ms\n", calculateOneWayDelay(minRtt))
// 	fmt.Printf("最大单向延迟: %.2f ms\n", calculateOneWayDelay(maxRtt))
// 	fmt.Printf("平均单向延迟: %.2f ms\n", calculateOneWayDelay(avgRtt))
// }
