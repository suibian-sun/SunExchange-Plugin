# SunExchange 插件 SDK 开发指南

本指南介绍如何在 SunExchange 生态内编写插件（进服验证、群服互通、经济联动、自动管理等）。SDK 借鉴 FIN（FunInterWork）插件框架的高自由度设计，但通信基础完全复用 SunExchange 自有的**鹊桥（Queqiao）协议**，无需额外基础设施。插件以独立进程运行，通过 WebSocket 与宿主双向通信。

> 适用对象：想在 SunExchange 生态内编写插件能力的开发者。SDK 源码位于仓库 `sdk/` 目录，独立 Go module，可被任意插件项目直接引用。

## 一、体系总览

```
┌─────────────────────── SunExchange 主进程 ───────────────────────┐
│  Web 管理端 · 任务工作进程 · 鹊桥 WS 服务(每任务一个)             │
└────────────────────────────────────┬─────────────────────────────┘
                                     │ ws://<host>:<port>/minecraft/ws
┌────────────────────────────────────▼─────────────────────────────┐
│  插件宿主（SunExchange 主仓 plugin/host）                         │
│  ├── Bridge  连接鹊桥，转发游戏事件 / 回写命令                    │
│  ├── Host    管理插件生命周期（Init → Start → Stop），分发事件   │
│  └── Loader  动态加载 .so 插件（免重编译安装）                   │
└────────────────────────────────────┬─────────────────────────────┘
                                     │ 调用 sdk.Plugin
┌────────────────────────────────────▼─────────────────────────────┐
│  你的插件（实现 sdk.Plugin）                                     │
│  事件监听 · 命令 · 游戏 API · 插件间 API · 广播 · 底层透传       │
└───────────────────────────────────────────────────────────────────┘
```

核心概念：**插件 = 一个实现 `sdk.Plugin` 接口的独立进程**，宿主负责加载它并注入一个 `sdk.Context`。所有能力都通过 `Context` 暴露。

## 二、快速开始

### 1. 引入 SDK

SDK 为独立模块（`github.com/suibian-sun/SunExchange-Plugin/sdk`），纯标准库、低 Go 版本（1.18+）即可编译：

```bash
module myplugin
go 1.18
require github.com/suibian-sun/SunExchange-Plugin/sdk v1.0.0
```

> 开发期若 SDK 尚未发版，可在 `go.mod` 加 `replace` 指向本地路径：
> `replace github.com/suibian-sun/SunExchange-Plugin/sdk => /path/to/SunExchange-Plugin/sdk`

### 2. 实现 Plugin 接口

```go
package main

import "github.com/suibian-sun/SunExchange-Plugin/sdk"

type MyPlugin struct{ ctx *sdk.Context }

func (p *MyPlugin) GetInfo() sdk.PluginInfo {
    return sdk.PluginInfo{
        Name: "myplugin", DisplayName: "我的插件",
        Version: "1.0.0", Description: "示例", Author: "你",
    }
}

func (p *MyPlugin) Init(ctx *sdk.Context) error {
    p.ctx = ctx
    // ...注册事件、命令、API
    return nil
}

func (p *MyPlugin) Start() error { return nil }
func (p *MyPlugin) Stop() error  { return nil }

func Plugin() sdk.Plugin { return &MyPlugin{} } // .so 导出符号
func main() {}
```

### 3. 编译为 .so

```bash
go build -buildmode=plugin -o myplugin.so .
```

### 4. 用宿主加载

```bash
# 连接某任务的鹊桥（端口取任务 WorkerPort），动态加载插件
sunexchange-plugin \
  -ws ws://127.0.0.1:23000/minecraft/ws \
  -server task-1-demo \
  -plugins ./myplugin.so
```

宿主二进制位于 SunExchange 主仓 `plugin/cmd/sunexchange-plugin`。

## 三、生命周期

宿主驱动三段生命周期，插件在其中建立 / 释放资源：

| 阶段 | 时机 | 典型用途 |
| --- | --- | --- |
| `Init(ctx)` | 插件加载时 | 注册事件、命令、插件 API；读配置 |
| `Start()` | 互通链路就绪后 | 启动后台协程、发送就绪广播 |
| `Stop()` | 卸载 / 热重载 / 退出 | 关协程、释放连接、保存状态 |

## 四、高自由度能力

### 4.1 事件监听与拦截

可监听任意事件，用**优先级**控制执行顺序。返回 `true` 表示消费事件、阻止后续监听器（拦截链）。

```go
// 带优先级：玩家进服先做白名单校验
ctx.ListenWithPriority("PlayerJoinEvent", func(ev any) bool {
    pe := ev.(sdk.PlayerEvent) // 宿主已解码为强类型
    if !p.whitelist[pe.Name] {
        ctx.Warnf("玩家 %s 未在白名单", pe.Name)
        return true // 拦截：阻止后续逻辑
    }
    return false
}, sdk.PriorityHigh)

// 语法糖：聊天事件（可取消）
ctx.ListenChat(func(c *sdk.ChatEvent) bool {
    if c.Message == "/mute" { c.Cancelled = true; return true }
    return false
})
```

内置强类型事件：`PlayerJoinEvent` / `PlayerQuitEvent` → `sdk.PlayerEvent`，`PlayerChatEvent` → `*sdk.ChatEvent`。

### 4.2 完全自由的原始事件

未识别的 `event_name` 会以原始 `map[string]any` 透传给插件。这让插件可以处理**任何**协议事件，无需 SDK 预定义——这是"高自由度"的关键。

```go
ctx.Listen("MyCustomEvent", func(ev any) bool {
    m := ev.(map[string]any)
    ctx.Infof("收到自定义事件: %v", m)
    return false
})
```

### 4.3 游戏 API

```go
ctx.Game().SendChat("§a全服公告")              // 全服广播
ctx.Game().SendWhisper("Steve", "私聊")       // 私聊
ctx.Game().SendCommand("/give @a diamond 64") // 执行命令
ctx.Game().SendTitle("标题", "副标题")
ctx.Game().OnlineNames()
```

### 4.4 命令注册

```go
ctx.RegisterCommand(sdk.Command{
    Name: "sxadmin", Triggers: []string{"sx"},
    ArgHint: "<玩家名>", Usage: "sxadmin <玩家名> - 加白",
    Handler: func(args []string) error {
        p.whitelist[args[0]] = true
        return nil
    },
})
```

### 4.5 插件间 API（跨插件能力共享）

注册自己的能力，供其他插件按**语义化版本**调用；宿主做 major 一致、请求 minor ≤ 实现 minor 的兼容校验。

```go
// 注册
ctx.RegisterAPI(&sdk.PluginAPI{Name: "entry.whitelist", Version: "1.0.0", Implement: p})

// 其他插件获取并调用
api, err := ctx.GetAPI("entry.whitelist", "1.0.0")
if e, ok := api.(*entry.Entry); ok { ok := e.CheckWhitelist("Steve") }
```

### 4.6 广播（发布 / 订阅）

```go
ctx.Subscribe("state.changed", func(b sdk.Broadcast) { /* 响应 */ })
ctx.Publish(sdk.Broadcast{Name: "state.changed", Data: "online=10"})
```

### 4.7 底层透传（任意协议）

直接向宿主发送任意 `kind + data`，可用于扩展协议、对接自有系统：

```go
ctx.Raw("hello", map[string]any{"from": "myplugin"})
```

### 4.8 状态机等待

```go
err := sdk.WaitFor(10*time.Second, func() bool { return p.ready }, 100*time.Millisecond)
```

## 五、用脚手架生成

推荐用 CLI 生成插件骨架，自动带好 go.mod 与正确依赖：

```bash
# 安装或直接 go run
go run github.com/suibian-sun/SunExchange-Plugin/cli new myplugin --end server
```

详见 [脚手架 CLI 指南](cli-guide.md)。

## 六、测试

示例插件自带完整能力，可 `go vet ./... && go build ./...` 校验。宿主（主仓）自带单元测试覆盖 `.so` 加载、事件解码分发、插件 API 与版本校验。

## 七、最佳实践

- **用优先级做拦截链**：需要"先校验后放行"的插件用 `PriorityHigh`，避免被后加载插件抢跑。
- **State 放 `Implement`**：插件间 API 暴露的是对象，把读方法写在暴露对象上，调用方断言后使用。
- **清理要彻底**：`Stop` 里关闭自己开的连接与协程，宿主热重载依赖它。
- **事件名要稳定**：`event_name` 是插件间契约，对外发布后别随意改名。
- **自由度是把双刃剑**：`Raw` 与原始事件透传很强，但尽量用内置强类型事件，保持可读性与可维护性。

## 八、常见问题

**Q：加载 .so 报"requires exactly one main package"？**
A：`.so` 必须由 `package main` 编译，并包含 `func Plugin() sdk.Plugin` 与 `func main(){}`。参见 `examples/server/entry/sub/main.go`。

**Q：事件处理器里 `ev.(sdk.PlayerEvent)` panic？**
A：宿主已把 `PlayerJoinEvent` 解码为 `sdk.PlayerEvent`。若用 `Listen`（非语法糖）监听自定义事件，请先断言为 `map[string]any` 再取值。

**Q：如何连接某个任务的鹊桥？**
A：鹊桥端口是任务的 `WorkerPort`，从任务详情获取，宿主用 `-ws ws://<host>:<port>/minecraft/ws -server <task server name>` 连接。