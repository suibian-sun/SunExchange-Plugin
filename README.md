# SunExchange 统一插件市场

这里维护 SunExchange 的**统一插件市场注册表**，合并了两类插件，通过**分类筛选**区分：

| kind | 说明 | 适用 |
| --- | --- | --- |
| `unibot` | **UniBot 扩展**（QQ 端） | 渲染引擎 / 模板 / 指令扩展，平台标记为 `all` |
| `server` | **Minecraft 服务器端插件** | 部署到服务器端，`platforms` 标记适用端 |

## 目录结构

```
SunExchange-Plugin/
├── registry/
│   ├── extensions.json      # 统一插件市场注册表（唯一数据源）
│   └── schema.json          # 注册表 JSON SCHEMA
├── plugin/                  # 插件源码与打包脚本
│   ├── server/              # 服务器端插件（国际版 Java/Bedrock 或常驻服务端的跨端插件）
│   │   └── fin-example/     # FIN 跨平台插件示例
│   └── netease/             # 网易我的世界端插件
├── docs/                    # 插件制作文档（指南 + SDK）
└── .github/workflows/       # 校验 + 打 Release 包
```

## 插件分类

`plugin/` 目录按**部署位置**分为两类：

| 分类 | 目录 | 说明 |
| --- | --- | --- |
| 服务器端 | `plugin/server/` | 部署在 Minecraft 服务器端：国际版 Java/Bedrock 服务端插件，或常驻服务端的跨端插件 |
| 我的世界端 | `plugin/netease/` | 部署在网易我的世界（netease）端：租赁服 funcs 脚本、网易专属扩展 |

> 全端插件（同时支持 netease/java/bedrock）通常作为常驻服务端的跨端插件归入 `plugin/server/`。

## 文档导航

| 文档 | 内容 |
| --- | --- |
| [插件制作完整指南](docs/plugin-development-guide.md) | 两类插件总览 + 筛选机制 + 发布流程 + 常见问题（推荐先读） |
| [SunExchange 插件 SDK 开发指南](docs/sdk-plugin-guide.md) | SunExchange 原生插件 SDK：生命周期 + 事件 + 命令 + 插件间 API + 高自由度扩展 |
| [服务器端插件制作指南](docs/server-plugin-guide.md) | Minecraft 服务器端插件：鹊桥协议 + FIN 框架 + 经典 Bukkit |
| [UniBot 扩展制作指南](docs/unibot-extension-guide.md) | UniBot 扩展：目录结构 + Extension.toml + 打包上架 |

## 注册表格式

每个插件条目通过两个维度参与筛选：

1. **分类（kind）**：`unibot`（UniBot 扩展）或 `server`（服务器端插件）。
2. **适用端（platforms）**：`netease` / `java` / `bedrock` / `all`。

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
  "official": true,
  "releases": [
    {
      "version": "1.0.0",
      "asset_url": "https://github.com/suibian-sun/SunExchange-Plugin/releases/download/v1.0.0/x.jar",
      "sha256": ""
    }
  ]
}
```

字段说明：

| 字段 | 必填 | 说明 |
| --- | --- | --- |
| `id` | 是 | 唯一标识，`^[A-Za-z0-9_.-]+$` |
| `name` | 是 | 插件显示名 |
| `kind` | 是 | `unibot` 或 `server` |
| `platforms` | 是 | 适用端数组，`server` 插件必须明确标注 |
| `category` | 是 | 功能分类（互通 / 经济 / 管理 / 模板 / 渲染引擎 / API） |
| `repo` | 是 | 插件源码仓库 `owner/repo` |
| `description` | 是 | 插件描述 |
| `official` | 否 | 是否官方 |
| `releases` | 是 | 版本列表，`version` + `asset_url` 必填 |

## 如何使用

SunExchange 主程序从本仓库的 `registry/extensions.json` 拉取市场数据。
登录 SunExchange 后进入「插件市场」页面，即可：

- 按**分类**（UniBot 扩展 / 服务器端插件）筛选；
- 按**适用端**（网易 / 国际 Java / 国际基岩）筛选；
- 查看插件详情、下载 `asset_url`、按需安装。

## 发布新插件

1. 在 `registry/extensions.json` 中新增条目（或添加版本）。
2. 插件源码放入 `plugin/<分类>/<id>/`（服务器端放 `plugin/server/`，网易端放 `plugin/netease/`），并配置 Release 打包。
3. 推送 `main` 分支 → GitHub Actions 自动校验注册表并构建 Release。
4. 在 Releases 页创建 `vX.Y.Z` 标签发布产物，`asset_url` 指向对应 Release。

## 本地校验

```bash
python -m json.tool registry/extensions.json > /dev/null && echo "JSON 合法"
```