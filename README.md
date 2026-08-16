# SunExchange 插件开发仓库

本仓库专注 **SunExchange 插件开发所需的一切文件**：插件 SDK、脚手架 CLI、模板、示例与文档。前端与后端（含注册表 API 服务）位于 [SunExchange](https://github.com/suibian-sun/SunExchange) 主仓库。

## 目录结构

```
SunExchange-Plugin/
├── sdk/          # 插件 SDK（独立 Go module，纯标准库）
├── cli/          # 脚手架 CLI（sunexchange-plugin new）
├── templates/    # 插件模板（按端分类：server / netease）
├── examples/     # 可运行示例（按端分类）
├── docs/         # 插件开发文档（SDK 参考 + 指南）
└── .gitignore
```

## 快速开始

```bash
# 1. 用 CLI 生成插件骨架
go run github.com/suibian-sun/SunExchange-Plugin/cli new myplugin --end server

# 2. 进入目录，开发你的逻辑
cd myplugin
go mod tidy

# 3. 编译为 .so（服务器端）
go build -buildmode=plugin -o myplugin.so .
```

## SDK 引用

插件 SDK 为独立 Go module，可被任意插件项目直接引用：

```bash
go get github.com/suibian-sun/SunExchange-Plugin/sdk@latest
```

开发期可用 `replace` 指向本地：

```
replace github.com/suibian-sun/SunExchange-Plugin/sdk => /path/to/SunExchange-Plugin/sdk
```

## 文档导航

| 文档 | 内容 |
| --- | --- |
| [SDK 开发指南](docs/sdk-guide.md) | 插件开发全流程（推荐先读） |
| [SDK API 参考](docs/sdk-reference.md) | `sdk` 包类型 / 方法速查 |
| [脚手架 CLI 指南](docs/cli-guide.md) | 用 CLI 一键生成插件 |

## 插件端分类

| 分类 | 目录 | 说明 |
| --- | --- | --- |
| 服务器端 | `templates/server/`、`examples/server/` | 国际版 Java/Bedrock 服务端插件，或常驻服务端的跨端插件 |
| 网易我的世界端 | `templates/netease/`、`examples/netease/` | 网易租赁服 funcs 脚本、网易专属扩展 |

## 关联仓库

- [SunExchange](https://github.com/suibian-sun/SunExchange)：主进程（前端 + 后端 + 注册表 API 服务）