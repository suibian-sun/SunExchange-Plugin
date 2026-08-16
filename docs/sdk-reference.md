# SunExchange 插件 SDK API 参考

`package sdk`（`github.com/suibian-sun/SunExchange-Plugin/sdk`）所有公开类型、常量与方法速查。纯标准库实现，Go 1.18+。

## 核心接口

### Plugin

插件必须实现的核心接口，由宿主驱动生命周期。

```go
type Plugin interface {
    Init(ctx *Context) error // 加载时调用，注入 Context，注册事件/命令/API
    Start() error            // 互通链路就绪后调用
    Stop() error             // 卸载/热重载/退出时调用
    GetInfo() PluginInfo     // 返回插件元信息
}
```

### PluginInfo

```go
type PluginInfo struct {
    Name        string // 唯一标识（如 sunexchange.entry）
    DisplayName string // 显示名
    Version     string // 版本
    Description string // 描述
    Author      string // 作者
}
```

## Context（插件唯一交互入口）

所有能力均通过 `sdk.Context` 暴露。`sdk.NewContext(opts Options)` 由宿主创建。

### Options

```go
type Options struct {
    ServerName  string                  // 鹊桥服务器名（握手身份）
    AccessToken string                  // 鹊桥接入 Token
    Broker      Broker                  // 插件间通信宿主
    Logger      Logger                  // 日志器
    Game        GameAPI                 // 游戏控制能力
    Send        func(api string, data any) (json.RawMessage, error) // 发请求到宿主
}
```

## 事件

### 监听

| 方法 | 说明 |
| --- | --- |
| `Listen(kind string, h Handler) func()` | 注册任意事件监听，返回取消函数 |
| `ListenWithPriority(kind string, h Handler, order int) func()` | 带优先级监听，order 越大越先执行 |
| `ListenPlayerJoin(h func(PlayerEvent) bool) func()` | 玩家进服 |
| `ListenPlayerLeave(h func(PlayerEvent) bool) func()` | 玩家离开 |
| `ListenChat(h func(*ChatEvent) bool) func()` | 玩家聊天（可取消） |

`Handler` 签名：`func(ev any) bool`，返回 `true` 表示消费事件 / 阻止后续监听器。

### 优先级常量

```go
const (
    PriorityLowest  = iota // 0
    PriorityLow           // 1
    PriorityNormal        // 2
    PriorityHigh          // 3
    PriorityHighest       // 4
)
```

### 事件类型

| 类型 | 事件名 | 结构 |
| --- | --- | --- |
| `PlayerEvent` | `PlayerJoinEvent` / `PlayerQuitEvent` | `Name, UUID, IsOp, ServerName, ServerType, ServerVer, Raw, EntryIndex` |
| `ChatEvent` | `PlayerChatEvent` | `Sender, Message, RawMessage, ServerName, ServerType, ServerVer, Raw, Cancelled` |

未识别的 `event_name` 以 `map[string]any` 透传（高自由度扩展口）。

## 命令

```go
type Command struct {
    Name        string              // 命令名（唯一）
    Triggers    []string            // 触发别名（可选，自动含 Name）
    ArgHint     string              // 参数提示，如 "<玩家> <数量>"
    Usage       string              // 用法说明
    Description string              // 描述
    Handler     func(args []string) error // 处理回调
}
```

`RegisterCommand(cmd Command) error`：注册命令，`Name` 不能为空。

## 游戏 API（GameAPI）

| 方法 | 说明 |
| --- | --- |
| `SendChat(msg any) error` | 全服广播 |
| `SendWhisper(target, msg any) error` | 向指定玩家私聊 |
| `SendCommand(cmd string) (string, error)` | 执行服务端命令，返回输出 |
| `OnlineNames() []string` | 在线玩家名 |
| `SendTitle(title, subtitle any) error` | 发送标题 |
| `SendActionbar(msg any) error` | 发送动作栏 |

通过 `ctx.Game()` 获取。

## 插件间 API（Broker）

```go
type PluginAPI struct {
    Name      string // API 名（全局唯一）
    Version   string // 语义化版本
    Implement any    // 实现对象，宿主通过反射调用
}
```

| 方法 | 说明 |
| --- | --- |
| `RegisterAPI(api *PluginAPI) error` | 注册插件 API |
| `GetAPI(name, version string) (any, error)` | 获取插件 API（版本兼容校验） |

版本兼容规则：major 必须相等，请求 minor ≤ 实现 minor。

## 广播

```go
type Broadcast struct {
    Name string // 主题
    Data any    // 载荷
}
type BroadcastHandler func(b Broadcast)
```

| 方法 | 说明 |
| --- | --- |
| `Subscribe(topic string, h BroadcastHandler)` | 订阅广播主题 |
| `Publish(b Broadcast) []any` | 发布广播，返回订阅者返回值 |

## 底层透传

```go
func (c *Context) Raw(kind string, data any) error
```

向宿主发送任意 `kind + data`，用于扩展协议 / 对接自有系统。

## 日志

`Context` 直接提供便捷日志方法（无独立字段）：

```go
ctx.Debugf(format string, a ...any)
ctx.Infof(format string, a ...any)
ctx.Successf(format string, a ...any)
ctx.Warnf(format string, a ...any)
ctx.Errorf(format string, a ...any)
```

`Logger` 接口（宿主实现）：`Debugf / Infof / Successf / Warnf / Errorf`。

## 工具函数

```go
func WaitFor(timeout time.Duration, cond func() bool, tick time.Duration) error
```

等待条件满足，超时返回 `ErrTimeout`，用于状态机式交互。

## 错误值

```go
var (
    ErrEmptyCommandName = errors.New("命令名不能为空")
    ErrNoBroker         = errors.New("宿主未提供插件间通信能力")
    ErrNoTransport      = errors.New("宿主未提供底层透传通道")
    ErrTimeout          = errors.New("等待超时")
)
```