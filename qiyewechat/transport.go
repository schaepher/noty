package qiyewechat

import (
	"encoding/json"
	"encoding/xml"
	"log"

	"github.com/gin-gonic/gin"
)

// VerifingHandler 用于企业微信 Agent 接收信息配置验证
func VerifingHandler(app Agent) func(c *gin.Context) {
	return func(c *gin.Context) {
		msgSignature := c.Query("msg_signature")
		timestamp := c.Query("timestamp")
		nonce := c.Query("nonce")
		echostr := c.Query("echostr")

		echoStr, err := app.VerifyURL(msgSignature, timestamp, nonce, echostr)
		if err != nil {
			log.Println("[ERROR] Verify: ", err)
		}

		c.Data(200, "plain", echoStr)
	}
}

type RawMessage struct {
	ToUserName string `json:"ToUserName" xml:"ToUserName"`
	AgentID    string `json:"AgentID" xml:"AgentID"`
	Encrypt    string `json:"Encrypt" xml:"Encrypt"`
}

type MsgContent struct {
	ToUsername   string `json:"ToUserName" xml:"ToUserName"`
	FromUsername string `json:"FromUserName" xml:"FromUserName"`
	CreateTime   uint32 `json:"CreateTime" xml:"CreateTime"`
	MsgType      string `json:"MsgType" xml:"MsgType"`
	Content      string `json:"Content" xml:"Content"`
	Title        string `json:"Title" xml:"Title"`
	Url          string `json:"Url" xml:"Url"`
	Msgid        uint64 `json:"MsgId" xml:"MsgId"`
	Agentid      uint32 `json:"AgentId" xml:"AgentId"`
}

type RawResponse struct {
	Encrypt      string `json:"encrypt" xml:"encrypt"`
	Timestamp    string `json:"timestamp" xml:"timestamp"`
	Nonce        string `json:"nonce" xml:"nonce"`
	MsgSignature string `json:"msgsignature" xml:"msgsignature"`
}

// MsgHandler 用于处理发送到 Agent 的消息
func MsgHandler(app Agent) func(*gin.Context) {
	return func(c *gin.Context) {
		msgSignature := c.Query("msg_signature")
		timestamp := c.Query("timestamp")
		nonce := c.Query("nonce")

		// 官方只发 xml 版本，但这里使用基于 json 的解码库，因此需要做一层转换。
		var raw RawMessage
		if err := xml.NewDecoder(c.Request.Body).Decode(&raw); err != nil {
			log.Println("[ERROR]", err)
			return
		}
		jsonBody, _ := json.Marshal(raw)

		content, err := app.DecryptMsg(msgSignature, timestamp, nonce, jsonBody)
		if err != nil {
			log.Println("[ERROR]", err)
		}
		err = app.HandleMsg(content)
		if err != nil {
			log.Println("[ERROR]", err)
		}

		c.Status(200)
	}
}

type TextMessageRequest struct {
	ToUsername string `json:"to_username"`
	Content    string `json:"content"`
}

// TextHandler 用于直接发送消息
func TextHandler(app Agent) func(*gin.Context) {
	return func(c *gin.Context) {
		var msg TextMessageRequest
		if err := json.NewDecoder(c.Request.Body).Decode(&msg); err != nil {
			log.Println("[ERROR]", err)
			return
		}

		err := app.SendTextMessage(Message{
			Touser: msg.ToUsername,
			Text: &TextMessage{
				Content: msg.Content,
			},
		})
		if err != nil {
			log.Println("[ERROR]", err)
			c.AbortWithError(500, err)
			return
		}

		c.Status(200)
	}
}
