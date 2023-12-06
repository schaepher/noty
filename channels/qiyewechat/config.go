package qiyewechat

type AgentConfig struct {
	ID             int64  `json:"id"`
	Secret         string `json:"secret"`
	Token          string `json:"token"`
	EncodingAESKey string `json:"encoding_aes_key"`
	Topic          string `json:"topic"`
}

type Config struct {
	ServerAddr string        `json:"addr"`
	CorpID     string        `json:"corp_id"`
	Agents     []AgentConfig `json:"agents"`
}
