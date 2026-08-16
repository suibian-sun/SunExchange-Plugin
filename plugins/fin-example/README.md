# FIN 跨平台插件示例

基于 [cat7street/FIN-plugin](https://github.com/cat7street/FIN-plugin) 的跨平台插件框架，
演示如何构建一个接入 SunExchange 鹊桥的服务器端插件。

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

## 说明

- 本示例展示玩家进服事件监听与游戏控制命令。
- 上层接入 SunExchange 鹊桥（WebSocket）进行进服验证 / 聊天互通 / 指令转发。
- 具体 SDK API 详见 FIN-plugin 的 `sdk/` 目录。

完整实现持续推进中，当前为占位骨架，供插件市场展示与后续开发对接。