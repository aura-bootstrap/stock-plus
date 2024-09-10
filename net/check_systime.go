package net

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"reflect"
	"regexp"

	"github.com/bootstrap-library/stock-plus/function"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

type CheckSystimeFunction struct {
	outputBuffer bytes.Buffer
}

func NewCheckSystimeFunction() function.Function {
	return &CheckSystimeFunction{}
}

func (f *CheckSystimeFunction) String() string {
	return reflect.TypeOf(f).Name()
}

func (f *CheckSystimeFunction) Main(input <-chan string, output chan<- string) {
	f.outputBuffer.Reset()
	defer func() {
		output <- f.outputBuffer.String()
	}()

	offset, err := f.getNTPTimeOffset()
	if err != nil {
		f.printf("获取 NTP 时间误差失败: %v\n", err)
		return
	}
	f.printf("NTP 时间误差: %v\n", offset)
}

func (f *CheckSystimeFunction) getNTPTimeOffset() (string, error) {
	cmd := exec.Command("w32tm", "/stripchart", "/computer:time.windows.com", "/dataonly", "/samples:1")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	// 将 GBK 编码转换为 UTF-8 编码
	utf8Output, err := f.gbkToUtf8(out.Bytes())
	if err != nil {
		return "", err
	}

	output := string(utf8Output)
	f.printf("w32tm 命令输出: %v\n", output)

	// 使用正则表达式匹配最后一行的最后一个逗号分隔符后的内容
	re := regexp.MustCompile(`(?m)^.*,\s*(-?\d+\.\d+s)$`)
	matches := re.FindStringSubmatch(output)
	if len(matches) == 2 {
		return matches[1], nil
	}

	return "", fmt.Errorf("未找到时间误差")
}

func (f *CheckSystimeFunction) gbkToUtf8(data []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(data), simplifiedchinese.GBK.NewDecoder())
	return io.ReadAll(reader)
}

func (f *CheckSystimeFunction) printf(format string, args ...any) {
	f.outputBuffer.WriteString(fmt.Sprintf(format, args...))
}
