// {{.Name}} 网易我的世界端插件骨架示例。
//
// 网易端插件通常以租赁服 funcs 脚本 / 网易专属能力形式存在。
// 此示例给出与 SunExchange SDK 对齐的入口结构，可按需扩展。
// 通过脚手架生成：sunexchange-plugin new <name> --end netease
package main

import (
	"github.com/suibian-sun/SunExchange-Plugin/sdk"
)

// MyPlugin 插件实例。
type MyPlugin struct {
	ctx *sdk.Context
}

// NewPlugin 构造函数。
func NewPlugin() sdk.Plugin { return &MyPlugin{} }

// GetInfo 返回插件元信息。
func (p *MyPlugin) GetInfo() sdk.PluginInfo {
	return sdk.PluginInfo{
		Name:        "netease-example",
		DisplayName: "netease-example",
		Version:     "0.1.0",
		Description: "网易我的世界端插件示例",
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