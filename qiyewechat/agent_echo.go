package qiyewechat

type EchoAgent struct {
	common CommonAgent
}

func NewEchoAgent(corpID string, client *QiyeWechatClient, config AgentConfig) Agent {
	return &EchoAgent{common: CommonAgent{corpID: corpID, client: client, config: config}}
}

func (app *EchoAgent) VerifyURL(msgSignature, timestamp, nonce, echostr string) ([]byte, error) {
	return app.common.VerifyURL(msgSignature, timestamp, nonce, echostr)
}

func (app *EchoAgent) DecryptMsg(msgSignature, timestamp, nonce string, jsonBody []byte) (msgContent MsgContent, err error) {
	return app.common.DecryptMsg(msgSignature, timestamp, nonce, jsonBody)
}

func (app *EchoAgent) HandleMsg(msg MsgContent) (err error) {
	return app.SendTextMessage(Message{
		Touser: msg.FromUsername,
		Text: &TextMessage{
			Content: msg.Content,
		},
	})
}

func (app *EchoAgent) SendTextMessage(msg Message) (err error) {
	msg.Agentid = app.common.config.ID
	msg.Msgtype = msg.Text.Type()
	return app.common.client.SendMessage(msg)
}
