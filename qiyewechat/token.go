package qiyewechat

import "time"

type Token struct {
	token    string
	expireAt time.Time
}

func (t Token) IsExpired() bool {
	return time.Now().After(t.expireAt)
}
