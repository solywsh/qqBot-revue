package cq

import (
	"strconv"
	"strings"
)

// MsgAddApiToken 添加token
func (cpf *PostForm) MsgAddApiToken() {
	//gdb := db.NewDB()
	if res, token := gdb.InsertRevueApiToken(strconv.Itoa(cpf.UserId), 4); res {
		_, _ = cpf.SendMsg(cpf.MessageType, "创建成功,你的token是:"+token+"\n注意,该token只能给自己发送消息")
	} else {
		_, _ = cpf.SendMsg(cpf.MessageType, "创建失败,你已经创建过了,token是:"+token)
	}
}

// MsgResetApiToken 重置token
func (cpf *PostForm) MsgResetApiToken() {
	if res, token := gdb.ResetRevueApiToken(strconv.Itoa(cpf.UserId)); res {
		_, _ = cpf.SendMsg(cpf.MessageType, "重置成功,你的token是:"+token+"\n注意,该token只能给自己发送消息")
	} else {
		_, _ = cpf.SendMsg(cpf.MessageType, "重置失败,请先创建token")
	}
}

// MsgDeleteApiToken 删除token
func (cpf *PostForm) MsgDeleteApiToken() {
	if res, token := gdb.DeleteRevueApiToken(strconv.Itoa(cpf.UserId)); res {
		_, _ = cpf.SendMsg(cpf.MessageType, token+"删除成功")
	} else {
		_, _ = cpf.SendMsg(cpf.MessageType, "删除失败,可能数据库没有对应的信息")
	}
}

func (cpf *PostForm) PrivateEvent() {
	switch {
	// 发送菜单
	case cpf.Message == "/help":
		cpf.SendMenu()
	// /getToken 创建对应token
	case cpf.Message == "/getToken":
		cpf.MsgAddApiToken() // 添加对应apiToken
	// /resetToken 重置对应token
	case cpf.Message == "/resetToken":
		cpf.MsgResetApiToken()
	// /deleteToken 删除对应token
	case cpf.Message == "/deleteToken":
		cpf.MsgDeleteApiToken()
	// 搜索答案
	case strings.HasPrefix(cpf.Message, "搜索答案"):
		cpf.GetAnswer()
	}
}
