package qiyewechat

import (
	"encoding/xml"
	"errors"

	"github.com/workweixin/weworkapi_golang/json_callback/wxbizjsonmsgcrypt"
)

type Agent interface {
	VerifyURL(msgSignature, timestamp, nonce, echostr string) ([]byte, error)
	DecryptMsg(msgSignature, timestamp, nonce string, jsonBody []byte) (msgContent MsgContent, err error)
	HandleMsg(msg MsgContent) (err error)
	SendTextMessage(msg Message) (err error)
}

type CommonAgent struct {
	corpID  string
	client  *QiyeWechatClient
	config  AgentConfig
	handler func(msg MsgContent) error
}

func NewComonAgent(corpID string, client *QiyeWechatClient, config AgentConfig, handler func(msg MsgContent) error) *CommonAgent {
	return &CommonAgent{
		corpID:  corpID,
		client:  client,
		config:  config,
		handler: handler,
	}
}

func (app *CommonAgent) VerifyURL(msgSignature, timestamp, nonce, echostr string) ([]byte, error) {
	wxcpt := wxbizjsonmsgcrypt.NewWXBizMsgCrypt(app.config.Token, app.config.EncodingAESKey, app.corpID, wxbizjsonmsgcrypt.JsonType)
	rs, err := wxcpt.VerifyURL(msgSignature, timestamp, nonce, echostr)
	if err != nil {
		return nil, errors.New(err.ErrMsg)
	}

	return rs, nil
}

func (app *CommonAgent) DecryptMsg(msgSignature, timestamp, nonce string, jsonBody []byte) (msgContent MsgContent, err error) {
	wxcpt := wxbizjsonmsgcrypt.NewWXBizMsgCrypt(app.config.Token, app.config.EncodingAESKey, app.corpID, wxbizjsonmsgcrypt.JsonType)
	msg, wxcpterr := wxcpt.DecryptMsg(msgSignature, timestamp, nonce, jsonBody)
	if wxcpterr != nil {
		err = errors.New("decrypting message: " + wxcpterr.ErrMsg)
		return
	}

	if err = xml.Unmarshal(msg, &msgContent); err != nil {
		err = errors.New("unmarshal decrypted messaage content: " + wxcpterr.ErrMsg)
		return
	}

	return
}

func (app *CommonAgent) HandleMsg(msg MsgContent) (err error) {
	return app.handler(msg)
}

func (app *CommonAgent) SendTextMessage(msg Message) (err error) {
	msg.Agentid = app.config.ID
	msg.Msgtype = msg.Text.Type()

	return app.client.SendMessage(msg)
}
