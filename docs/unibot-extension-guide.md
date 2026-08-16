# UniBot 扩展制作指南

本指南说明如何为 SunExchange 统一插件市场制作 **UniBot 扩展**。UniBot 扩展运行在 QQ 机器人端，负责指令解析、模板渲染与数据查询，是市场上 `kind=unibot` 的一类插件。

## 一、扩展是什么

UniBot 扩展是放在 `unibot/extensions/<id>/` 目录下的一个独立扩展包，包含清单、模板与资源。它不区分适用端，注册表里 `platforms` 标记为 `all`。

典型用途：

- 指令模板（默认模板、帮助、绑定、展示列表）
- 渲染引擎（如 Html2Pic，把 HTML 转成图片）
- 占位符桥（如 Placeholder，通过 RCON 查询服务器数据）
- 群服互通时配合服务器端插件，把游戏事件渲染成图片发出

## 二、目录结构

```
<id>/
├── Extension.toml      # 扩展清单（必填）
├── Templates/          # 图片/文字模板（HTML 或占位符模板）
├── Resources/          # 字体、背景图等静态资源
└── 其他业务文件
```

## 三、Extension.toml 清单

```toml
[extension]
id = "MyExtension"          # 必填，与注册表 id 一致
name = "我的扩展"
version = "1.0.0"
description = "扩展功能说明"
author = "YourName"
unibot_version = "*"        # 兼容的 UniBot 版本，* 表示全部
```

| 字段 | 必填 | 说明 |
| --- | --- | --- |
| `id` | 是 | 扩展标识，安装目录名，需与注册表 `id` 一致 |
| `name` | 是 | 显示名 |
| `version` | 是 | 扩展版本 |
| `description` | 是 | 功能说明 |
| `author` | 否 | 作者 |
| `unibot_version` | 否 | 兼容的 UniBot 版本 |

## 四、模板与资源

- `Templates/` 下的 HTML 文件由渲染引擎（如 Html2Pic）转成图片发出，可在 HTML 里用占位符语法嵌入动态数据。
- `Resources/` 存放字体、背景图等被模板引用的静态资源。
- 需要查询游戏内数据时，通过 Placeholder 占位符桥（RCON）或由服务器端插件经鹊桥回写的数据驱动渲染。

## 五、在注册表登记

在 `registry/extensions.json` 的 `extensions` 里加条目，`kind` 填 `unibot`，`platforms` 填 `["all"]`：

```json
{
  "id": "MyExtension",
  "name": "我的扩展",
  "kind": "unibot",
  "platforms": ["all"],
  "category": "模板",
  "author": "YourName",
  "repo": "yourname/MyExtension",
  "description": "扩展功能说明",
  "official": false,
  "releases": [
    {
      "version": "1.0.0",
      "asset_url": "https://github.com/yourname/MyExtension/releases/download/v1.0.0/MyExtension-1.0.0.zip",
      "sha256": "可选，建议填写",
      "unibot_version": "*"
    }
  ]
}
```

## 六、打包与发布

1. 把整个扩展目录打成 zip，命名建议 `Id-版本.zip`。
2. 作为 Release 资产上传，`asset_url` 指向该 zip。
3. `sha256` 建议填写，SunExchange 安装时校验文件完整性。

安装时，SunExchange 会下载 zip 并解压到 `unibot/extensions/<id>/`，随后 UniBot 即可加载该扩展。

完整流程见 [插件制作完整指南](./plugin-development-guide.md)。