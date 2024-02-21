// ==============================================================================================================
// 缓存数据key和tag定义.
// 登录用户信息和一些经常要查数据库的数据, 可以存到缓存中, 提升效率.
// 对于同一类数据, 可以添加设置标签方便清除缓存数据.
// 当数据增、删、该时, 注意同步db数据和缓存数据.
// ==============================================================================================================

package global

const (
	// CachePrefix 缓存key前缀
	CachePrefix = "mpuaps_"

	// CacheLoginUserPrefix gtoken登录用户缓存前缀(完整key: CacheLoginUserPrefix + userKey)
	CacheLoginUserPrefix = CachePrefix + "gtoken:"

	// CacheUserNoPassTimePrefix 用户登录验证未通过次数缓存key前缀(完整key: CacheUserNoPassTimePrefix + loginName)
	CacheUserNoPassTimePrefix = CachePrefix + "user_no_pass_"

	// CacheUserLockPrefix 账户/密码错误达最大次数,用户锁定缓存key前缀(完整key: CacheUserLockPrefix + loginName)
	CacheUserLockPrefix = CachePrefix + "user_lock_"

	// CacheSysRole 角色信息缓存key
	CacheSysRole = CachePrefix + "sysRole"

	// CacheMpuChannel 通道信息缓存
)
