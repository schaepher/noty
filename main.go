package main

import (
	"context"
	"encoding/json"
	"flag"
	"os"
	"path"
	"strconv"

	"noty/agents"
	"noty/channels/qiyewechat"
	"noty/log"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func init() {
	f, err := os.OpenFile(getPwd()+"/noty.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}

	log.Init(&log.Config{
		Filename: f.Name(),
	})
}

func getPwd() string {
	err := os.Chdir(path.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}

	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return pwd
}

func main() {
	var typ, to, msg string
	flag.StringVar(&typ, "type", "server", "send|server")
	flag.StringVar(&msg, "msg", "ping", "message")
	flag.StringVar(&to, "to", "username", "to user")
	flag.Parse()

	logger := log.GetLogger()
	defer logger.Sync()

	f, err := os.Open(getPwd() + "/config.json")
	if err != nil {
		logger.Error("open config", zap.Error(err))
		panic(err)
	}

	var config qiyewechat.Config
	err = json.NewDecoder(f).Decode(&config)
	if err != nil {
		logger.Error("decode config", zap.Error(err))
		panic(err)
	}

	ctx := context.Background()
	engin := gin.Default()

	for _, agent := range config.Agents {
		strID := strconv.FormatInt(agent.ID, 10)

		client, err := qiyewechat.NewQiyeWechatClien(ctx, config.CorpID, agent.Secret)
		if err != nil {
			logger.Error("create qiye wechat client", zap.Error(err))
			panic(err)
		}

		echoAgent := agents.NewEchoAgent(config.CorpID, client, agent)
		engin.GET("/qiye-wechat/agents/"+strID, qiyewechat.VerifingHandler(echoAgent))
		engin.POST("/qiye-wechat/agents/"+strID, qiyewechat.MsgHandler(echoAgent))

		// 需要在 Proxy 控制该接口的访问，避免被恶意访问。Nginx 配置参考 nginx 文件夹里的 noty.conf
		engin.POST("/qiye-wechat/text-senders/"+strID, qiyewechat.TextHandler(echoAgent))
	}

	if err = engin.Run(config.ServerAddr); err != nil {
		logger.Error("engin run", zap.Error(err))
		panic(err)
	}
}
