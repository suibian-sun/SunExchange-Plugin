module example.com/sunexchange/entry

go 1.18

// 独立示例插件模块：演示 SunExchange 插件 SDK 的完整能力。
// 开发期用 replace 指向本仓库根下 sdk；发布后替换为版本 tag。

require github.com/suibian-sun/SunExchange-Plugin/sdk v0.0.0

replace github.com/suibian-sun/SunExchange-Plugin/sdk => ../../../sdk
