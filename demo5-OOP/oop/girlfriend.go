package oop

import (
	"fmt"
	"strconv"
)

type GirlFriend struct {
	height    int    // 身高
	weight    int    // 体重
	age       int    // 年龄
	name      string // 姓名
	content   string // 结束语
	greetings string // 问候语
}

func (gf *GirlFriend) SetHeight(v int) {
	gf.height = v
}

func (gf *GirlFriend) SetWeight(v int) {
	gf.weight = v
}

func (gf *GirlFriend) SetAge(v int) {
	gf.age = v
}

func (gf *GirlFriend) SetName(v string) {
	gf.name = v
}

func (gf *GirlFriend) SetContent(v string) *GirlFriend {
	gf.content = v
	return gf
}

func (gf *GirlFriend) SetGreeting(v string) *GirlFriend {
	gf.greetings = v
	return gf
}

func NewOne() *GirlFriend {
	return &GirlFriend{}
}

// 注意这里小写字母开头属性是可以的哟
func (gf *GirlFriend) Show() {
	fmt.Println(gf.greetings + "我是" + gf.name + ",今年" + strconv.Itoa(gf.age) + "岁，身高" + strconv.Itoa(gf.height) + "cm,体重" + strconv.Itoa(gf.weight) + "kg。" + gf.content)
}
