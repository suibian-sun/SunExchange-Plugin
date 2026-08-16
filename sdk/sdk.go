// Package sdk 提供 SunExchange 插件开发接口。
//
// 插件以独立进程运行，通过鹊桥（queqiao）WebSocket 协议与 SunExchange
// 主进程双向通信。插件实现 Plugin 接口，由宿主加载并驱动生命周期。
//
// 本包为纯标准库实现、无外部依赖，可被任意插件模块直接引用。
// 设计借鉴 FIN（FunInterWork）插件框架的高自由度特性，但通信基础复用
// SunExchange 自有的鹊桥协议，无需引入额外基础设施。
package sdk

import (
	"encoding/json"
	"time"
)

// PluginInfo 插件元信息，由 GetInfo 返回并在宿主侧展示。
type PluginInfo struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Author      string `json:"author"`
}

// Plugin 是插件必须实现的核心接口。
// 生命周期由宿主驱动：Init → Start →（运行中）→ Stop。
type Plugin interface {
	// Init 在插件加载时调用，注入 Context。在此注册事件监听、命令与插件 API。
	Init(ctx *Context) error
	// Start 在互通链路就绪后调用，启动正式逻辑。
	Start() error
	// Stop 在卸载 / 热重载 / 退出时调用，用于释放资源。
	Stop() error
	// GetInfo 返回插件元信息。
	GetInfo() PluginInfo
}

// PlayerEvent 玩家事件（进服 / 离开）。
type PlayerEvent struct {
	Name         string `json:"name"`
	UUID         string `json:"uuid"`
	IsOp         bool   `json:"is_op"`
	ServerName   string `json:"server_name"`
	ServerType   string `json:"server_type"`
	ServerVer    string `json:"server_version"`
	Raw          any    `json:"raw,omitempty"`
	EntryIndex   int    `json:"entry_index"`
}

// ChatEvent 玩家聊天事件。Cancelled 置 true 可拦截消息，阻止其继续传播。
type ChatEvent struct {
	Sender    string `json:"sender"`
	Message   string `json:"message"`
	RawMessage string `json:"raw_message"`
	ServerName string `json:"server_name"`
	ServerType string `json:"server_type"`
	ServerVer  string `json:"server_version"`
	Raw       any    `json:"raw,omitempty"`
	Cancelled bool   `json:"-"`
}

// Command 控制台 / 群聊命令注册。
type Command struct {
	// Name 命令名（唯一）。
	Name string `json:"name"`
	// Triggers 触发别名（可选），自动含 Name。
	Triggers []string `json:"triggers,omitempty"`
	// ArgHint 参数提示，如 "<玩家> <数量>"。
	ArgHint string `json:"arg_hint,omitempty"`
	// Usage 用法说明。
	Usage string `json:"usage,omitempty"`
	// Description 命令描述。
	Description string `json:"description,omitempty"`
	// Handler 命令处理回调。
	Handler func(args []string) error `json:"-"`
}

// Broadcast 插件间广播（订阅 / 发布）。Name 标识主题，Data 为载荷。
type Broadcast struct {
	Name string `json:"name"`
	Data any    `json:"data,omitempty"`
}

// 事件优先级：数值越大越先执行，用于拦截链。
const (
	PriorityLowest = iota
	PriorityLow
	PriorityNormal
	PriorityHigh
	PriorityHighest
)

// OnEvent 事件回调接口。所有事件类型均可通过类型断言分发。
// 返回 true 表示已消费（阻止后续监听器）；false 继续传播。
type OnEvent func(ev any) bool

// Handler 事件处理器：接收具体事件类型，返回是否消费。
type Handler func(ev any) bool

// rawMessage 底层透传消息（高自由度扩展口）。
type rawMessage struct {
	Kind string          `json:"kind"`
	Data json.RawMessage `json:"data"`
}

// Context 是插件与宿主的唯一交互入口，封装事件、命令、游戏 API 与日志。
// 所有方法均为并发安全；Stop 后自动注销监听与等待器。
type Context struct {
	// 事件监听（注册后由宿主驱动）
	onEvent map[string][]handlerEntry
	// 命令注册
	commands map[string]*Command
	// 插件 API 注册
	apis map[string]*PluginAPI
	// 广播订阅
	subs map[string][]BroadcastHandler
	// 底层透传
	raw chan rawMessage

	// 宿主能力（由宿主注入）
	game   GameAPI
	logger Logger
	sendc  func(api string, data any) (json.RawMessage, error)
	broker Broker
}

// handlerEntry 事件处理器条目（含优先级）。
type handlerEntry struct {
	handler Handler
	order   int
}

// BroadcastHandler 广播处理回调。
type BroadcastHandler func(b Broadcast)

// PluginAPI 插件间 API 描述。
type PluginAPI struct {
	// Name API 名（全局唯一，如 "economy"）。
	Name string `json:"name"`
	// Version 语义化版本，供版本兼容校验。
	Version string `json:"version"`
	// Implement 实现对象（任意 Go 值），宿主通过反射调用。
	Implement any `json:"-"`
}

// Logger 插件日志接口。
type Logger interface {
	Debugf(format string, args ...any)
	Infof(format string, args ...any)
	Successf(format string, args ...any)
	Warnf(format string, args ...any)
	Errorf(format string, args ...any)
}

// Broker 宿主提供的插件间通信能力。
type Broker interface {
	// TriggerBroadcast 触发广播，返回订阅者返回值。
	TriggerBroadcast(b Broadcast) []any
	// RegisterAPI 注册插件 API。
	RegisterAPI(api *PluginAPI) error
	// GetAPI 获取插件 API（按名 + 版本兼容）。
	GetAPI(name string, version string) (any, error)
}

// GameAPI 游戏控制能力（由宿主桥接到底层平台）。
type GameAPI interface {
	// SendChat 全服广播消息。
	SendChat(msg any) error
	// SendWhisper 向指定玩家发送私聊。
	SendWhisper(target, msg any) error
	// SendCommand 执行服务端命令，返回输出。
	SendCommand(cmd string) (string, error)
	// OnlineNames 返回在线玩家名列表。
	OnlineNames() []string
	// SendTitle 发送标题。
	SendTitle(title, subtitle any) error
	// SendActionbar 发送动作栏。
	SendActionbar(msg any) error
}

// Options 插件初始化选项。
type Options struct {
	// ServerName 鹊桥服务器名（握手身份）。
	ServerName string
	// AccessToken 鹊桥接入 Token。
	AccessToken string
	// Broker 插件间通信宿主。
	Broker Broker
	// Logger 日志器。
	Logger Logger
	// Game 游戏控制能力。
	Game GameAPI
	// Send 发请求到宿主（api + 数据），返回原始响应。
	Send func(api string, data any) (json.RawMessage, error)
}

// NewContext 由宿主创建插件 Context。
func NewContext(opts Options) *Context {
	c := &Context{
		onEvent:  make(map[string][]handlerEntry),
		commands: make(map[string]*Command),
		apis:     make(map[string]*PluginAPI),
		subs:     make(map[string][]BroadcastHandler),
		raw:      make(chan rawMessage, 64),
		game:     opts.Game,
		logger:   opts.Logger,
		sendc:    opts.Send,
		broker:   opts.Broker,
	}
	return c
}

// ---- 事件监听 ----

// Listen 注册任意事件监听器。kind 为事件类型名（如 "PlayerJoinEvent"）。
// 返回取消函数。
func (c *Context) Listen(kind string, h Handler) func() {
	return c.listen(kind, h, PriorityNormal)
}

// ListenWithPriority 注册带优先级的事件监听器。order 越大越先执行。
func (c *Context) ListenWithPriority(kind string, h Handler, order int) func() {
	return c.listen(kind, h, order)
}

func (c *Context) listen(kind string, h Handler, order int) func() {
	key := normalizeKind(kind)
	c.onEvent[key] = append(c.onEvent[key], handlerEntry{handler: h, order: order})
	// 降序：order 越大越先执行
	sortHandlers(c.onEvent[key])
	return func() { c.removeHandler(key, h) }
}

func sortHandlers(entries []handlerEntry) {
	for i := 1; i < len(entries); i++ {
		for j := i; j > 0 && entries[j-1].order < entries[j].order; j-- {
			entries[j-1], entries[j] = entries[j], entries[j-1]
		}
	}
}

func (c *Context) removeHandler(kind string, h Handler) {
	entries := c.onEvent[kind]
	for i, e := range entries {
		if equalHandler(e.handler, h) {
			c.onEvent[kind] = append(entries[:i], entries[i+1:]...)
			return
		}
	}
}

// ListenPlayerJoin 玩家进服。
func (c *Context) ListenPlayerJoin(h func(PlayerEvent) bool) func() {
	return c.Listen("PlayerJoinEvent", func(ev any) bool { return h(ev.(PlayerEvent)) })
}

// ListenPlayerLeave 玩家离开。
func (c *Context) ListenPlayerLeave(h func(PlayerEvent) bool) func() {
	return c.Listen("PlayerQuitEvent", func(ev any) bool { return h(ev.(PlayerEvent)) })
}

// ListenChat 玩家聊天。
func (c *Context) ListenChat(h func(*ChatEvent) bool) func() {
	return c.Listen("PlayerChatEvent", func(ev any) bool { return h(ev.(*ChatEvent)) })
}

// ---- 命令 ----

// RegisterCommand 注册命令。
func (c *Context) RegisterCommand(cmd Command) error {
	if cmd.Name == "" {
		return ErrEmptyCommandName
	}
	c.commands[cmd.Name] = &cmd
	return nil
}

// ---- 游戏 API ----

// Game 返回游戏控制能力。
func (c *Context) Game() GameAPI { return c.game }

// ---- 插件 API ----

// RegisterAPI 注册插件间 API。
func (c *Context) RegisterAPI(api *PluginAPI) error {
	if c.broker != nil {
		return c.broker.RegisterAPI(api)
	}
	return ErrNoBroker
}

// GetAPI 获取插件 API。
func (c *Context) GetAPI(name, version string) (any, error) {
	if c.broker != nil {
		return c.broker.GetAPI(name, version)
	}
	return nil, ErrNoBroker
}

// ---- 广播 ----

// Subscribe 订阅广播主题。
func (c *Context) Subscribe(topic string, h BroadcastHandler) {
	c.subs[topic] = append(c.subs[topic], h)
}

// Publish 发布广播，返回订阅者返回值。
func (c *Context) Publish(b Broadcast) []any {
	if c.broker == nil {
		return nil
	}
	return c.broker.TriggerBroadcast(b)
}

// ---- 底层透传 ----

// Raw 发送底层原始消息（高自由度：可透传任意协议载荷）。
func (c *Context) Raw(kind string, data any) error {
	if c.sendc == nil {
		return ErrNoTransport
	}
	_, err := c.sendc("plugin.raw", map[string]any{"kind": kind, "data": data})
	return err
}

// ---- 日志 ----

func (c *Context) Debugf(f string, a ...any) { if c.logger != nil { c.logger.Debugf(f, a...) } }
func (c *Context) Infof(f string, a ...any)  { if c.logger != nil { c.logger.Infof(f, a...) } }
func (c *Context) Successf(f string, a ...any) { if c.logger != nil { c.logger.Successf(f, a...) } }
func (c *Context) Warnf(f string, a ...any)  { if c.logger != nil { c.logger.Warnf(f, a...) } }
func (c *Context) Errorf(f string, a ...any) { if c.logger != nil { c.logger.Errorf(f, a...) } }

// ---- 内部 ----

// Dispatch 由宿主调用，向插件分发事件。返回是否被消费。
func (c *Context) Dispatch(kind string, ev any) bool {
	key := normalizeKind(kind)
	entries := c.onEvent[key]
	for _, e := range entries {
		if e.handler(ev) {
			return true
		}
	}
	return false
}

func normalizeKind(k string) string {
	// 兼容带包名 / 无前缀的事件名
	if i := lastIndexByte(k, '.'); i >= 0 {
		k = k[i+1:]
	}
	return k
}

func lastIndexByte(s string, b byte) int {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == b {
			return i
		}
	}
	return -1
}

func equalHandler(a, b Handler) bool {
	return fmtPointer(a) == fmtPointer(b)
}

func fmtPointer(h Handler) uintptr {
	// 通过反射取函数指针，用于去重比较
	return reflectValueOf(h)
}

// WaitFor 等待一个条件满足（超时返回 error）。用于状态机式交互。
func WaitFor(timeout time.Duration, cond func() bool, tick time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if cond() {
			return nil
		}
		time.Sleep(tick)
	}
	return ErrTimeout
}