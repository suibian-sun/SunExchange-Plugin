# FIN 跨平台插件示例

基于 [maoqijie/FIN-plugin](https://github.com/maoqijie/FIN-plugin) 的跨平台插件框架，演示如何构建一个接入 SunExchange 鹊桥的服务器端进服验证插件。

## 结构

```
plugins/fin-example/
├── go.mod
├── main.go
└── README.md
```

## 构建

```bash
cd plugins/fin-example
go build -o sunexchange-entry .
```

产物 `sunexchange-entry` 即为可运行的服务器端插件（FIN 使用 gRPC 与宿主进程通信，跨平台：Windows/Linux/macOS/Android）。

顶层 `registry/extensions.json` 中 `sunexchange.entry` 的 `asset_url` 指向由本示例构建的 Release 资产。

## 插件要点

- 实现 `Init / Start / Stop / GetInfo` 四个生命周期方法。
- `Init` 中 `ListenPlayerJoin` 触发进服验证，`ListenChat` 转发聊天到鹊桥（群服互通）。
- `GameUtils` 提供游戏控制：`SendChat` 发消息、`SendCommand` 执行服务端命令。
- 通过 `plugin.Serve` 暴露 `sdk.PluginGRPC`，由 FIN 宿主进程加载。

详见 [插件制作完整指南](../../docs/plugin-development-guide.md) 与 [服务器端插件制作指南](../../docs/server-plugin-guide.md)。