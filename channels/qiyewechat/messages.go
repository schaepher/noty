package qiyewechat

type Message struct {
	Touser                 string        `json:"touser,omitempty"`
	Toparty                string        `json:"toparty,omitempty"`
	Totag                  string        `json:"totag,omitempty"`
	Msgtype                string        `json:"msgtype"`
	Agentid                int64         `json:"agentid"`
	Text                   *TextMessage  `json:"text,omitempty"`
	Image                  *ImageMessage `json:"image,omitempty"`
	Safe                   int64         `json:"safe,omitempty"`
	EnableIDTrans          int64         `json:"enable_id_trans,omitempty"`
	EnableDuplicateCheck   int64         `json:"enable_duplicate_check,omitempty"`
	DuplicateCheckInterval int64         `json:"duplicate_check_interval,omitempty"`
}

type TextMessage struct {
	Content string `json:"content"`
}

func (t TextMessage) Type() string {
	return "text"
}

type ImageMessage struct {
	MediaID string `json:"media_id"`
}

func (t ImageMessage) Type() string {
	return "img"
}
