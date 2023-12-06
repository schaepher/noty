package qiyewechat

import (
	"encoding/json"
	"net/url"
	"strings"
	"time"

	"github.com/flytam/filenamify"
)

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

func getCommandName(s string) string {
	cmd := strings.SplitN(s, " ", 2)
	return cmd[0]
}

func (app *EchoAgent) HandleMsg(msg MsgContent) (err error) {
	app.SendTextMessage(Message{
		Touser: msg.FromUsername,
		Text: &TextMessage{
			Content: "received at " + time.Now().Format("2006-01-02 15:04:05"),
		},
	})

	go func() {
		if msg.MsgType == "link" {
			msg.Title, err = filenamify.FilenamifyV2(msg.Title, func(options *filenamify.Options) {
				options.Replacement = "_"
			})
			if err != nil {
				msg.Title = time.Now().Format("20060102-150405")
			}
			url, _ := url.QueryUnescape(msg.Url)
			msg.Content = "pdf " + url
			msg.Url = ""
		}

		if strings.HasPrefix(msg.Content, "http") {
			msg.Content = "pdf " + msg.Content
		}

		cmd := getCommandName(msg.Content)
		switch cmd {
		case "pdf":
			cfg := PDFConvertConfig{
				URL:      "",
				Username: "",
				Password: "",
				PDFDir:   "",
			}
			NewPDFHandler(app, cfg).Handle(msg)
		default:
			a, _ := json.Marshal(msg)
			app.SendTextMessage(Message{
				Touser: msg.FromUsername,
				Text: &TextMessage{
					Content: string(a),
				},
			})
		}
	}()

	return nil
}

func (app *EchoAgent) SendTextMessage(msg Message) (err error) {
	msg.Agentid = app.common.config.ID
	msg.Msgtype = msg.Text.Type()
	return app.common.client.SendMessage(msg)
}
