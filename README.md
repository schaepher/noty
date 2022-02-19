# noty

在个人微信上借助个人微信与企业微信的绑定达到与服务器交互的效果。无需真实拥有企业。

个人微信 -> 企业微信 -> 个人服务器  
个人微信 <- 企业微信 <- 个人服务器  

## 大致步骤

1. 在企业微信上创建企业
2. 通过个人微信扫码绑定上一步创建的企业的微信插件
3. 创建应用
4. 启动服务
5. 设置应用接收信息的 API 地址
6. 在个人微信上收发消息

## 1. 在企业微信上创建企业

以手机企业微信为例，菜单依次如下：

我 | 设置 | 管理企业 | 全新创建企业 | 企业

- 填写企业名称
- 行业类型根据需要选择
- 员工规模选择1-50人
- 真实姓名填写自己的姓名

## 2. 绑定插件

电脑浏览器访问：

[https://work.weixin.qq.com/wework_admin/frame#profile/wxPlugin](https://work.weixin.qq.com/wework_admin/frame#profile/wxPlugin)

用个人微信扫【邀请关注】一栏的二维码。关注后，仅会绑定该企业，不会绑定其他企业。

## 3. 创建应用

电脑浏览器访问：

[https://work.weixin.qq.com/wework_admin/frame#apps](https://work.weixin.qq.com/wework_admin/frame#apps)

依次访问：

应用 | 自建 | 创建应用

填写页面要求的信息，记得在可见范围里把自己选上。

创建成功后，进入应用页面可以看到以下信息：

- AgentId
- Secret
- 功能 | 接收消息 | 设置API接收

其中 Secret 是企业的密码。在请求企业微信接口的时候，要使用企业 ID 和应用密码获取 Token。企业 ID 通过以下链接页面底部获取：

[https://work.weixin.qq.com/wework_admin/frame#profile](https://work.weixin.qq.com/wework_admin/frame#profile)

接着进入【设置API接收】，URL 按以下格式填写：

```
http://IP:端口/qiye-wechat/agents/应用的AgentId
```

Token 和 EncodingAESKey 随机获取或者自己填都行。

填完之后先不点保存，启动服务后再保存。

## 4. 启动服务

复制一份 config.json.example 到 config.json。

corp_id 是企业的 ID。上一步有说明获取方式。

在 agents 里面，填写刚刚创建的应用的信息。其中 secret 是通过应用界面的 Secret 一栏进入获取的。

type 是 agent 类型，与企业微信无关，与本项目 qiyewechat 文件夹里 agent.go 的 AgentFactory 有关。

```bash
go build
./noty
```

## 5. 配置企业微信

回到刚才的界面，点击保存。此时企业微信会发一条验证信息到服务器，如果通过，就能成功保存。

## 6. 个人微信收发消息

在个人微信的【我的企业及企业联系人】分组中，找到企业。进入后可以看到应用，发送消息即可。