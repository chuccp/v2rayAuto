\# v2rayAuto\
生成订阅地址，自动刷新v2ray 端口\
\
直接引用 github.com/v2fly/v2ray-core 作为开发包使用，只需要单个可执行文件，就可科学上网,\
降低部署难度

\
实现：证书申请，订阅地址生成，端口每天随机变更等功能\
\
目前实现 websocket+tls+vmess。使用前需要申请域名\
\
有问题 请加入 <https://t.me/+POaWrs8YGFQ4MDdl>\
下载发布包至服务器解压\
运行\
v2rayAuto start\
停止\
v2rayAuto stop\
使用命令窗口运行\
v2rayAuto run\
配置文件 config.ini 说明\
`[core]`\
`##随机端口范围`\
`range_port = 8390-9000`\
`##要访问的域名`\
`host = 127.0.0.1`\
`[tls]`\
`##申请证书的域名`\
`domain = 127.0.0.1`\
`##申请证书的邮件`\
`email = `[`cooge123@gmail.com`](mailto\:cooge123@gmail.com)\
`[web]`\
`##订阅端口，不要与其它端口冲突`\
`port = 8336`\
`##订阅地址路径 例子 `[`https://127.0.0.1:8336/abc.md`](https://127.0.0.1:8336/abc.md)\
`subscribe = /abc.md`\
`[vmess_ws]`\
`##创建链接数量`\
`create_num = 5`