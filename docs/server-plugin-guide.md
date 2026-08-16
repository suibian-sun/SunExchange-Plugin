# 服务器端插件制作指南

本指南说明如何为 SunExchange 统一插件市场制作 **Minecraft 服务器端插件**。
服务器端插件参考 [cat7street/FIN-plugin](https://github.com/cat7street/FIN-plugin) 的跨平台插件框架设计，
同时支持经典服务端插件（Spark/Paper、基岩 Nukkit/BedrockServer、网易）。

## 一、插件分类与适用端

| 适用端 | 平台标识 | 服务端类型 | 插件形态 |
| --- | --- | --- | --- |
| 网易我的世界 | `netease` | FunShuttler / 网易租赁服 | 网易 funcs 插件 / 脚本 |
| 国际 Java 版 | `java` | Spigot / Paper / Purpur（Bukkit） | `.jar` |
| 国际基岩版 | `bedrock` | BDS / Nukkit / PocketMine | `.jar` / `.js` / 插件文件夹 |

## 二、插件的核心职责（群服互通）

服务器端插件在游戏内完成三件事，与 SunExchange 的鹊桥（queqiao）协议对接：

1. **进服验证**：玩家加入时，向 SunExchange 校验该玩家是否已购买并拥有有效卡槽/任务。
2. **聊天互通**：把玩家聊天转发为 `PlayerChatEvent`，并接收来自 QQ 群的指令向玩家/服务器广播。
3. **指令与事件转发**：接收 `SendCommand` / `SendWhisper`，执行服务端命令并回传结果。

## 三、事件协议（鹊桥 WebSocket）

服务器端插件通过 WebSocket 连接 SunExchange 的鹊桥服务：

```jsonc
// 玩家聊天 → 上行
{ "post_type": "message", "event_name": "PlayerChatEvent",
  "server_name": "...", "server_type": "java", "server_ver": "java",
  "message": "hello", "player": { "name": "Steve", "uid": "..." } }

// 群服指令 → 下行
{ "event_name": "SendCommandEvent", "command": "/give @a diamond 64" }
```

## 四、使用 FIN 跨平台框架（推荐）

FIN 插件框架（`cat7street/FIN-plugin`）支持跨平台插件系统（Windows/Linux/macOS/Android）与传统 `.so` 插件。

```go
// 监听玩家加入
ctx.ListenPlayerJoin(func(event sdk.PlayerEvent) {
    ctx.LogSuccess("玩家 %s 加入", event.Name)
})

// 游戏控制
gu := ctx.GameUtils()
gu.SayTo("玩家名", "§a你好！")
gu.SendCommand("/give @a diamond 64")
```

参考本仓库 `plugins/fin-example/` 的完整示例。

## 五、发布到市场

1. 将插件源码放入 `plugins/<id>/`。
2. 在 `registry/extensions.json` 添加条目，`kind` 填 `server`，`platforms` 标注适用端。
3. 配置 `.github/workflows/release.yml` 构建产物并挂到 Release。
4. `asset_url` 指向 Release 下载地址，`sha256` 可留空（由平台校验）。

## 六、经典服务端插件（Java/Bukkit）

```java
public final class EntryPlugin extends JavaPlugin implements Listener {
    @Override public void onEnable() {
        getServer().getPluginManager().registerEvents(this, this);
    }
    @EventHandler public void onJoin(PlayerJoinEvent e) {
        // 进服验证：调用 SunExchange 校验接口
        validateAndBridge(e.getPlayer());
    }
    @EventHandler public void onChat(AsyncPlayerChatEvent e) {
        // 聊天互通：转发到鹊桥
        bridgeChat(e.getPlayer(), e.getMessage());
    }
}
```

发布为 `.jar` 后按上述流程上架即可。