// Package entry 示例插件：进服验证 + 群服互通。
//
// 演示 SunExchange 插件 SDK 的高自由度能力：
//   - 生命周期：Init / Start / Stop / GetInfo
//   - 事件监听：玩家进服、离开、聊天（带优先级）
//   - 命令注册：控制台 / 群聊命令
//   - 游戏 API：全服广播、私聊、执行命令、标题
//   - 插件间 API：注册 / 获取其他插件能力
//   - 底层透传：Raw 任意协议载荷
package entry

import (
	"fmt"
	"strings"

	"github.com/suibian-sun/SunExchange-Plugin/sdk"
)

// Entry 进服验证插件（全端）。
type Entry struct {
	ctx       *sdk.Context
	whitelist map[string]bool
}

// NewEntry 创建插件实例。
func NewEntry() sdk.Plugin { return &Entry{} }

// Init 初始化：注册事件监听、命令、插件 API。
func (p *Entry) Init(ctx *sdk.Context) error {
	p.ctx = ctx
	p.whitelist = map[string]bool{}

	// 1. 事件监听：玩家进服（最高优先级，先做校验）
	ctx.ListenWithPriority("PlayerJoinEvent", func(ev any) bool {
		pe := ev.(sdk.PlayerEvent)
		if !p.whitelist[pe.Name] {
			ctx.Warnf("玩家 %s 进服，但不在白名单，广播提示", pe.Name)
			_ = ctx.Game().SendChat("§e[验证] " + pe.Name + " 请先通过商城购买获取进服资格")
		} else {
			ctx.Successf("玩家 %s 进服，白名单校验通过", pe.Name)
		}
		return false
	}, sdk.PriorityHigh)

	// 2. 事件监听：玩家聊天（转发到群服互通）
	ctx.ListenChat(func(c *sdk.ChatEvent) bool {
		ctx.Infof("聊天 %s: %s", c.Sender, c.Message)
		return false
	})

	// 3. 命令注册：群聊指令 /admin 手动加白
	_ = ctx.RegisterCommand(sdk.Command{
		Name:        "sxadmin",
		Triggers:    []string{"sx", "进服验证"},
		ArgHint:     "<玩家名>",
		Usage:       "sxadmin <玩家名> - 将玩家加入进服白名单",
		Description: "进服验证白名单管理",
		Handler: func(args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("用法: %s", "sxadmin <玩家名>")
			}
			p.whitelist[args[0]] = true
			ctx.Successf("已将 %s 加入进服白名单", args[0])
			return nil
		},
	})

	// 4. 插件 API：向其他插件暴露"进服资格校验"
	_ = ctx.RegisterAPI(&sdk.PluginAPI{
		Name:      "entry.whitelist",
		Version:   "1.0.0",
		Implement: p,
	})

	// 5. 底层透传：发送一条自定义协议消息
	_ = ctx.Raw("hello", map[string]any{"from": "sunexchange-entry"})

	return nil
}

// Start 启动逻辑。
func (p *Entry) Start() error {
	p.ctx.Successf("进服验证插件已启动")
	_ = p.ctx.Game().SendChat("§aSunExchange 进服验证已就绪")
	return nil
}

// Stop 释放资源。
func (p *Entry) Stop() error {
	p.ctx.Warnf("进服验证插件已停止")
	return nil
}

// GetInfo 返回元信息。
func (p *Entry) GetInfo() sdk.PluginInfo {
	return sdk.PluginInfo{
		Name:        "sunexchange-entry",
		DisplayName: "进服验证插件（全端）",
		Version:     "1.0.0",
		Description: "玩家进服校验并回写事件到鹊桥，支撑群服互通",
		Author:      "SunExchange",
	}
}

// CheckWhitelist 供其他插件调用的进服资格校验。
func (p *Entry) CheckWhitelist(name string) bool {
	return p.whitelist[strings.TrimSpace(name)]
}

var _ sdk.Plugin = (*Entry)(nil)