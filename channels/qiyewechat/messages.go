package qiyewechat

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"noty/eventbus"
)

type Message struct {
	Touser                 string        `json:"touser,omitempty"`
	Toparty                string        `json:"toparty,omitempty"`
	Totag                  string        `json:"totag,omitempty"`
	Msgtype                string        `json:"msgtype"`
	Agentid                int64         `json:"agentid"`
	Text                   *TextMessage  `json:"text,omitempty"`
	Image                  *ImageMessage `json:"image,omitempty"`
	Safe                   int64         `json:"safe,omitempty"`
	EnableIDTrans          int64         `json:"enable_id_trans,omitempty"`
	EnableDuplicateCheck   int64         `json:"enable_duplicate_check,omitempty"`
	DuplicateCheckInterval int64         `json:"duplicate_check_interval,omitempty"`
}

type TextMessage struct {
	Content string `json:"content"`
}

func (t TextMessage) Type() string {
	return "text"
}

type ImageMessage struct {
	MediaID string `json:"media_id"`
}

func (t ImageMessage) Type() string {
	return "img"
}

type SendMessageResponse struct {
	Errcode      int64  `json:"errcode"`
	Errmsg       string `json:"errmsg"`
	Invaliduser  string `json:"invaliduser"`
	Invalidparty string `json:"invalidparty"`
	Invalidtag   string `json:"invalidtag"`
	Msgid        string `json:"msgid"`
	ResponseCode string `json:"response_code"`
}

func (c *QiyeWechatAgent) SendTextMessage(msg Message) (err error) {
	msg.Agentid = c.cfg.ID
	msg.Msgtype = msg.Text.Type()

	return c.sendMessage(msg)
}

func (c *QiyeWechatAgent) sendMessage(msg Message) (err error) {
	accessToken := c.accessToken.Load().(AccessToken)

	j, _ := json.Marshal(msg)
	resp, err := http.Post(c.baseUrl+"/message/send?access_token="+accessToken.token,
		"application/json", bytes.NewReader(j))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var r SendMessageResponse
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return
	}

	if r.Errcode != 0 {
		err = errors.New(r.Errmsg)
		return
	}

	return
}

func (c *QiyeWechatAgent) HandleMsg(msg MsgContent) (err error) {
	return eventbus.Publish(c.cfg.FromUserTopic, msg)
}
