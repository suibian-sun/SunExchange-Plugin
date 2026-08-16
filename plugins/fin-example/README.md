# FIN 跨平台插件示例

基于 [cat7street/FIN-plugin](https://github.com/cat7street/FIN-plugin) 的跨平台插件框架，演示如何构建一个接入 SunExchange 鹊桥的服务器端插件。

## 结构

```
plugins/fin-example/
├── go.mod
└── main.go
```

## 构建

```bash
cd plugins/fin-example
go build -o sunexchange-entry .
```

## main.go 骨架

```go
package main

import (
    "github.com/cat7street/FIN-plugin/sdk"
    "github.com/cat7street/FIN-plugin/sdk/plugin"
)

type Entry struct{}

func (p *Entry) OnLoad(ctx *sdk.Context) error {
    // 玩家进服：触发进服验证与事件回写
    ctx.ListenPlayerJoin(func(event sdk.PlayerEvent) {
        ctx.LogSuccess("玩家 %s 加入，触发进服验证", event.Name)
        // 调用 SunExchange 校验接口并回写事件到鹊桥
    })

    // 玩家聊天：转发到鹊桥（群服互通）
    ctx.ListenChat(func(event *sdk.ChatEvent) {
        ctx.LogInfo("聊天 %s: %s", event.Sender, event.Message)
        // bridge.EmitChat(event.Sender, event.Message)
    })

    // 游戏控制
    gu := ctx.GameUtils()
    gu.Broadcast("§eSunExchange 进服验证已就绪")

    return nil
}

func main() { plugin.Serve(&Entry{}) }
```

## 说明

- 本示例展示玩家进服事件监听与游戏控制命令。
- 上层接入 SunExchange 鹊桥（WebSocket）进行进服验证 / 聊天互通 / 指令转发。
- 具体 SDK API 详见 FIN-plugin 的 `sdk/` 目录。

完整方案见 [插件制作完整指南](../../docs/plugin-development-guide.md) 与 [服务器端插件制作指南](../../docs/server-plugin-guide.md)。