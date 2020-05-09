package config

const UserConfigTemplate = "{\"ver\": 0, \"users\": {}}"

// UserConfig 用户配置文件结构
type UserConfig struct {
    Ver   int                  `json:"ver"`
    Users map[string]*UserInfo `json:"users"`
}

// UserInfo 用户信息结构
type UserInfo struct {
    Password string `json:"password"` // 密码
    Expired  string `json:"expired"`  // 过期时间
    MaxRate  int    `json:"max_rate"` // 最大传输速率
}
