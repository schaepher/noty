package tasks

// import (
// 	"noty/channels"
// 	"noty/channels/qiyewechat"
// 	"noty/eventbus"
// 	"time"
// )

// func init() {
// 	eventbus.Subscribe("pdf-in", CreatePDF)
// }

// func CreatePDF(msg qiyewechat.MsgContent) {
// 	eventbus.Publish("pdf-out", channels.Message{
// 		From:      "",
// 		Receivers: []string{msg.FromUsername},
// 		Title:     "PDF creation task received",
// 		Body:      "received at " + time.Now().Format("2006-01-02 15:04:05"),
// 	})

// 	go func() {
// 		if msg.MsgType == "link" {
// 			msg.Title, err = filenamify.FilenamifyV2(msg.Title, func(options *filenamify.Options) {
// 				options.Replacement = "_"
// 			})
// 			if err != nil {
// 				msg.Title = time.Now().Format("20060102-150405")
// 			}
// 			url, _ := url.QueryUnescape(msg.Url)
// 			msg.Content = "pdf " + url
// 			msg.Url = ""
// 		}

// 		if strings.HasPrefix(msg.Content, "http") {
// 			msg.Content = "pdf " + msg.Content
// 		}

// 		cmd := getCommandName(msg.Content)
// 		switch cmd {
// 		case "pdf":
// 			cfg := PDFConvertConfig{
// 				URL:      "http://127.0.0.1:8088/convert/html2pdf",
// 				Username: "",
// 				Password: "",
// 				PDFDir:   "",
// 			}
// 			NewPDFHandler(app, cfg).Handle(msg)
// 		default:
// 			a, _ := json.Marshal(msg)
// 			app.SendTextMessage(qiyewechat.Message{
// 				Touser: msg.FromUsername,
// 				Text: &qiyewechat.TextMessage{
// 					Content: string(a),
// 				},
// 			})
// 		}
// 	}()
// }
