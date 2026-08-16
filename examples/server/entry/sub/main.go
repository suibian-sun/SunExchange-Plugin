// Package main 将 entry 示例插件编译为可动态加载的 .so（Go plugin）。
//
// 编译（在 examples/server/entry 目录下执行）：
//
//	go build -buildmode=plugin -o entry.so ./sub
//
// 产物 entry.so 放入插件目录即可被 SunExchange 宿主动态加载。
// 完整 SDK 开发方法见 docs/ 下的指南。
package main

import (
	"github.com/suibian-sun/SunExchange-Plugin/sdk"
	"example.com/sunexchange/entry"
)

// Plugin 是 .so 必须导出的符号，返回 sdk.Plugin 实现。
func Plugin() sdk.Plugin {
	return entry.NewEntry()
}

// main 空函数：Go 插件包要求包含 main。
func main() {}