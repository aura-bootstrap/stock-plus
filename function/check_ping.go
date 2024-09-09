package function

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

type CheckPingFunction struct {
	output bytes.Buffer
}

func (f *CheckPingFunction) String() string {
	return "CheckPingFunction"
}

func (f *CheckPingFunction) Main(input <-chan string, output chan<- string) {
	f.output.Reset()

	defer func() {
		output <- f.output.String()
	}()

	host := "61.144.239.92"
	cmdOutput, err := f.ping(host)
	if err != nil {
		f.Output(fmt.Sprintf("执行 ping 命令失败: %v\n", err))
		return
	}

	// 打印 ping 命令的输出
	f.Output(fmt.Sprintln("ping 命令输出:"))
	f.Output(fmt.Sprintln(cmdOutput))

	minRtt, maxRtt, avgRtt, err := f.parsePingOutput(cmdOutput)
	if err != nil {
		f.Output(fmt.Sprintf("解析 ping 输出失败: %v\n", err))
		return
	}

	f.Output(fmt.Sprintf("最小单向延迟: %.2f ms\n", f.calculateOneWayDelay(minRtt)))
	f.Output(fmt.Sprintf("最大单向延迟: %.2f ms\n", f.calculateOneWayDelay(maxRtt)))
	f.Output(fmt.Sprintf("平均单向延迟: %.2f ms\n", f.calculateOneWayDelay(avgRtt)))
}

func (f *CheckPingFunction) Output(s string) {
	log.Printf(s)
	f.output.WriteString(s)
}

// 执行 ping 命令
func (*CheckPingFunction) ping(host string) (string, error) {
	cmd := exec.Command("ping", host)
	output, _ := cmd.Output()

	// 将 GBK 编码转换为 UTF-8 编码
	utf8Output, err := gbkToUtf8(output)
	if err != nil {
		return "", err
	}

	return string(utf8Output), nil
}

// 将 GBK 编码转换为 UTF-8 编码
func (*CheckPingFunction) gbkToUtf8(data []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(data), simplifiedchinese.GBK.NewDecoder())
	return io.ReadAll(reader)
}

// 解析 ping 输出
func (*CheckPingFunction) parsePingOutput(output string) (minRtt, maxRtt, avgRtt int, err error) {
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
func (*CheckPingFunction) calculateOneWayDelay(rtt int) float64 {
	return float64(rtt) / 2
}
