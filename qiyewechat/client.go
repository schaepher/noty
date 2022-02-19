package qiyewechat

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

// QiyeWechatClient 用于与企业微信 API 交互
type QiyeWechatClient struct {
	baseUrl    string
	corpID     string
	corpSecret string
	token      Token
}

func NewQiyeWechatClien(baseUrl string, corpID string, corpSecret string) *QiyeWechatClient {
	return &QiyeWechatClient{
		baseUrl: baseUrl, 
		corpID: corpID, 
		corpSecret: corpSecret,
	}
}

func (c *QiyeWechatClient) RefreshToken() (err error) {
	if c.token.token == "" || c.token.IsExpired() {
		token, expireIn, err := c.GetToken(c.corpID, c.corpSecret)
		if err != nil {
			return err
		}
		c.token.token = token
		c.token.expireAt = time.Now().Add(time.Duration(expireIn)*time.Second - time.Minute)
	}

	return
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

func (c QiyeWechatClient) GetToken(corpID string, corpSecret string) (token string, expireIn int64, err error) {
	req := getTokenRequest{corpID, corpSecret}
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

func (c QiyeWechatClient) SendMessage(msg Message) (err error) {
	j, _ := json.Marshal(msg)
	resp, err := http.Post(c.baseUrl+"/message/send?access_token="+c.token.token,
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
