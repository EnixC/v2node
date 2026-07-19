# v2node
A v2board backend base on moddified xray-core.
一个基于修改版xray内核的V2board节点服务端。

**注意： 本项目需要搭配[修改版V2board](https://github.com/wyx2685/v2board)**

## 软件安装

### Docker（推荐：一节点一容器）

每个容器只需要提供面板地址、节点 ID 和通讯密钥。日志级别、请求超时、重试次数、时区和 GeoData 路径均由镜像提供默认值。

```bash
docker run -d \
  --name v2node-1 \
  --restart unless-stopped \
  --network host \
  -e V2NODE_API_HOST='https://panel.example.com/' \
  -e V2NODE_NODE_ID='1' \
  -e V2NODE_API_KEY='your-api-key' \
  -v v2node-1:/etc/v2node \
  ghcr.io/enixc/v2node:latest
```

安装多个节点时，为每个节点使用不同的容器名和数据卷名称。使用 host 网络时，各节点在面板中配置的监听端口不能冲突。

可选环境变量：

| 变量 | 默认值 | 说明 |
| --- | --- | --- |
| `V2NODE_LOG_LEVEL` | `warning` | 日志级别：`debug`、`info`、`warn`、`warning`、`error` |
| `V2NODE_TIMEOUT` | `15` | 面板请求超时，单位为秒 |
| `V2NODE_RETRY_COUNT` | `1` | 面板请求重试次数 |
| `TZ` | `Asia/Shanghai` | 容器时区 |

仍可挂载传统的 `/etc/v2node/config.json` 启动；没有设置三个必填环境变量时，程序会自动使用文件配置。

### 原生一键安装

```
wget -N https://raw.githubusercontent.com/wyx2685/v2node/master/script/install.sh && bash install.sh
```

## 构建
``` bash
GOEXPERIMENT=jsonv2 go build -v -o build_assets/v2node -trimpath -ldflags "-X 'github.com/wyx2685/v2node/cmd.version=$version' -s -w -buildid="
```

## Stars 增长记录

[![Stargazers over time](https://starchart.cc/wyx2685/v2node.svg?variant=adaptive)](https://starchart.cc/wyx2685/v2node)
