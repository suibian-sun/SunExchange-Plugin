# 插件制作完整指南

本指南面向希望为 SunExchange 统一插件市场提交插件的开发者，覆盖两类插件从设计、编码、打包到上架发布的全流程。文中所有协议字段均以 SunExchange 鹊桥服务的实际实现为准，代码示例可直接运行。

## 一、市场机制

SunExchange 插件市场把两类插件合并到一个注册表，靠**分类**与**适用端**两个维度区分和筛选。

| kind | 是什么 | 运行位置 | 由谁加载 | 适用端 |
| --- | --- | --- | --- | --- |
| `unibot` | UniBot 扩展（指令、模板、渲染引擎） | QQ 机器人端 | UniBot | 不区分端，标记 `all` |
| `server` | Minecraft 服务器端插件 | 游戏服务器 | 服务端 | `netease` / `java` / `bedrock` |

玩家在「插件市场」页面按这两个维度筛选，再决定下载还是安装。因此你生产的插件最终只做两件事：**写好插件本身**，**在注册表里正确标注它是哪类、适用于哪些端**。

三类适用端对应关系：

| 适用端 | 平台标识 | 典型服务端 | 插件形态 |
| --- | --- | --- | --- |
| 网易我的世界 | `netease` | FunShuttler / 网易租赁服 | 网易 funcs 插件 / 脚本 |
| 国际 Java 版 | `java` | Spigot / Paper / Purpur | `.jar` |
| 国际基岩版 | `bedrock` | BDS / Nukkit / PocketMine | `.jar` / `.js` / 插件文件夹 |

## 二、选择插件类型

先回答一个问题：你的功能需要在哪里生效？

- 功能发生在 **QQ 群聊**（如把 HTML 渲染成图片、解析群指令、查询数据）→ 做 **UniBot 扩展**。
- 功能发生在 **游戏内**（如进服校验、踢人、发指令、改经济）→ 做 **服务器端插件**。
- 功能横跨两端（如"群友发指令，游戏内执行并回显"）→ 需要**服务器端插件** + UniBot 扩展配合，通过鹊桥打通。

## 三、制作 UniBot 扩展

UniBot 扩展是放在 `unibot/extensions/<id>/` 下的一个目录，包含配置、模板与资源。

### 目录结构

```
<id>/
├── Extension.toml      # 扩展清单（必填）
├── Templates/          # 图片/文字模板（HTML 或占位符模板）
├── Resources/          # 字体、背景图等静态资源
└── 其他业务文件
```

### Extension.toml 清单

```toml
[extension]
id = "MyExtension"          # 必填，与注册表 id 一致
name = "我的扩展"
version = "1.0.0"
description = "扩展功能说明"
author = "YourName"
unibot_version = "*"        # 兼容的 UniBot 版本，* 表示全部
```

### 模板与资源

- `Templates/` 下的 HTML 文件由渲染引擎（如 Html2Pic）转成图片发出，可用占位符语法嵌入动态数据。
- `Resources/` 存放字体、背景图等被模板引用的静态资源。
- 需要查询游戏数据时，通过 Placeholder 占位符桥（RCON）或服务器端插件回写的数据驱动渲染。

### 打包

把整个目录打成 zip，命名建议 `Id-版本.zip`，并作为 Release 资产上传，`asset_url` 指向该 zip。

## 四、制作服务器端插件

服务器端插件是市场的主体，也是群服互通的核心。推荐用 **FIN 跨平台框架** 开发，它能一套代码同时跑在三个适用端。

### 4.1 用 FIN 框架（推荐）

FIN（FunInterWork）插件框架用 Go 编写，支持 Linux/macOS/Windows/Android，通过 gRPC 与宿主机通信，进程隔离、跨端一致。

```go
package main

import (
    "github.com/cat7street/FIN-plugin/sdk"
    "github.com/cat7street/FIN-plugin/sdk/plugin"
)

type MyPlugin struct{}

func (p *MyPlugin) OnLoad(ctx *sdk.Context) error {
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
    // 监听聊天
    ctx.ListenChat(func(event *sdk.ChatEvent) {
        ctx.LogInfo("%s: %s", event.Sender, event.Message)
    })
    return nil
}

func main() { plugin.Serve(&MyPlugin{}) }
```

### 4.2 游戏控制工具

```go
gu := ctx.GameUtils()
gu.SayTo("玩家名", "§a你好！")        // 向指定玩家发消息
gu.Broadcast("§e全服公告")             // 全服广播
gu.SendCommand("/give @a diamond 64") // 执行服务端命令
```

### 4.3 接入鹊桥（群服互通）

服务器端插件通过 WebSocket 连接 SunExchange 的鹊桥服务，地址为 `ws://<sunexchange>:<port>/minecraft/ws`。每次连接用插件配置里的 `server_name`、`server_type` 声明身份。

**上行事件（游戏 → 群）**，JSON 结构固定含 `event_name`、`server_*`、`player` 等字段：

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

常用事件：

| event_name | post_type | 触发时机 |
| --- | --- | --- |
| `PlayerJoinEvent` | notice | 玩家进服 |
| `PlayerQuitEvent` | notice | 玩家离开 |
| `PlayerChatEvent` | message | 玩家聊天 |

**下行指令（群 → 游戏）**，服务器端插件收到后执行并把结果回包：

```json
{ "event_name": "SendCommandEvent", "command": "/give @a diamond 64" }
{ "event_name": "SendWhisperEvent", "player": "Steve", "message": "你好" }
```

### 4.4 三种适用端的实现要点

| 适用端 | 说明 |
| --- | --- |
| `java` | 若不想用 FIN，可直接写 Spigot/Paper 插件（见下）。 |
| `bedrock` | FIN 或 BDS `.js` / Nukkit `.jar`。 |
| `netease` | 网易租赁服加载机制特殊，通常用 FIN 或网易 funcs 脚本对接。 |

### 4.5 经典 Java（Bukkit）插件

不依赖 FIN 框架时，直接写标准 Bukkit 插件，核心是 `onEnable` 注册监听器，再把事件转发到鹊桥：

```java
public final class EntryPlugin extends JavaPlugin implements Listener {
    private BridgeClient bridge;

    @Override public void onEnable() {
        getServer().getPluginManager().registerEvents(this, this);
        bridge = new BridgeClient(getConfig().getString("bridge.ws"));
        bridge.connect();
    }

    @EventHandler public void onJoin(PlayerJoinEvent e) {
        bridge.emitJoin(e.getPlayer().getName());
    }

    @EventHandler public void onChat(AsyncPlayerChatEvent e) {
        bridge.emitChat(e.getPlayer().getName(), e.getMessage());
    }
}
```

## 五、在注册表登记插件

写好后，在 `registry/extensions.json` 的 `extensions` 数组里加一个条目。必填字段以 `kind` 和 `platforms` 为核心：

```json
{
  "id": "sunexchange.entry",
  "name": "进服验证插件（全端）",
  "kind": "server",
  "platforms": ["netease", "java", "bedrock"],
  "category": "互通",
  "author": "SunExchange",
  "repo": "suibian-sun/SunExchange-Plugin",
  "description": "玩家进服校验并回写事件到鹊桥",
  "official": false,
  "releases": [
    {
      "version": "1.0.0",
      "asset_url": "https://github.com/owner/repo/releases/download/v1.0.0/entry.jar",
      "sha256": "可选，建议填写"
    }
  ]
}
```

字段速查：

| 字段 | 必填 | 说明 |
| --- | --- | --- |
| `id` | 是 | 唯一标识，`^[A-Za-z0-9_.-]+$` |
| `name` | 是 | 显示名 |
| `kind` | 是 | `unibot` 或 `server` |
| `platforms` | 是 | 适用端数组；`server` 必须明确标注，不支持的端不要写 |
| `category` | 是 | 功能分类（互通 / 经济 / 管理 / 模板 / 渲染引擎 / API） |
| `repo` | 是 | 源码仓库 `owner/repo` |
| `description` | 是 | 插件描述 |
| `official` | 否 | 是否官方 |
| `releases` | 是 | 版本列表，至少一项 `version` + `asset_url` |

`platforms` 的合法值：`netease`、`java`、`bedrock`、`all`。`all` 表示不区分端（UniBot 扩展通常用 `all`）。

## 六、发布流程

1. 插件源码放入独立仓库（或本仓库 `plugins/<id>/`），配置 Release 打包。
2. 在 `registry/extensions.json` 新增条目。
3. 推送 `main` 分支，GitHub Actions 会自动校验注册表（JSON 合法性 + schema）。
4. 打 `vX.Y.Z` 标签，生成 Release 资产，把 `asset_url` 指向该资产下载地址。
5. `sha256` 建议填写，SunExchange 安装时会校验，防止下载文件被篡改。

本地校验注册表：

```bash
python -m json.tool registry/extensions.json > /dev/null && echo "JSON 合法"
```

## 七、最佳实践

- **明确适用端**：只在 `platforms` 列出真正测试过的端，避免误导用户。
- **填写 sha256**：让安装过程可校验完整性。
- **提供描述**：`description` 说清楚功能与依赖，用户仅凭它决定是否下载。
- **保持 id 稳定**：`id` 是安装目录和去重的依据，发布后不要随意改动。
- **先做最小版本**：一个端验证通过再扩展其他端，减少跨端联调成本。
- **群服互通建议用 FIN**：跨端一致、进程隔离，避免为每个端各写一套。

## 八、常见问题

**Q：为什么群里的指令没在游戏内执行？**
A：确认服务器端插件已连上鹊桥，且 `server_name` / `server_type` 与任务配置一致；下行事件需在插件里处理 `SendCommandEvent`。

**Q：UniBot 扩展和服务器端插件必须成对出现吗？**
A：不一定。纯 QQ 端功能只需要 UniBot 扩展；纯游戏内功能只需要服务器端插件；只有群服互通才需要两者配合。

**Q：`platforms=all` 是什么意思？**
A：表示不区分适用端，主要用于不依赖具体游戏端的 UniBot 扩展。服务器端插件应明确标注 `netease` / `java` / `bedrock`。