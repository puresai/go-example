package util

import (
	"fmt"
	"bytes"
	"strings"
	"time"    
	"unicode/utf8"
	"strconv"
	"crypto/md5"
	"net/http"
	"net/url"
	"io/ioutil"
)


// 获取本周一时间戳
func GetFirstDateOfWeekTS() (ts int64) {
	now := time.Now()
 
	offset := int(time.Monday - now.Weekday())
	if offset > 0 {
		offset = -6
	}
 
	ts = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offset).Unix()
	return
}

// 获取周几时间戳
func GetWeekTS(day int64) (ts int64) {
	thisWeekMonday := GetFirstDateOfWeekTS()
	ts = thisWeekMonday + (day)*86400
	return
}

// 字符串长度
func StrLen(str string) int {
	return utf8.RuneCountInString(str)
}

// 截取字符串
func StrSub(str string, sub ...int) string {
	start := sub[0]
	length := 0
	if len(sub) > 1 {
		length = sub[1]
	}

	if length < 1 {
		return string(([]rune(str))[start:])
	}
    return string(([]rune(str))[start:length])
}

// 合并字符串
func StrCombine(str ...string) string {
    var bt bytes.Buffer
    for _, arg := range str {
        bt.WriteString(arg)
    }
    //获得拼接后的字符串
    return bt.String()
}

func UnicodeIndex(str,substr string) int {
	// 子串在字符串的字节位置
	result := strings.Index(str,substr)  
	if result >= 0 {
	  // 获得子串之前的字符串并转换成[]byte
	  prefix := []byte(str)[0:result]  
	  // 将子串之前的字符串转换成[]rune
	  rs := []rune(string(prefix))  
	  // 获得子串之前的字符串的长度，便是子串在字符串的字符位置
	  result = len(rs)
	}
	
	return result
}

func ToUnicode(str string) string {
	textQuoted := strconv.QuoteToASCII(str)
    return textQuoted[1 : len(textQuoted)-1]
}

func UnicodeTo(str string) string {
    sUnicodev := strings.Split(str, "\\u")
    var context string
    for _, v := range sUnicodev {
        if len(v) < 1 {
            continue
        }
        temp, err := strconv.ParseInt(v, 16, 32)
        if err != nil {
            panic(err)
        }
        context += fmt.Sprintf("%c", temp)
    }
    return context
}

// 获取当天0点时间抽
func TodayTS() int64 {
	now := time.Now()
    return GetZeroTime(now).Unix()
}

func TodayDate()  string {
	return time.Now().Format("2006/01/02")
}

//获取传入的时间所在月份的第一天，即某月第一天的0点。如传入time.Now(), 返回当前月份的第一天0点时间。
func GetFirstDateOfMonth(d time.Time) time.Time {
	d = d.AddDate(0, 0, -d.Day() + 1)
	return GetZeroTime(d)
}
//获取传入的时间所在月份的最后一天，即某月最后一天的0点。如传入time.Now(), 返回当前月份的最后一天0点时间。
func GetLastDateOfMonth(d time.Time) time.Time {
	return GetFirstDateOfMonth(d).AddDate(0, 1, -1)
}
 
//获取某一天的0点时间
func GetZeroTime(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}

func Md5(str string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(str)));
}

func HttpGet(iUrl string, params map[string]string) ([]byte, error) {
	query := ""
    for k, v := range params {
        query = fmt.Sprintf("%s%s=%s&", query, k, url.QueryEscape(v))
	}
	
	client := http.Client{Timeout: 5 * time.Second} //创建客户端
	resp, err := client.Get(fmt.Sprintf("%s?%s", iUrl, query))
	fmt.Println(iUrl,query)

    if err != nil {
        return nil, err
	}
	defer resp.Body.Close()
	res, err := ioutil.ReadAll(resp.Body)
	if err!=nil{
		fmt.Println(err)
		return nil,err
    }
	
	return res, nil
}