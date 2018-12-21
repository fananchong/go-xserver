package custom

// ConfigLogin : 配置 login 节
type ConfigLogin struct {
	Listen string `default:":8080" desc:"登录服务器监听地址"`
	Sign1  string `default:"5UY6$f$h" desc:"用于登录验证的部分签名段"`
	Sign2  string `default:"3wokZB%q" desc:"用于登录验证的部分签名段"`
	Sign3  string `default:"%2Fi9TRf" desc:"用于登录验证的部分签名段"`
}
