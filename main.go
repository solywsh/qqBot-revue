package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/solywsh/qqBot-revue/conf"
	"github.com/solywsh/qqBot-revue/cq"
	"github.com/solywsh/qqBot-revue/db"
	"github.com/thedevsaddam/gojsonq"
	"io/ioutil"
	"net/http"
	"strconv"
)

// 定义全局变量
var (
	yamlConf *conf.Config
	gdb      *db.GormDb
)

// 初始化配置
func init() {
	yamlConf, _ = conf.NewConf("./config.yaml") // 得到配置文件
	gdb = db.NewDB()                            // 初始化操作数据库
}

// revue发送消息api,用于接收
type revueApiPost struct {
	Token      string `json:"token"`       // 加密后的密钥
	UserId     string `json:"user_id"`     // qq号
	Message    string `json:"message"`     // 消息内容
	AutoEscape string `json:"auto_escape"` // 消息内容是否作为纯文本发送(即不解析CQ码)
}

// 监听go-cqhttp动作并以此做出反应
func listenFromCqhttp(c *gin.Context) {
	var cpf cq.PostForm
	if c.ShouldBindBodyWith(&cpf, binding.JSON) == nil {
		// 对message事件进行响应
		if cpf.PostType == "message" {
			go cpf.MsgEvent()
		}
	}
}

//  SHA1 加密进行鉴权
func hmacSHA1Encrypt(encryptKey string, encryptText []byte) string {
	key := []byte(encryptKey)
	mac := hmac.New(sha1.New, key)
	mac.Write(encryptText)
	var str = hex.EncodeToString(mac.Sum(nil))
	return str
}

// gin中间件,如果开启反向鉴权(reverseAuthentication)时,对数据进行验证
func ginReverseAuthentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		if yamlConf.ReverseAuthentication.Enable {
			body, _ := ioutil.ReadAll(c.Request.Body)
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body)) // 重设body
			headerXSignature := c.Request.Header.Get("X-Signature")  // sha1签名
			//headerXSelfId := c.Request.Header.Get("X-Self-Id")      // 发送消息的qq
			if headerXSignature[len("sha1="):] != hmacSHA1Encrypt(yamlConf.ReverseAuthentication.Secret, body) {
				c.JSON(
					http.StatusForbidden,
					gin.H{
						"code": http.StatusForbidden,
						"msg":  "密钥错误",
					},
				)
				c.Abort()
			}
		}
	}
}

//  监听revue发送私聊消息的接口
func listenFromSendPrivateMsg(c *gin.Context) {
	// 如果revue没开启直接结束
	if !yamlConf.Revue.Enable {
		c.Abort()
	}
	var rap revueApiPost
	if c.ShouldBindBodyWith(&rap, binding.JSON) == nil {
		// do event
		var cpf cq.PostForm
		cpf.UserId, _ = strconv.Atoi(rap.UserId)
		msg, err := cpf.SendMsg("private", cpf.Message)
		if err != nil {
			c.JSON(
				http.StatusInternalServerError,
				gin.H{
					"code": http.StatusInternalServerError,
					"msg":  msg, // 返回错误信息
				},
			)
		} else {
			c.JSON(
				http.StatusOK,
				gin.H{
					"code": http.StatusOK,
					"msg":  msg, // 正确返回msg id
				},
			)
		}

	} else {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"code": http.StatusBadRequest,
				"msg":  "请求参数不能识别",
			},
		)
	}
}

// revue接口中间件,对发送的token进行验证
func ginRevueAuthentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		if yamlConf.ReverseAuthentication.Enable {
			body, _ := ioutil.ReadAll(c.Request.Body)
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body)) // 重设body
			pJson := gojsonq.New().JSONString(string(body))
			pToken := pJson.Reset().Find("token").(string)
			pUserId := pJson.Reset().Find("user_id").(string)
			if res, permission := gdb.SearchRevueApiToken(pUserId, pToken); !res || permission < 4 {
				c.JSON(
					http.StatusForbidden,
					gin.H{
						"code": http.StatusForbidden,
						"msg":  "密钥错误",
					},
				)
				c.Abort() // 结束会话
			}
		}
	}
}

func main() {
	gin.DisableConsoleColor() // 不显示彩色日志
	router := gin.Default()
	// 监听动作并做出反应
	router.POST("/", ginReverseAuthentication(), listenFromCqhttp)
	// 监听revue提供发送消息的接口
	router.POST("/send_private_msg", ginRevueAuthentication(), listenFromSendPrivateMsg)
	err := router.Run("0.0.0.0:" + yamlConf.ListenPort)
	if err != nil {
		return
	}
}
