package qiyewechat

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"noty/log"
	"time"

	"go.uber.org/zap"
)

const qyAPIBaseUrl = "https://qyapi.weixin.qq.com/cgi-bin"

// QiyeWechatClient 用于与企业微信 API 交互
type QiyeWechatClient struct {
	ctx  context.Context
	stop context.CancelFunc

	baseUrl     string
	corpID      string
	agentSecret string
	apiToken    Token

	logger *zap.Logger
}

func NewQiyeWechatClien(ctx context.Context, corpID string, agentSecret string) (*QiyeWechatClient, error) {
	ctx, cancel := context.WithCancel(ctx)

	client := &QiyeWechatClient{
		ctx:         ctx,
		stop:        cancel,
		baseUrl:     qyAPIBaseUrl,
		corpID:      corpID,
		agentSecret: agentSecret,
		logger:      log.GetLogger(),
	}

	err := client.start()
	if err != nil {
		client.logger.Error("start client", zap.Error(err))
		return nil, err
	}

	return client, nil
}

func (c *QiyeWechatClient) start() error {
	logger := c.logger.With(zap.String("_loc", "[QiyeWechatClient.start]"))

	if err := c.refreshToken(); err != nil {
		logger.Error("get first token", zap.Error(err))
		return err
	}

	go func() {
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := c.refreshToken(); err != nil {
					logger.Error("refresh token", zap.Error(err))
				}
			case <-c.ctx.Done():
				return
			}
		}
	}()

	return nil
}

func (c *QiyeWechatClient) refreshToken() (err error) {
	logger := c.logger.With(zap.String("_loc", "[QiyeWechatClient.refreshToken]"))

	if c.apiToken.token == "" || c.apiToken.IsExpired() {
		token, expireIn, err := c.GetToken(c.corpID, c.agentSecret)
		if err != nil {
			logger.Error("get token by API", zap.Error(err))
			return err
		}
		c.apiToken.token = token
		c.apiToken.expireAt = time.Now().Add(time.Duration(expireIn)*time.Second - time.Minute)
	}

	return nil
}

func (c *QiyeWechatClient) Close() {
	c.stop()
}

type getTokenRequest struct {
	CorpID     string `json:"corpid"`
	CorpSecret string `json:"corpsecret"`
}

type getTokenResponse struct {
	ErrCode     int    `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expire_in"`
}

func (c *QiyeWechatClient) GetToken(corpID string, agentSecret string) (token string, expireIn int64, err error) {
	req := getTokenRequest{corpID, agentSecret}
	reqJ, _ := json.Marshal(req)
	resp, err := http.Post(c.baseUrl+"/gettoken", "application/json", bytes.NewReader(reqJ))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var respData getTokenResponse
	err = json.NewDecoder(resp.Body).Decode(&respData)
	if err != nil {
		return
	}

	if respData.ErrCode != 0 {
		err = errors.New(respData.ErrMsg)
		return
	}

	return respData.AccessToken, respData.ExpiresIn, nil
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

func (c *QiyeWechatClient) SendMessage(msg Message) (err error) {
	j, _ := json.Marshal(msg)
	resp, err := http.Post(c.baseUrl+"/message/send?access_token="+c.apiToken.token,
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
