package qiyewechat

import (
	"context"
	"noty/eventbus"
	"noty/log"
	"time"

	"go.uber.org/atomic"

	"go.uber.org/zap"
)

const qyAPIBaseUrl = "https://qyapi.weixin.qq.com/cgi-bin"

// QiyeWechatAgent 用于与企业微信应用交互
type QiyeWechatAgent struct {
	ctx  context.Context
	stop context.CancelFunc

	cfg AgentConfig

	baseUrl     string
	accessToken *atomic.Value
	logger      *zap.Logger
}

func NewQiyeWechatAgent(ctx context.Context, cfg AgentConfig) (*QiyeWechatAgent, error) {
	ctx, cancel := context.WithCancel(ctx)

	client := &QiyeWechatAgent{
		ctx:         ctx,
		stop:        cancel,
		cfg:         cfg,
		baseUrl:     qyAPIBaseUrl,
		logger:      log.GetLogger(),
		accessToken: &atomic.Value{},
	}

	client.accessToken.Store(AccessToken{})

	err := client.start()
	if err != nil {
		client.logger.Error("start client", zap.Error(err))
		return nil, err
	}

	return client, nil
}

func (c *QiyeWechatAgent) start() error {
	logger := c.logger.With(zap.String("_loc", "[QiyeWechatAgent.start]"))

	if err := c.refreshToken(); err != nil {
		logger.Error("get first token", zap.Error(err))
		return err
	}

	go func() {
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := c.refreshToken(); err != nil {
					logger.Error("refresh token", zap.Error(err))
				}
			case <-c.ctx.Done():
				return
			}
		}
	}()

	eventbus.Subscribe(c.cfg.ToUserTopic, c.SendTextMessage)

	return nil
}

func (c *QiyeWechatAgent) refreshToken() (err error) {
	logger := c.logger.With(zap.String("_loc", "[QiyeWechatAgent.refreshToken]"))

	accessToken := c.accessToken.Load().(AccessToken)

	if accessToken.token == "" || accessToken.IsExpired() {
		token, expireIn, err := c.getToken(c.cfg.CorpID, c.cfg.Secret)
		if err != nil {
			logger.Error("get token by API", zap.Error(err))
			return err
		}

		accessToken = AccessToken{
			token:    token,
			expireAt: time.Now().Add(time.Duration(expireIn)*time.Second - time.Minute),
		}
		c.accessToken.Store(accessToken)
	}

	return nil
}

func (c *QiyeWechatAgent) Close() {
	c.stop()
}
