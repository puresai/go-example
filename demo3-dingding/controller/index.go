package controller

import (
	"fmt"
	"strings"
	"regexp"
	"time"
	"errors"
	"strconv"
	"io/ioutil"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	redis "github.com/go-redis/redis/v8"

	"local.com/sai0556/demo3-dingding/util"
	"local.com/sai0556/demo3-dingding/db"
)

/*
{
	"conversationId":"xxx",
	"atUsers":[
		{
			"dingtalkId":"xxx"
		}],
	"chatbotUserId":"xxx",
	"msgId":"xxx",
	"senderNick":"sai0556",
	"isAdmin":false,
	"sessionWebhookExpiredTime":1594978626787,
	"createAt":1594973226742,
	"conversationType":"2",
	"senderId":"xxx",
	"conversationTitle":"智能备忘录",
	"isInAtList":true,
	"sessionWebhook":"xxx",
	"text":{
		"content":" hello demo3-dingding"
	},
	"msgtype":"text"
}
curl 'http://127.0.0.1:8080/dingding' -H 'Content-Type: application/json' -d '{"content":{"conversationId":"cidxFg0R3qiLXU/CPOUwAtuRg==","atUsers":[{"dingtalkId":"$:LWCP_v1:$RoLrlkUFpZV+V2It1DLeEnasfhMsuID8"}],"chatbotUserId":"$:LWCP_v1:$RoLrlkUFpZV+V2It1DLeEnasfhMsuID8","msgId":"msgqxbfDGcsQs6hqV4THRu7Xw==","senderNick":"Sai 王泽涛","isAdmin":false,"sessionWebhookExpiredTime":1594978585443,"createAt":1594973185398,"conversationType":"2","senderId":"$:LWCP_v1:$bwfM2M8yXVRSCX80c7MNQA==","conversationTitle":"智能备忘录","isInAtList":true,"sessionWebhook":"https://oapi.dingtalk.com/robot/sendBySession?session=5fb94838e8d31e7e976faf450eff58be","text":{"content":"7天哈哈哈哈"},"msgtype":"text"}}
*/

type DingDingMsgContent struct {
	SenderNick string `json:"senderNick"`
	SenderId string `json:"senderId"`
	Text struct {
		Content string `json:"content"`
	} `json:"text"`
}

// 接收钉钉消息，最好加入验证，防止伪造消息
func DingDing(c *gin.Context) {
	data, _ := ioutil.ReadAll(c.Request.Body)
	form := DingDingMsgContent{}
	err := json.Unmarshal([]byte(data), &form)
	// err := c.ShouldBindJSON(&form)
	if  err != nil {
		fmt.Println(err)
		util.SendDD("你说的我听不懂")
		return
	}

	if er := parseContent(form); er != nil {
		// fmt.Println(er)
		util.SendDD(er.Error())
		return
	} 
	if err := tips(form); err != nil { 
		text, err := util.RobotApi(form.Text.Content)
		if err != nil {
			util.SendDD("你说的我听不懂")
			return
		}
		util.SendDD(text)
	}
	ApiResponse(c, 0, "success", nil)
}

func parseContent(form DingDingMsgContent) (err error) {
	str := form.Text.Content
	redis := db.RedisClient
	fmt.Println(str)

	// 要先绑定哟，不然无法@到对应的人
	index := strings.Index(str, "绑定手机")
	if index > -1 {
		reg := regexp.MustCompile("1[0-9]{10}")
		res := reg.FindAllString(str, 1)
		if len(res) < 1 || res[0] == "" {
			err = errors.New("手机格式不正确")
			return
		}
		redis.HSet(db.Ctx, util.KeyDingDingID, form.SenderId, res[0])
		util.SendDD("绑定成功")
		return
	}

	hExist := redis.HExists(db.Ctx, util.KeyDingDingID, form.SenderId)
	if !hExist.Val() {
		err = errors.New("绑定手机号才能精确提醒哦，发送--绑定手机 13456567878--@我即可")
		return 
	}

	index = strings.Index(util.StrSub(str, 0, 10), "我的提醒")
	fmt.Println(index, "---", util.StrSub(str, 0, 6))
	if index > -1 {
		www := util.QueryAllQueue(form.SenderId);
		if len(www) < 1 {
			err = errors.New("暂无任务")
			return
		} 
		msg := ""
		for key,value := range www {
			fmt.Println(strings.Index(value, "@"))
			value := value[0:strings.Index(value, "@")]
			fmt.Println(value)
			msg = util.StrCombine(msg, "任务id：", key, "，任务内容：", value, "{br}")
		}
		err = errors.New(msg)
		return
	}

	index = strings.Index(util.StrSub(str, 0, 10), "查看任务")
	fmt.Println(index, "---", util.StrSub(str, 0, 6))
	if index > -1 {
		www := util.QueryAllQueue(form.SenderId);
		if len(www) < 1 {
			err = errors.New("暂无任务")
			return
		} 
		msg := ""
		for key,value := range www {
			fmt.Println(strings.Index(value, "@"))
			value := value[0:strings.Index(value, "@")]
			fmt.Println(value)
			msg = util.StrCombine(msg, "任务id：", key, "，任务内容：", value, "{br}")
		}
		err = errors.New(msg)
		return
	}

	index = strings.Index(util.StrSub(str, 0, 10), "取消所有任务")
	fmt.Println(index, "---", util.StrSub(str, 0, 6))
	if index > -1 {
		if er := util.CancelAllQueue(form.SenderId); er != nil {
			err = er
			return
		}
		err = errors.New("取消成功")
		return
	}

	index = strings.Index(util.StrSub(str, 0, 10), "取消")
	if index > -1 {
		reg := regexp.MustCompile("[a-z0-9]{32}")
		res := reg.FindAllString(str, 1)
		if len(res) < 1 {
			err = errors.New("任务id不正确")
			return
		}
		if er := util.CancelQueue(res[0], form.SenderId); er != nil {
			err = er
			return
		}
		err = errors.New("取消成功")
		return

	}

	return
}

// 提醒内容
func tips(form DingDingMsgContent) (err error)  {
	rd := db.RedisClient
	str := form.Text.Content

	mobile := rd.HGet(db.Ctx, util.KeyDingDingID, form.SenderId).Val()
	key := util.KeyWords
	list, _ := rd.HGetAll(db.Ctx, key).Result()
	now := time.Now().Unix()
	tipsType := 1
	k := ""
	v := ""
	fmt.Println("str", str)

	index := 0

	for key, value := range list {
		index = util.UnicodeIndex(str, key)
		if index > -1 && util.StrLen(key) > util.StrLen(k) {
			fmt.Println("index", index, str, key, value)
			k = key
			v = value
		}
	}

	msg := ""
	var score int64
	if k != "" {
		kLen := util.StrLen(k)
		msg = util.StrSub(str, index+kLen)

		val := strings.Split(v, "|")
		unit := val[1]
		units,_ := strconv.Atoi(unit)

		switch val[0] {
			// 多少时间后
			case "1":
				reg := regexp.MustCompile("[0-9]{1,2}")
				res := reg.FindAllString(str, 1)
				minute, _ := strconv.Atoi(res[0])
				score = now + int64(units*minute)
			// 每周
			case "2":
				reg := regexp.MustCompile("[0-9]{1,2}")
				res := reg.FindAllString(util.StrSub(msg, 0, 7), -1)
				hour := 9
				minute := 0
				if len(res) > 0 {
					hour, _ = strconv.Atoi(res[0])
				}
				if len(res) > 1 {
					minute, _ = strconv.Atoi(res[1])
				}
				now = util.GetWeekTS(int64(units))
				score = now + int64(60*minute + 3600*hour)
				tipsType = 2
				
			// 下周
			case "3":
				reg := regexp.MustCompile("[0-9]{1,2}")
				res := reg.FindAllString(util.StrSub(msg, 0, 7), -1)
				hour := 9
				minute := 0
				if len(res) > 0 {
					hour, _ = strconv.Atoi(res[0])
				}
				if len(res) > 1 {
					minute, _ = strconv.Atoi(res[1])
				}
				now = util.TodayTS()
				score = now + int64(60*minute + 3600*hour + units*86400)
			case "4":
				reg := regexp.MustCompile("[0-9]{1,2}")
				res := reg.FindAllString(util.StrSub(msg, 0, 7), -1)
				hour := 9
				minute := 0
				if len(res) > 0 {
					hour, _ = strconv.Atoi(res[0])
				}
				if len(res) > 1 {
					minute, _ = strconv.Atoi(res[1])
				}
				now = util.TodayTS() + 86400*int64(units)
				score = now + int64(60*minute + 3600*hour)
			case "-1": 
				reg := regexp.MustCompile("[0-9]{1,10}")
				res := reg.FindAllString(util.StrSub(msg, 0, 7), -1)
				fmt.Println("res", res)
				hour := 9
				minute := 0
				if len(res) > 0 {
					hour, _ = strconv.Atoi(res[0])
				}
				if len(res) > 1 {
					minute, _ = strconv.Atoi(res[1])
				}
				now = util.TodayTS() + 86400
				score = now + int64(60*minute + 3600*hour)
				fmt.Println(now, score, minute, hour)
				tipsType = 3
			default:
		}
	} else {
		reg := regexp.MustCompile("(([0-9]{4})[-|/|年])?([0-9]{1,2})[-|/|月]([0-9]{1,2})日?")
		pi := reg.FindAllStringSubmatch(str, -1)
		if (len(pi) > 0 ) {
			date := pi[0]
			if date[2] == "" {
				date[2] = "2020"
			}
			location, _ := time.LoadLocation("Asia/Shanghai")
			tm2, _ := time.ParseInLocation("2006/01/02", fmt.Sprintf("%s/%s/%s", date[2], date[3], date[4]), location)
			score = util.GetZeroTime(tm2).Unix()

			msg = reg.ReplaceAllString(str, "")
			fmt.Println(msg)
			
		} else {
			msg = str
			score = util.TodayTS()
		}
		
		reg = regexp.MustCompile("[0-9]{1,10}")
		res := reg.FindAllString(util.StrSub(msg, 0, 7), -1)
		fmt.Println("res", res)
		hour := 9
		minute := 0
		if len(res) >= 1 {
			hour, _ = strconv.Atoi(res[0])
			fmt.Println("hour", hour, minute)
		}
		if len(res) > 1 {
			minute, _ = strconv.Atoi(res[1])
		}
		score += int64(60*minute + 3600*hour)
	}

	if msg == "" {
		err = errors.New("你说啥")
		return
	}
	index = util.UnicodeIndex(msg, "提醒我")
	index2 := util.UnicodeIndex(msg, "提醒")
	if index2 < 0 {
		err = errors.New("大哥，要我提醒你干啥呢？请发送--下周一13点提醒我写作业")
		return
	}

	if index < 0 && index2 > -1 {
		msg = util.StrSub(msg, index2+2)
	} else {
		msg = util.StrSub(msg, index+3)
	}
 
	fmt.Println(msg, mobile)
	msg = util.StrCombine(msg, "@", mobile)

	fmt.Println(score, msg, tipsType, err)
	if err != nil {
		util.SendDD(err.Error())
		return
	}

	member := util.StrCombine(strconv.Itoa(tipsType), msg)
	rd.ZAdd(db.Ctx, util.KeyCrontab, &redis.Z{
		Score: float64(score),
		Member: member,
	})

	uniqueKey := util.Md5(member)
	rd.HSet(db.Ctx, util.StrCombine(util.KeyUserCron, form.SenderId), uniqueKey, member)
	util.SendDD(fmt.Sprintf("设置成功(取消请回复：取消任务%s)--%s提醒您%s", uniqueKey, time.Unix(score, 0).Format("2006/01/02 15:04:05"), msg))
	return 
}


// 返回
type Response struct {
    Code int `json:"code"`
    Message string `json:"message"`
    Data interface{} `json:"data"`
}

// api返回结构
func ApiResponse(c *gin.Context, code int, message string, data interface{}) {
    c.JSON(http.StatusOK, Response{
        Code: code,
        Message: message,
        Data: data,
    })
}

func HealthCheck(c *gin.Context) {
    ApiResponse(c, 0, "success", nil)
}