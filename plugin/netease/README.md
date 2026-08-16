# 我的世界端插件

本目录存放部署在 **网易我的世界（netease）** 端的插件，包括网易租赁服 `funcs` 脚本、网易专属扩展等。

## 目录约定

```
plugin/netease/
└── <plugin-id>/
    ├── README.md     # 插件说明
    ├── build.sh      # 构建 / 打包脚本（产物用于 Release）
    └── 源码
```

## 现有插件

暂无。本分类用于放置适用于网易我的世界端的插件。

> 提示：若插件同时支持网易端与国际版，通常作为常驻服务端的跨端插件放入 [`plugin/server/`](../server/)（参见全端插件归属约定）。

## 上架

构建产物对应的 `asset_url` 登记到顶层 `registry/extensions.json`。