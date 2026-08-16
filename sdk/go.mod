module github.com/suibian-sun/SunExchange-Plugin/sdk

go 1.18

// SunExchange 插件 SDK 独立模块。
//
// 纯标准库实现、无外部依赖，插件可直接 require 本模块开发，
// 并可用较低 Go 版本（1.18+）编译。通信协议由宿主注入，SDK 本身不启动网络。