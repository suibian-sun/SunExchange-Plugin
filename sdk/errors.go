package sdk

import (
	"errors"
	"reflect"
)

// SDK 错误定义。
var (
	ErrEmptyCommandName = errors.New("命令名不能为空")
	ErrNoBroker         = errors.New("宿主未提供插件间通信能力")
	ErrNoTransport      = errors.New("宿主未提供底层透传通道")
	ErrTimeout          = errors.New("等待超时")
)

// reflectValueOf 取函数指针（用于处理器去重比较）。
func reflectValueOf(h Handler) uintptr {
	return reflect.ValueOf(h).Pointer()
}