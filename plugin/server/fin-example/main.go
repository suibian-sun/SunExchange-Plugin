package main

import (
	"github.com/maoqijie/FIN-plugin/sdk"
	"github.com/hashicorp/go-plugin"
)

// Entry 进服验证插件（全端）
type Entry struct{ ctx *sdk.Context }

// Init 是插件入口：注册事件监听
func (p *Entry) Init(ctx *sdk.Context) error {
	p.ctx = ctx

	// 玩家进服：触发进服验证与事件回写
	ctx.ListenPlayerJoin(func(event sdk.PlayerEvent) {
		ctx.LogSuccess("玩家 %s 加入，触发进服验证", event.Name)
		// 此处调用 SunExchange 校验接口，确认玩家已购买并拥有有效卡槽/任务
	})

	// 玩家聊天：转发到鹊桥（群服互通）
	ctx.ListenChat(func(event *sdk.ChatEvent) {
		ctx.LogInfo("聊天 %s: %s", event.Sender, event.Message)
		// 此处将消息转发到 SunExchange 鹊桥，再由 QQ 群侧回显
	})

	// 游戏控制：向服务端广播就绪信息
	if gu := ctx.GameUtils(); gu != nil {
		gu.SendChat("§eSunExchange 进服验证已就绪")
		gu.SendCommand("say SunExchange 进服验证已就绪")
	}

	return nil
}

func (p *Entry) Start() error { return nil }
func (p *Entry) Stop() error  { return nil }

// GetInfo 返回插件元信息
func (p *Entry) GetInfo() sdk.PluginInfo {
	return sdk.PluginInfo{
		Name:        "sunexchange-entry",
		DisplayName: "进服验证插件",
		Version:     "1.0.0",
		Description: "玩家进服校验并回写事件到鹊桥，支撑群服互通",
		Author:      "SunExchange",
	}
}

func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: sdk.HandshakeConfig,
		Plugins: map[string]plugin.Plugin{
			"plugin": &sdk.PluginGRPC{Impl: &Entry{}},
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
}