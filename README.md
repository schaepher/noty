# noty

以企业微信为代理，打通个人微信和自建服务器的双向通信通道。

个人注册企业时，不会要求验证企业信息。

# 实现后的消息流向

个人微信 -> 企业微信 -> 个人服务器  

个人微信 <- 企业微信 <- 个人服务器

# 企业微信配置

见 [README_config_qiyewechat.md](./README_config_qiyewechat.md)

# PDF 插件

如果发送一个链接，会让服务器将网页内容转成 PDF 保存。

需要先安装和启动 doctron 服务：https://github.com/lampnick/doctron

它支持把微信公众号文章保存为 PDF，包含所有图片。

该选项的配置示例：

```json
{
    "agents": [
        {
            "pdf_convert": {
                "url": "http://127.0.0.1:8080/convert/html2pdf",
                "username": "",
                "password": "",
                "pdf_dir": ""
            }
        }
    ]
}
```