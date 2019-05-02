package context

// ITime : 时间类接口
type ITime interface {
	GetTickCount() int64  // 毫秒数
	SetDelta(delta int64) // 设置时间差（单位毫秒），用来快进或后退当前时间（调试时间相关功能时，会用到）
}
