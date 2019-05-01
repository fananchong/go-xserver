package context

// IRole2Account : 提供`根据角色名查找账号`的功能；角色名重名检查也可以用该接口
type IRole2Account interface {
	Add(role, account string)                 // 加入本地缓存
	AddAndInsertDB(role, account string) bool // 加入本地缓存并保存数据库
	GetAndActive(role string) string          // 根据角色名，查找账号。先从本地缓存中找，没有则数据库中找
}
