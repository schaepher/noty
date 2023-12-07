package qiyewechat

type AgentConfig struct {
	ID             int64  `json:"id"`
	Secret         string `json:"secret"`
	Token          string `json:"token"`
	EncodingAESKey string `json:"encoding_aes_key"`
	ProducerTopic  string `json:"producer_topic"` // 来自 agent 的消息会发送到这个 Topic
	ReceiverTopic  string `json:"receiver_topic"` // 想要发送给 agent 的消息需要发送到这个 Topic
}

type Config struct {
	ServerAddr string        `json:"addr"`
	CorpID     string        `json:"corp_id"`
	Agents     []AgentConfig `json:"agents"`
}
