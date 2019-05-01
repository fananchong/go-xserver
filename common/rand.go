package common

// IRand : 随机类接口
type IRand interface {
	// Int63n returns, as an int64, a non-negative pseudo-random number in [0,n).
	// It panics if n <= 0.
	Int63n(n int64) int64

	// Float64 returns, as a float64, a pseudo-random number in [0.0,1.0).
	Float64() float64
}
