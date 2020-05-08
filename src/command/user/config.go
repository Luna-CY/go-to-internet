package user

// Config 用户配置结构
type Config struct {
    List bool // 打印用户列表
    Add  bool // 添加用户
    Upd  bool // 更新用户
    Del  bool // 删除用户

    Username string // 用户名
    Password string // 密码
    Expired  string // 有效期，格式: yyyy-MM-dd HH:mm:ss
    MaxRate  int64  // 最大速率，单位kb，格式: 1024 * 1024
}
