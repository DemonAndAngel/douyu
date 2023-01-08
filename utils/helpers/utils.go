package helpers

import (
	"encoding/json"
	"github.com/dop251/goja"
	"github.com/go-co-op/gocron"
	"io/ioutil"
	"net/url"
	"strings"
	"time"
)

var Sign1 func(roomId int64, did string, tt int64) string
var Sign2 func(roomId int64, did string, tt int64) string

func ParseSignStr(signStr string) (string, string) {
	uu, _ := url.ParseQuery(signStr)
	sign := ""
	v := ""
	if signA, ok := uu["sign"]; ok {
		if len(signA) > 0 {
			sign = signA[0]
		}
	}
	if vA, ok := uu["v"]; ok {
		if len(vA) > 0 {
			v = vA[0]
		}
	}
	return sign, v
}

func init() {
	vm := goja.New() // 创建engine实例
	b, _ := ioutil.ReadFile("./js/crypto-js.min.js")
	_, err := vm.RunScript("crypto-js.min.js", string(b))
	if err != nil {
		panic("js代码异常;信息:" + err.Error())
	}
	b, _ = ioutil.ReadFile("./js/sign1.js")
	_, err = vm.RunString(string(b))
	if err != nil {
		panic("js代码异常;信息:" + err.Error())
	}
	b, _ = ioutil.ReadFile("./js/sign2.js")
	_, err = vm.RunString(string(b))
	if err != nil {
		panic("js代码异常;信息:" + err.Error())
	}
	err = vm.ExportTo(vm.Get("sign1"), &Sign1) // 将执行的结果转换为Golang对应的类型
	if err != nil {
		panic("映射Go函数失败;信息:" + err.Error())
	}
	err = vm.ExportTo(vm.Get("sign2"), &Sign2) // 将执行的结果转换为Golang对应的类型
	if err != nil {
		panic("映射Go函数失败;信息:" + err.Error())
	}
}

func JsonMarshal(m interface{}) string {
	byteData, _ := json.Marshal(m)
	str := string(byteData)
	if str == "null" {
		str = ""
	}
	return str
}

// CreateTimeFormat 日期转时间
func CreateTimeFormat(format string, now string) time.Time {
	if format == time.RFC3339 {
		t, _ := time.ParseInLocation(time.RFC3339, now, time.Local)
		return t
	} else {
		format = strings.Replace(format, "Y", "2006", -1)
		format = strings.Replace(format, "m", "01", -1)
		format = strings.Replace(format, "d", "02", -1)
		format = strings.Replace(format, "H", "15", -1)
		format = strings.Replace(format, "i", "04", -1)
		format = strings.Replace(format, "s", "05", -1)
		t, _ := time.ParseInLocation(format, now, time.Local)
		return t
	}
}

// TimeFormat 时间转日期
func TimeFormat(format string, t time.Time) string {
	if format != time.RFC3339 {
		format = strings.Replace(format, "Y", "2006", -1)
		format = strings.Replace(format, "m", "01", -1)
		format = strings.Replace(format, "d", "02", -1)
		format = strings.Replace(format, "H", "15", -1)
		format = strings.Replace(format, "i", "04", -1)
		format = strings.Replace(format, "s", "05", -1)
	}
	if t.IsZero() {
		return ""
	}
	return t.Format(format)
}

// GetScheduler 获取调度器
func GetScheduler(unique bool) *gocron.Scheduler {
	s := gocron.NewScheduler(time.Local)
	if unique {
		s.TagsUnique()
	}
	s.StartAsync()
	return s
}
