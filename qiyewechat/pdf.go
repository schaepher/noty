package qiyewechat

import (
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"
)

const MB = 1024 * 1024

type PDFConvertConfig struct {
	URL      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
	PDFDir   string `json:"pdf_dir"`
}

type PDFHandler struct {
	app *EchoAgent
	cfg PDFConvertConfig
}

func NewPDFHandler(app *EchoAgent, cfg PDFConvertConfig) *PDFHandler {
	return &PDFHandler{
		app: app,
		cfg: cfg,
	}
}

func (h *PDFHandler) Handle(msg MsgContent) error {
	startTime := time.Now()

	cmd := strings.SplitN(msg.Content, " ", 2)

	go func() {
		time.Sleep(2 * time.Second)
		h.app.SendTextMessage(Message{
			Touser: msg.FromUsername,
			Text: &TextMessage{
				Content: "creating pdf: " + cmd[1],
			},
		})
	}()

	name := msg.Title
	if name == "" {
		name = time.Now().Format("20060102-150405")
	} else {
		name = time.Now().Format("20060102-15-") + name
	}

	v := url.Values{}
	v.Add("u", h.cfg.Username)
	v.Add("p", h.cfg.Password)
	v.Add("url", cmd[1])
	v.Add("WaitingTime", "100000")
	v.Add("uploadKey", name+".pdf")

	req := fmt.Sprintf("%s?%s", h.cfg.URL, v.Encode())

	resp, err := http.Get(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	all, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	filename := fmt.Sprintf("%s/%s", h.cfg.PDFDir, name)
	err = os.WriteFile(filename, all, fs.ModePerm)
	if err != nil {
		return err
	}

	h.app.SendTextMessage(Message{
		Touser: msg.FromUsername,
		Text: &TextMessage{
			Content: fmt.Sprintf("pdf created\nfilename: %s\nsize: %0.2f MB\ncost time: %d seconds", path.Base(filename), float64(len(all))/MB, int(time.Since(startTime).Seconds())),
		},
	})

	return nil
}
