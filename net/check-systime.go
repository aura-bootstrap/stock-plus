package main

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"regexp"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// 将 GBK 编码转换为 UTF-8 编码
func gbkToUtf8(data []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(data), simplifiedchinese.GBK.NewDecoder())
	return io.ReadAll(reader)
}

// 获取 NTP 时间误差
func getNTPTimeOffset() (string, error) {
	cmd := exec.Command("w32tm", "/stripchart", "/computer:time.windows.com", "/dataonly", "/samples:1")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	// 将 GBK 编码转换为 UTF-8 编码
	utf8Output, err := gbkToUtf8(out.Bytes())
	if err != nil {
		return "", err
	}

	output := string(utf8Output)
	fmt.Println("w32tm 命令输出:")
	fmt.Println(output)

	// 使用正则表达式匹配最后一行的最后一个逗号分隔符后的内容
	re := regexp.MustCompile(`(?m)^.*,\s*(-?\d+\.\d+s)$`)
	matches := re.FindStringSubmatch(output)
	if len(matches) == 2 {
		return matches[1], nil
	}

	return "", fmt.Errorf("未找到时间误差")
}

// func main() {
// 	// 获取 NTP 时间误差
// 	offset, err := getNTPTimeOffset()
// 	if err != nil {
// 		fmt.Printf("获取 NTP 时间误差失败: %v\n", err)
// 		return
// 	}
// 	fmt.Printf("NTP 时间误差: %v\n", offset)
// }
