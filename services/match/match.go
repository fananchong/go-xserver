package main

// Match : 匹配服务器
type Match struct {
}

// NewMatch : 构造函数
func NewMatch() *Match {
	match := &Match{}
	return match
}

// Start : 启动
func (match *Match) Start() bool {
	return true
}

// Close : 关闭
func (match *Match) Close() {
}
