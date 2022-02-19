package main

import (
	"encoding/json"
	"flag"
	"log"
	"noty/qiyewechat"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func init() {
	f, err := os.OpenFile(getPwd()+"/noty.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if err != nil {
		log.Fatalln(err)
	}

	log.SetOutput(f)
}

func getPwd() string {
	err := os.Chdir(path.Dir(os.Args[0]))
	if err != nil {
		log.Fatalln(err)
	}

	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}

	return pwd
}

func main() {
	var typ, to, msg string
	flag.StringVar(&typ, "type", "server", "send|server")
	flag.StringVar(&msg, "msg", "ping", "message")
	flag.StringVar(&to, "to", "username", "to user")
	flag.Parse()

	f, err := os.Open(getPwd() + "/config.json")
	if err != nil {
		log.Fatalln(err)
	}

	var config qiyewechat.Config
	err = json.NewDecoder(f).Decode(&config)
	if err != nil {
		log.Fatalln(err)
	}

	if config.BaseURL == "" {
		config.BaseURL = "https://qyapi.weixin.qq.com/cgi-bin"
	}

	engin := gin.Default()

	agentFactory := new(qiyewechat.AgentFactory)
	for _, agent := range config.Agents {
		strID := strconv.FormatInt(agent.ID, 10)

		client := qiyewechat.NewQiyeWechatClien(config.BaseURL, config.CorpID, agent.Secret)
		if err = client.RefreshToken(); err != nil {
			log.Fatalln("[ERROR]", agent.ID, err)
		}
		ticker := time.NewTicker(time.Minute)
		go func(agent qiyewechat.AgentConfig) {
			for range ticker.C {
				if err = client.RefreshToken(); err != nil {
					log.Println("[ERROR]", agent.ID, err)
				}
			}
		}(agent)

		app := agentFactory.Create(config.CorpID, client, agent)
		engin.GET("/qiye-wechat/agents/"+strID, qiyewechat.VerifingHandler(app))
		engin.POST("/qiye-wechat/agents/"+strID, qiyewechat.MsgHandler(app))

		// 需要在 Proxy 控制该接口的访问，避免被恶意访问。Nginx 配置参考 nginx 文件夹里的 noty.conf
		engin.POST("/qiye-wechat/text-senders/"+strID, qiyewechat.TextHandler(app))
	}

	if err = engin.Run(config.Addr); err != nil {
		log.Println("[ERROR]", err)
	}
}
