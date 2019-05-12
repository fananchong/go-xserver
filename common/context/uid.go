package context

// IUID : 唯一ID接口
type IUID interface {
	GetUID(key string) (uint64, error) // 根据 KEY , 获取 UID
}
