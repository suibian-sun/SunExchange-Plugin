// Command sunexchange-plugin 是 SunExchange 插件脚手架命令行工具。
//
// 用法：
//
//	# 生成一个服务器端插件骨架
//	sunexchange-plugin new myplugin --end server
//
//	# 生成一个网易我的世界端插件骨架
//	sunexchange-plugin new myfunc --end netease
//
//	# 省略 --end 时默认 server
//	sunexchange-plugin new foo
//
// 生成完毕后，进入目标目录按 README 完成开发与构建。
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func usage() {
	fmt.Fprint(os.Stderr, `SunExchange 插件脚手架

用法:
  sunexchange-plugin new <插件名> [--end server|netease] [--dir <目标目录>]

示例:
  sunexchange-plugin new myplugin --end server
  sunexchange-plugin new neteasemod --end netease

`)
}

func main() {
	raw := os.Args[1:]
	// 定位子命令：第一个参数必须是 new
	if len(raw) < 2 || raw[0] != "new" {
		usage()
		os.Exit(1)
	}
	name := raw[1]

	// 解析剩余 flags（--end / --dir）
	end := "server"
	dir := ""
	for i := 2; i < len(raw); i++ {
		switch raw[i] {
		case "--end":
			if i+1 < len(raw) {
				end = raw[i+1]
				i++
			}
		case "--dir":
			if i+1 < len(raw) {
				dir = raw[i+1]
				i++
			}
		default:
			fmt.Fprintf(os.Stderr, "未知参数 %q\n", raw[i])
			usage()
			os.Exit(1)
		}
	}

	var tpl *Template
	switch end {
	case "server":
		tpl = serverTemplate
	case "netease":
		tpl = neteaseTemplate
	default:
		fmt.Fprintf(os.Stderr, "未知的插件端 %q（可选 server / netease）\n", end)
		os.Exit(1)
	}

	target := dir
	if target == "" {
		target = name
	}
	if err := scaffold(name, tpl, target); err != nil {
		fmt.Fprintf(os.Stderr, "生成插件失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✔ 已在 %s 生成 %s 端插件骨架 %q\n", target, tpl.EndName, name)
	fmt.Printf("  下一步：cd %s && 阅读 README.md 完成开发\n", target)
}

// File 待生成的一个文件。
type File struct {
	Path    string // 相对目标目录的路径
	Mode    os.FileMode
	Content string
}

// Template 一类插件的骨架描述。
type Template struct {
	EndName string // 端的中文名（服务器端 / 网易我的世界端）
	Files   []File
}

var serverTemplate = &Template{
	EndName: "服务器端",
	Files: []File{
		{Path: "go.mod", Content: serverGoMod},
		{Path: "main.go", Content: serverMainGo},
		{Path: "README.md", Content: serverReadme},
	},
}

var neteaseTemplate = &Template{
	EndName: "网易我的世界端",
	Files: []File{
		{Path: "main.go", Content: neteaseMainGo},
		{Path: "README.md", Content: neteaseReadme},
	},
}

// scaffold 把模板写入目标目录，替换 {{.Name}} 占位符。
func scaffold(name string, tpl *Template, target string) error {
	if _, err := os.Stat(target); err == nil {
		return fmt.Errorf("目标目录 %s 已存在", target)
	}

	data := struct{ Name string }{Name: name}
	for _, f := range tpl.Files {
		path := filepath.Join(target, filepath.FromSlash(f.Path))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return err
		}
		t, err := template.New(f.Path).Parse(f.Content)
		if err != nil {
			return fmt.Errorf("解析模板 %s: %v", f.Path, err)
		}
		var b strings.Builder
		if err := t.Execute(&b, data); err != nil {
			return fmt.Errorf("渲染模板 %s: %v", f.Path, err)
		}
		mode := f.Mode
		if mode == 0 {
			mode = 0o644
		}
		if err := os.WriteFile(path, []byte(b.String()), mode); err != nil {
			return err
		}
	}
	return nil
}

const serverGoMod = `module {{.Name}}

go 1.18

// 依赖 SunExchange 插件 SDK（独立模块）。发布后 SDK 有版本 tag 时替换 v1.0.0。
require github.com/suibian-sun/SunExchange-Plugin/sdk v1.0.0
`

const serverMainGo = `// {{.Name}} SunExchange 服务器端插件骨架。
//
// 实现 sdk.Plugin 接口，由宿主加载并驱动生命周期。
// 默认作为 .so 动态加载（go build -buildmode=plugin 导出 Plugin 符号）。
package main

import (
	"github.com/suibian-sun/SunExchange-Plugin/sdk"
)

// Plugin 插件实例。
type MyPlugin struct {
	ctx *sdk.Context
}

// NewPlugin 构造函数。
func NewPlugin() sdk.Plugin { return &MyPlugin{} }

// GetInfo 返回插件元信息。
func (p *MyPlugin) GetInfo() sdk.PluginInfo {
	return sdk.PluginInfo{
		Name:        "{{.Name}}",
		DisplayName: "{{.Name}}",
		Version:     "0.1.0",
		Description: "由 SunExchange 脚手架生成的服务器端插件",
		Author:      "you",
	}
}

// Init 插件加载时调用：注册事件、命令、插件 API。
func (p *MyPlugin) Init(ctx *sdk.Context) error {
	p.ctx = ctx

	// 示例：监听玩家聊天
	ctx.ListenChat(func(c *sdk.ChatEvent) bool {
		ctx.Infof("聊天 %s: %s", c.Sender, c.Message)
		return false
	})

	return nil
}

// Start 互通链路就绪后调用。
func (p *MyPlugin) Start() error {
	p.ctx.Successf("%s 已启动", p.GetInfo().DisplayName)
	return nil
}

// Stop 卸载 / 退出时调用，释放资源。
func (p *MyPlugin) Stop() error {
	p.ctx.Warnf("%s 已停止", p.GetInfo().DisplayName)
	return nil
}

// Plugin 是 .so 动态加载必须导出的符号。
func Plugin() sdk.Plugin { return NewPlugin() }

// main 空函数：Go 插件包要求包含 main。
func main() {}
`

const serverReadme = `# {{.Name}}

SunExchange 服务器端插件骨架（由脚手架生成）。

## 开发

1. 按需修改 'main.go' 中的 'Init'：注册事件监听、命令、插件 API。
2. 能力参考：事件监听/拦截、游戏 API、命令、插件间 API、广播、底层透传（见 SDK 文档）。

## 构建为 .so

    go mod tidy
    go build -buildmode=plugin -o {{.Name}}.so .

产物 '{{.Name}}.so' 放入插件目录即可被 SunExchange 宿主动态加载。

## 本地测试

    go vet ./... && go build ./...
`

const neteaseMainGo = `// {{.Name}} 网易我的世界端插件骨架。
//
// 网易端插件通常以租赁服 funcs 脚本 / 网易专属能力形式存在。
// 此骨架给出与 SunExchange SDK 对齐的入口结构，可按需扩展。
package main

import (
	"github.com/suibian-sun/SunExchange-Plugin/sdk"
)

// Plugin 插件实例。
type MyPlugin struct {
	ctx *sdk.Context
}

// NewPlugin 构造函数。
func NewPlugin() sdk.Plugin { return &MyPlugin{} }

// GetInfo 返回插件元信息。
func (p *MyPlugin) GetInfo() sdk.PluginInfo {
	return sdk.PluginInfo{
		Name:        "{{.Name}}",
		DisplayName: "{{.Name}}",
		Version:     "0.1.0",
		Description: "由 SunExchange 脚手架生成的网易我的世界端插件",
		Author:      "you",
	}
}

// Init 注册事件与能力。
func (p *MyPlugin) Init(ctx *sdk.Context) error {
	p.ctx = ctx
	// 网易端能力在此扩展（funcs 脚本桥接、网易专属事件等）
	return nil
}

// Start 启动逻辑。
func (p *MyPlugin) Start() error {
	p.ctx.Successf("%s 已启动", p.GetInfo().DisplayName)
	return nil
}

// Stop 释放资源。
func (p *MyPlugin) Stop() error {
	p.ctx.Warnf("%s 已停止", p.GetInfo().DisplayName)
	return nil
}

func Plugin() sdk.Plugin { return NewPlugin() }

func main() {}
`

const neteaseReadme = `# {{.Name}}

网易我的世界端插件骨架（由脚手架生成）。

## 说明

- 网易端（netease）插件面向网易租赁服：funcs 脚本、网易专属扩展等。
- 骨架基于 SunExchange SDK，复用「生命周期 + 事件 + 命令」模型。

## 构建

    go mod tidy
    go build -o {{.Name}}. .

## 部署

将产物部署到网易租赁服对应目录，并在 SunExchange 中登记后即可纳入群服互通能力。
`