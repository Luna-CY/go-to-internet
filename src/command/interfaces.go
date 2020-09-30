package command

// Config 子命令配置接口
type Config interface {

    // Usage 方法返回该配置的使用文档
    Usage()

    // Validate 方法校验用户输入的配置是否满足子命令的运行
    Validate() bool
}

// SubCommand 子命令接口
type SubCommand interface {

    // 执行命令
    Exec() error
}
