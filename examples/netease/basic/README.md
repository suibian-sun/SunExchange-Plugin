# netease-example 示例插件（网易我的世界端）

网易我的世界端插件骨架示例，复用 SunExchange SDK 的「生命周期 + 事件 + 命令」模型。

## 结构

```
examples/netease/basic/
├── go.mod   # replace 指向 ../../../sdk
└── main.go  # 插件实现（网易端能力扩展点）
```

## 构建

```bash
cd examples/netease/basic
go mod tidy
go build -o netease-example .
```

## 说明

- 网易端（netease）插件面向网易租赁服：funcs 脚本、网易专属扩展等。
- 用脚手架生成更便捷：`sunexchange-plugin new <name> --end netease`。