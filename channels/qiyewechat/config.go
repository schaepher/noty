package qiyewechat

type AgentConfig struct {
	ID             int64  `json:"id"`
	CorpID         string `json:"corp_id"`
	Secret         string `json:"secret"`
	Token          string `json:"token"`
	EncodingAESKey string `json:"encoding_aes_key"`
	FromUserTopic  string `json:"from_user_topic"` // 应用收到个人微信的消息会发送到这个 Topic
	ToUserTopic    string `json:"to_user_topic"`   // 服务端发送给个人微信的消息需要发送到这个 Topic
}

type Config struct {
	ServerAddr string        `json:"addr"`
	Agents     []AgentConfig `json:"agents"`
}
