package qiyewechat

type AgentConfig struct {
	ID             int64  `json:"id"`
	Secret         string `json:"secret"`
	Token          string `json:"token"`
	EncodingAESKey string `json:"encoding_aes_key"`
	Type           string `json:"type"`
}

type Config struct {
	Addr    string        `json:"addr"`
	BaseURL string        `json:"base_url"`
	CorpID  string        `json:"corp_id"`
	Agents  []AgentConfig `json:"agents"`
}
