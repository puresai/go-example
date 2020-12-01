package util

import (
	"fmt"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"bytes"
	// "unsafe"
	"strings"
	"time"
	"strconv"
	"errors"
	"net/url"

	"github.com/spf13/viper"
	redis "github.com/go-redis/redis/v8"
	"github.com/robfig/cron"

	"local.com/sai0556/demo3-dingding/db"
)

const (
	KeyTime = "time"
	KeyDingDingID = "ddid"
	KeyWords = "keywords"
	KeyCrontab = "crontab"
	KeyUserCron = "dd_cron:"
)

type RobotApiRes struct {
	Code int `json:"result"`
	Content string `json:"content"`
}

func RobotApi(content string) (string, error) {
	client := http.Client{Timeout: 5 * time.Second} //创建客户端
	resp, err := client.Get(fmt.Sprintf(viper.GetString("robot_api"), url.QueryEscape(content)))

    if err != nil {
        fmt.Printf("FilterWords.Get%v", err)
        return "", err
	}
	defer resp.Body.Close()
	res, err := ioutil.ReadAll(resp.Body)
	if err!=nil{
		fmt.Println(err)
		return "",err
    }
	var s RobotApiRes
    if err := json.Unmarshal([]byte(res), &s); err!=nil{
        fmt.Println(err)
		return "", err
    }

	if s.Code != 0 {
		return content, errors.New("接口错误")
	}
    
	return s.Content, nil

}

func SendDD(msg string) {
	// 打印出来看看是个啥
	fmt.Println("dingding-----------")
	fmt.Println(msg)
	tips := make(map[string]interface{})
	content := make(map[string]interface{})
	tips["msgtype"] = "markdown"
    // @ 是用来提醒群里对应的人
	arr := strings.Split(msg, "@")
    // [提醒]是机器人关键字，个人建议设置机器人限制ip或使用token，比较靠谱
	content["text"] = fmt.Sprintf("%s", strings.Replace(arr[0], "{br}", " \n\n", -1))
	content["title"] = "鹅鹅鹅"

	if len(arr) > 1 {
		mobile := make([]string, 0)
		at := make(map[string]interface{})
		mobile = append(mobile, arr[1])
		at["atMobiles"] = mobile
		tips["at"] = at
		content["text"] = fmt.Sprintf("%s @%s", content["text"], arr[1])
	}

	tips["markdown"] = content

    bytesData, err := json.Marshal(tips)
    if err != nil {
        fmt.Println(err.Error() )
        return
    }
    reader := bytes.NewReader(bytesData)
    url := viper.GetString("dingding_url")
    request, err := http.NewRequest("POST", url, reader)
    if err != nil {
        return
    }
    request.Header.Set("Content-Type", "application/json;charset=UTF-8")
    client := http.Client{}
    _, err = client.Do(request)
    if err != nil {
        fmt.Println(err.Error())
        return
	}
	// 偷懒不重试了
    // respBytes, err := ioutil.ReadAll(resp.Body)
    // if err != nil {
    //     fmt.Println(err.Error())
    //     return
    // }
    // //byte数组直接转成string，优化内存
    // str := (*string)(unsafe.Pointer(&respBytes))
    // fmt.Println(*str)
}

func UpdateKeywords() {
	redis := db.RedisClient.Pipeline()
	key := KeyWords
	// val里签名是类型，后面可以理解成单位
	redis.HSet(db.Ctx, key, "分钟后", "1|60")
	redis.HSet(db.Ctx, key, "时后", "1|3600")
	redis.HSet(db.Ctx, key, "天后", "1|86400")
	redis.HSet(db.Ctx, key, "每天", "-1|1")
	// redis.HSet(db.Ctx, key, "每月", "-2")
	// redis.HSet(db.Ctx, key, "每个月", "4|2")
	redis.HSet(db.Ctx, key, "每周一", "2|0")
	redis.HSet(db.Ctx, key, "每周二", "2|1")
	redis.HSet(db.Ctx, key, "每周三", "2|2")
	redis.HSet(db.Ctx, key, "每周四", "2|3")
	redis.HSet(db.Ctx, key, "每周五", "2|4")
	redis.HSet(db.Ctx, key, "每周六", "2|5")
	redis.HSet(db.Ctx, key, "每周日", "2|6")
	redis.HSet(db.Ctx, key, "周一", "3|0")
	redis.HSet(db.Ctx, key, "周二", "3|1")
	redis.HSet(db.Ctx, key, "周三", "3|2")
	redis.HSet(db.Ctx, key, "周四", "3|3")
	redis.HSet(db.Ctx, key, "周五", "3|4")
	redis.HSet(db.Ctx, key, "周六", "3|5")
	redis.HSet(db.Ctx, key, "周日", "3|6")
	redis.HSet(db.Ctx, key, "下周一", "3|7")
	redis.HSet(db.Ctx, key, "下周二", "3|8")
	redis.HSet(db.Ctx, key, "下周三", "3|9")
	redis.HSet(db.Ctx, key, "下周四", "3|10")
	redis.HSet(db.Ctx, key, "下周五", "3|11")
	redis.HSet(db.Ctx, key, "下周六", "3|12")
	redis.HSet(db.Ctx, key, "下周日", "3|13")
	redis.HSet(db.Ctx, key, "下星期一", "3|7")
	redis.HSet(db.Ctx, key, "下星期二", "3|8")
	redis.HSet(db.Ctx, key, "下星期三", "3|9")
	redis.HSet(db.Ctx, key, "下星期四", "3|10")
	redis.HSet(db.Ctx, key, "下星期五", "3|11")
	redis.HSet(db.Ctx, key, "下星期六", "3|12")
	redis.HSet(db.Ctx, key, "下星期日", "3|13")
	redis.HSet(db.Ctx, key, "今天", "4|0")
	redis.HSet(db.Ctx, key, "明天", "4|1")
	redis.HSet(db.Ctx, key, "后天", "4|2")
	redis.HSet(db.Ctx, key, "取消", "0|0")
	redis.Exec(db.Ctx)
}

func Cron() {
	c := cron.New()
    spec := "*/10 * * * * ?"
	c.AddJob(spec, Queue{})
    c.Start()
}

type Queue struct {
}

func (q Queue) Run() {
	now := time.Now().Unix()
	rd := db.RedisClient
	op := &redis.ZRangeBy{
        Min: "0",
        Max: strconv.FormatInt(now, 10),
    }
    ret, err := rd.ZRangeByScoreWithScores(db.Ctx, KeyCrontab, op).Result()
    if err != nil {
        fmt.Printf("zrangebyscore failed, err:%v\n", err)
        return
	}
    for _, z := range ret {
		fmt.Println(z.Member.(string), z.Score)
		QueueDo(z.Member.(string), z.Score)
    }
}

func QueueDo(msg string, score float64) {
	msgType := msg[0:1]
	SendDD(msg[1:])
	rd := db.RedisClient
	rd.ZRem(db.Ctx, KeyCrontab, msg)

	switch msgType {
		case "2":
			rd.ZAdd(db.Ctx, KeyCrontab, &redis.Z{
				Score: score + 7*86400,
				Member: msg,
			})
		case "3":
			rd.ZAdd(db.Ctx, KeyCrontab, &redis.Z{
				Score: score + 86400,
				Member: msg,
			})
		default:
			rd.ZRem(db.Ctx, KeyCrontab, msg)
	}
}

// 取消提醒
func CancelQueue(uniqueKey string, SenderId string) (err error) {
	rd := db.RedisClient
	member := rd.HGet(db.Ctx, StrCombine(KeyUserCron, SenderId), uniqueKey).Val()
	if member == "" {
		fmt.Println(StrCombine(KeyUserCron, SenderId), uniqueKey)
		err = errors.New("没有此任务")
		return
	}
	fmt.Println(member, "member")
	rd.ZRem(db.Ctx, KeyCrontab, member)
	rd.HDel(db.Ctx, StrCombine(KeyUserCron, SenderId), uniqueKey)
	err = errors.New("取消成功")
	return 
}

// 取消所有
func CancelAllQueue(SenderId string) (err error) {
	rd := db.RedisClient
	list, _ := rd.HGetAll(db.Ctx, StrCombine(KeyUserCron, SenderId)).Result()
	for _, value := range list {
		rd.ZRem(db.Ctx, KeyCrontab, value)
	}
	
	rd.Del(db.Ctx, StrCombine(KeyUserCron, SenderId))
	err = errors.New("已经取消所有提醒任务")
	return 
}

func QueryAllQueue(SenderId string) (map[string]string) {
	rd := db.RedisClient
	list, _ := rd.HGetAll(db.Ctx, StrCombine(KeyUserCron, SenderId)).Result()
	// fmt.Println(list)
	return list
}
