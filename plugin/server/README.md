# 服务器端插件

本目录存放部署在 **Minecraft 服务器端** 的插件（适用于国际版 Java / Bedrock 服务端，或作为常驻服务端的跨端插件）。

## 目录约定

```
plugin/server/
└── <plugin-id>/
    ├── README.md     # 插件说明
    ├── build.sh      # 构建 / 打包脚本（产物用于 Release）
    └── 源码
```

## 现有插件

| 插件 | 说明 |
| --- | --- |
| [fin-example](./fin-example/) | FIN 跨平台服务器端插件示例（进服验证 + 群服互通） |

## 构建

```bash
cd plugin/server/<plugin-id>
./build.sh   # 产物输出到 dist/
```

## 上架

构建产物对应的 `asset_url` 登记到顶层 `registry/extensions.json`。