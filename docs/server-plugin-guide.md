# 服务器端插件制作指南

本指南说明如何为 SunExchange 统一插件市场制作 **Minecraft 服务器端插件**。服务器端插件参考 [maoqijie/FIN-plugin](https://github.com/maoqijie/FIN-plugin) 的跨平台插件框架设计，同时支持经典 Bukkit 服务端插件与网易 FunShuttler 脚本。

## 一、插件分类与适用端

| 适用端 | 平台标识 | 服务端类型 | 插件形态 |
| --- | --- | --- | --- |
| 网易我的世界 | `netease` | 网易租赁服 | 网易 funcs 插件 / 脚本 |
| 国际 Java 版 | `java` | Spigot / Paper / Purpur（Bukkit） | `.jar` |
| 国际基岩版 | `bedrock` | BDS / Nukkit / PocketMine | `.jar` / `.js` / 插件文件夹 |

在注册表登记时，`kind` 填 `server`，`platforms` 只列你实际支持的端。

## 二、插件的核心职责（群服互通）

服务器端插件在游戏内完成三件事，与 SunExchange 的鹊桥（queqiao）协议对接：

1. **进服验证**：玩家加入时，向 SunExchange 校验该玩家是否已购买并拥有有效卡槽/任务。
2. **聊天互通**：把玩家聊天转发为 `PlayerChatEvent`，并接收来自 QQ 群的指令向玩家/服务器广播。
3. **指令与事件转发**：接收 `SendCommand` / `SendWhisper`，执行服务端命令并回传结果。

## 三、鹊桥 WebSocket 协议

服务器端插件通过 WebSocket 连接 SunExchange 的鹊桥服务，地址为 `ws://<sunexchange>:<port>/minecraft/ws`。连接建立后用插件配置声明服务器身份。

### 上行事件（游戏 → 群）

统一的事件结构（字段以 SunExchange 鹊桥实现为准）：

```json
{
  "timestamp": 1760000000000,
  "post_type": "notice",
  "event_name": "PlayerJoinEvent",
  "server_name": "my-server",
  "server_type": "java",
  "server_version": "1.20.1",
  "sub_type": "player_join",
  "player": { "nickname": "Steve", "uuid": "uuid-here", "is_op": false }
}
```

| event_name | post_type | sub_type | 说明 |
| --- | --- | --- | --- |
| `PlayerJoinEvent` | notice | `player_join` | 玩家进服 |
| `PlayerQuitEvent` | notice | `player_quit` | 玩家离开 |
| `PlayerChatEvent` | message | `player_chat` | 玩家聊天 |

玩家聊天示例：

```json
{
  "post_type": "message",
  "event_name": "PlayerChatEvent",
  "server_name": "my-server",
  "server_type": "java",
  "server_version": "1.20.1",
  "sub_type": "player_chat",
  "message_id": "uuid",
  "raw_message": "hello",
  "message": "hello",
  "player": { "nickname": "Steve", "uuid": "uuid-here", "is_op": false }
}
```

### 下行事件（群 → 游戏）

```json
{ "event_name": "SendCommandEvent", "command": "/give @a diamond 64" }
{ "event_name": "SendWhisperEvent", "player": "Steve", "message": "你好" }
```

插件收到后执行命令 / 私聊，并把结果按 `echo` 回包。

### API 请求 / 响应

适配器可向服务器端发起带 `echo` 的 API 请求，服务器端做出响应：

```json
// 请求
{ "api": "get_online_players", "data": {}, "echo": "abc123" }
// 响应
{ "code": 200, "api": "get_online_players", "post_type": "response",
  "status": "SUCCESS", "message": "success", "data": [], "echo": "abc123" }
```

## 四、使用 FIN 跨平台框架（推荐）

FIN 插件框架（`maoqijie/FIN-plugin`）支持跨平台插件系统（Windows/Linux/macOS/Android），一套 Go 代码跑通三个适用端，进程隔离、跨端一致。

### 4.1 生命周期与事件监听

```go
package main

import (
    "github.com/hashicorp/go-plugin"
    "github.com/maoqijie/FIN-plugin/sdk"
)

type MyPlugin struct{ ctx *sdk.Context }

// Init 是插件入口：注册命令、监听事件均在此时完成
func (p *MyPlugin) Init(ctx *sdk.Context) error {
    p.ctx = ctx
    // 注册控制台命令
    ctx.RegisterConsoleCommand(sdk.ConsoleCommand{
        Name: "mycmd",
        Handler: func(args []string) error {
            ctx.LogInfo("命令执行")
            return nil
        },
    })
    // 监听玩家加入
    ctx.ListenPlayerJoin(func(event sdk.PlayerEvent) {
        ctx.LogSuccess("玩家 %s 加入", event.Name)
    })
    // 监听聊天（*ChatEvent 可置 Cancelled=true 拦截转发）
    ctx.ListenChat(func(event *sdk.ChatEvent) {
        ctx.LogInfo("%s: %s", event.Sender, event.Message)
    })
    return nil
}

func (p *MyPlugin) Start() error { return nil }
func (p *MyPlugin) Stop() error  { return nil }
func (p *MyPlugin) GetInfo() sdk.PluginInfo {
    return sdk.PluginInfo{Name: "MyPlugin", DisplayName: "我的插件", Version: "1.0.0", Description: "", Author: ""}
}

func main() {
    plugin.Serve(&plugin.ServeConfig{
        HandshakeConfig: sdk.HandshakeConfig,
        Plugins: map[string]plugin.Plugin{
            "plugin": &sdk.PluginGRPC{Impl: &MyPlugin{}},
        },
        GRPCServer: plugin.DefaultGRPCServer,
    })
}
```

### 4.2 游戏控制

```go
gu := ctx.GameUtils()
gu.SayTo("玩家名", "§a你好！")
gu.SendChat("§e全服公告")
gu.SendCommand("/give @a diamond 64")
```

### 4.3 控制台命令与配置

```go
ctx.RegisterConsoleCommand(sdk.ConsoleCommand{
    Name: "mycmd",
    Handler: func(args []string) error { ctx.LogInfo("执行"); return nil },
})

config := ctx.Config()
config.SetDefault("key", "value")
value := config.GetString("key")
```

参考本仓库 `plugin/server/fin-example/` 的完整示例。

## 五、经典服务端插件（Java/Bukkit）

不依赖 FIN 框架时，直接写标准 Bukkit 插件，核心是 `onEnable` 注册监听器，再把事件转发到鹊桥 HTTP/WebSocket 通道：

```java
public final class EntryPlugin extends JavaPlugin implements Listener {
    private BridgeClient bridge;

    @Override public void onEnable() {
        getServer().getPluginManager().registerEvents(this, this);
        bridge = new BridgeClient(getConfig().getString("bridge.ws"));
        bridge.connect();
    }

    @Override public void onDisable() { bridge.disconnect(); }

    @EventHandler public void onJoin(PlayerJoinEvent e) {
        bridge.emitJoin(e.getPlayer().getName());
    }

    @EventHandler public void onChat(AsyncPlayerChatEvent e) {
        bridge.emitChat(e.getPlayer().getName(), e.getMessage());
    }
}
```

发布为 `.jar` 后按上架流程提交即可。

## 六、发布到市场

1. 将插件源码放入 `plugin/server/<id>/`（服务器端）或 `plugin/netease/<id>/`（网易端），或独立仓库。
2. 在 `registry/extensions.json` 添加条目，`kind` 填 `server`，`platforms` 标注适用端。
3. 配置 `.github/workflows/release.yml` 构建产物并挂到 Release。
4. `asset_url` 指向 Release 下载地址，`sha256` 建议填写（平台安装时校验）。

完整流程见 [插件制作完整指南](./plugin-development-guide.md)。