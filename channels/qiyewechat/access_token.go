package qiyewechat

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type AccessToken struct {
	token    string
	expireAt time.Time
}

func (t AccessToken) IsExpired() bool {
	return time.Now().After(t.expireAt)
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

func (c *QiyeWechatAgent) getToken(corpID string, agentSecret string) (token string, expireIn int64, err error) {
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