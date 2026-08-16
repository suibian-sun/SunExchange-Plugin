# 脚手架 CLI 指南

`sunexchange-plugin` 是 SunExchange 插件脚手架命令行工具，一条命令生成完整的插件项目骨架（正确依赖、构建脚本、开发说明）。

## 安装

```bash
# 方式一：go install（推荐）
go install github.com/suibian-sun/SunExchange-Plugin/cli@latest

# 方式二：直接 go run（无需安装）
go run github.com/suibian-sun/SunExchange-Plugin/cli new myplugin
```

## 用法

```bash
sunexchange-plugin new <插件名> [--end server|netease] [--dir <目标目录>]
```

| 参数 | 默认 | 说明 |
| --- | --- | --- |
| `<插件名>` | 必填 | 插件名，同时作为 Go module 名与默认目录名 |
| `--end` | `server` | 插件端：`server`（服务器端）或 `netease`（网易我的世界端） |
| `--dir` | 插件名 | 生成的目标目录 |

## 示例

```bash
# 生成服务器端插件 myplugin
sunexchange-plugin new myplugin --end server

# 生成网易我的世界端插件 neteasemod
sunexchange-plugin new neteasemod --end netease

# 指定目录
sunexchange-plugin new foo --dir ./my-custom-dir
```

## 生成内容

### 服务器端（--end server）

```
myplugin/
├── go.mod     # 依赖 github.com/suibian-sun/SunExchange-Plugin/sdk
├── main.go    # 实现 sdk.Plugin，导出 Plugin 符号（.so）
└── README.md  # 构建与开发说明
```

- 默认监听玩家聊天事件，扩展点在 `Init`。
- `go build -buildmode=plugin -o myplugin.so .` 编译为 .so。

### 网易我的世界端（--end netease）

```
neteasemod/
├── go.mod
├── main.go    # 网易端能力扩展点
└── README.md
```

- 面向网易租赁服（funcs 脚本、网易专属扩展）。
- `go build -o neteasemod .` 编译为可执行文件。

## 下一步

生成后进入目标目录：

```bash
cd myplugin
go mod tidy          # 拉取依赖（SDK 发版后无需此步的 replace）
# 编辑 main.go 实现你的逻辑
go build -buildmode=plugin -o myplugin.so .   # 服务器端
```

部署：把产物放入插件目录，在 SunExchange 中登记后纳入群服互通能力。完整开发方法见 [SDK 开发指南](sdk-guide.md)。