# entry 示例插件（服务器端·进服验证）

演示 SunExchange 插件 SDK 的完整能力：

| 能力 | 说明 |
| --- | --- |
| 生命周期 | `Init / Start / Stop / GetInfo` |
| 事件监听 | 玩家进服（带优先级拦截）、聊天 |
| 命令注册 | 群聊/控制台命令 `sxadmin <玩家名>` 加白 |
| 游戏 API | 全服广播、私聊、执行命令、标题 |
| 插件间 API | 注册 `entry.whitelist` 供其他插件 `GetAPI` 调用 |
| 底层透传 | `Raw` 发送任意协议载荷 |

## 结构

```
examples/server/entry/
├── go.mod     # replace 指向 ../../../sdk
├── entry.go   # 插件实现（package entry）
└── sub/       # .so 导出包装（package main，导出 Plugin 符号）
    └── main.go
```

## 编译为 .so

```bash
cd examples/server/entry
go mod tidy
go build -buildmode=plugin -o entry.so ./sub
```

## 核心代码

```go
// 监听进服，最高优先级做白名单校验
ctx.ListenWithPriority("PlayerJoinEvent", func(ev any) bool {
    pe := ev.(sdk.PlayerEvent)
    if !p.whitelist[pe.Name] {
        ctx.Warnf("玩家 %s 未在白名单", pe.Name)
        return true // 拦截
    }
    return false
}, sdk.PriorityHigh)

// 注册插件间 API
ctx.RegisterAPI(&sdk.PluginAPI{Name: "entry.whitelist", Version: "1.0.0", Implement: p})
```

更多能力说明见 `docs/sdk-guide.md`。