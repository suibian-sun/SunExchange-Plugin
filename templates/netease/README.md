# 网易我的世界端插件模板

本目录提供 **网易我的世界端（netease）插件** 的静态模板，与 CLI（`sunexchange-plugin new`）生成的内容一致，供手动复制或参考。

```
templates/netease/basic/
├── main.go     # 实现 sdk.Plugin，网易端能力扩展点
└── README.md   # 构建与部署说明
```

> 网易端插件面向网易租赁服（funcs 脚本、网易专属扩展）。复用了 SunExchange SDK 的「生命周期 + 事件 + 命令」模型。