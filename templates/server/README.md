# 服务器端插件模板

本目录提供 **服务器端插件** 的静态模板，与 CLI（`sunexchange-plugin new`）生成的内容一致，供手动复制或参考。

```
templates/server/basic/
├── go.mod      # 依赖 SunExchange 插件 SDK
├── main.go     # 实现 sdk.Plugin，导出 Plugin 符号（.so）
└── README.md   # 构建与开发说明
```

## 使用方式

方式一（推荐）：用 CLI 生成

```bash
sunexchange-plugin new myplugin --end server
```

方式二：手动复制本目录到项目，改 `module` 名与插件名。

## 能力

模板默认监听玩家聊天事件，扩展点见 `main.go` 的 `Init`：事件监听/拦截、命令、游戏 API、插件间 API、广播、底层透传（详见 `docs/sdk-guide.md`）。