package cq

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/thedevsaddam/gojsonq"
	"strconv"
	"strings"
)

// GroupEvent 群消息事件
func (cpf *PostForm) GroupEvent() {
	// cpf.RepeatOperation() // 对adminUSer复读防止风控
	//fmt.Println("收到群消息:", cpf.Message, cpf.UserId)
	switch {
	// demo
	case cpf.Message == "叫两声":
		_, _ = cpf.SendGroupMsg("汪汪")
	case strings.HasPrefix(cpf.Message, "查找音乐"):
		// 查找音乐
		cpf.FindMusicEvent()
	case cpf.Message == "开始添加":
		// 触发添加自动回复
		cpf.KeywordsReplyAddEvent(1, 0)
	case strings.HasPrefix(cpf.Message, "删除自动回复:"):
		// 删除自动回复
		cpf.KeywordsReplyDeleteEvent()
	case strings.HasPrefix(cpf.Message, "搜索答案"):
		// 搜索答案
		cpf.GetAnswer()
	default:
		// 添加自动回复(关键词/回复内容)
		if res, kr := gdb.GetKeywordsReplyFlag(strconv.Itoa(cpf.UserId)); res {
			cpf.KeywordsReplyAddEvent(kr.Flag+1, kr.ID)
		} else {
			// 自动回复
			cpf.AutoGroupMsg()
		}
	}
}

// MsgEvent 对message事件进行相应
func (cpf *PostForm) MsgEvent() {
	// 判断是否为adminUser且为命令
	if cpf.JudgmentAdminUser() {
		cpf.AdminEvent() //执行对应admin命令事件
		return           // 如果是执行之后直接返回，不再继续响应
	}
	// 群消息进行响应
	if cpf.MessageType == "group" {
		if ok := cpf.JudgeListenGroup(); ok {

			// 发送菜单
			if cpf.Message == "/help" {
				cpf.SendMenu()
				return
			}
			// 如果是监听qq群列表的才做出相应
			cpf.GroupEvent()
		}
		// 对不是监听qq群列表的消息做出相应

		// do event
	}
	// 对私聊进行响应
	if cpf.MessageType == "private" {
		// 发送菜单
		if cpf.Message == "/help" {
			cpf.SendMenu()
		}
		// /getToken 创建对应token
		if cpf.Message == "/getToken" {
			cpf.MsgAddApiToken() // 添加对应apiToken
			return
		}
		// /resetToken 重置对应token
		if cpf.Message == "/resetToken" {
			cpf.MsgResetApiToken()
			return
		}
		// /deleteToken 删除对应token
		if cpf.Message == "/deleteToken" {
			cpf.MsgDeleteApiToken()
			return
		}
	}
}

// SendMenu 发送命令菜单
func (cpf *PostForm) SendMenu() {
	s := "revue提供以下命令:\n"
	if cpf.MessageType == "private" {
		s += "revueApi 相关(私聊执行命令):\n"
		s += "[/getToken] 获取token\n"
		s += "[/resetToken] 重置token\n"
		s += "[/deleteToken] 删除token"
		_, _ = cpf.SendMsg(cpf.MessageType, s)
	} else if cpf.MessageType == "group" {
		s += "群聊菜单:\n"
		s += "[开始添加] 添加自动回复\n"
		s += "[删除自动回复:{关键词}] 删除自动回复\n"
		s += "[查找音乐{关键词}] 查找音乐(暂时只支持163)\n"
		s += "[搜索答案{关键词}] 搜索答案"
		_, _ = cpf.SendMsg(cpf.MessageType, s)
	}
}

// ProblemRepository 搜索题库
func ProblemRepository(question string) string {
	client := resty.New()
	post, err := client.R().SetQueryParams(map[string]string{
		"q": question,
	}).Post("http://api.902000.xyz:88/wkapi.php")
	if err != nil {
		return "题目请求失败"
	}
	postJSON := gojsonq.New().JSONString(post.String())
	if postJSON.Reset().Find("code").(float64) == float64(1) {
		return fmt.Sprintf("问题:" + postJSON.Reset().Find("tm").(string) + "\n" +
			"答案:" + postJSON.Reset().Find("answer").(string))
	} else {
		return "没有找到相关问题"
	}
}

// GetAnswer 搜索题目答案
func (cpf *PostForm) GetAnswer() {
	question := strings.TrimPrefix(cpf.Message, "搜索答案")
	ans := ProblemRepository(question)
	_, _ = cpf.SendMsg(cpf.MessageType, ans)
}