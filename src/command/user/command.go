package user

import (
    "encoding/json"
    "errors"
    "fmt"
    "gitee.com/Luna-CY/go-to-internet/src/config"
    "gitee.com/Luna-CY/go-to-internet/src/utils"
    "golang.org/x/crypto/bcrypt"
    "io/ioutil"
    "os"
    "path"
)

// Exec 处理命令
func Exec(config *Config) error {
    cmd := &userCmd{cmdInputConfig: config}
    if err := cmd.validateConfigFilePath(); nil != err {
        return err
    }

    userConfig, err := cmd.load()
    if nil != err {
        return err
    }

    cmd.fileConfig = userConfig

    return cmd.exec()
}

// userCmd 用户子命令结构
type userCmd struct {
    cmdInputConfig *Config
    fileConfig     *config.UserConfig
}

// validateConfigFilePath 检查配置文件
func (u *userCmd) validateConfigFilePath() error {
    if "" == u.cmdInputConfig.Config {
        return errors.New("配置文件路径不能为空")
    }

    u.cmdInputConfig.Config = path.Clean(u.cmdInputConfig.Config)
    info, err := os.Stat(u.cmdInputConfig.Config)
    if nil == err {
        if info.IsDir() {
            return errors.New("用户配置文件路径必须是一个文件")
        }

        return nil
    }

    if !os.IsNotExist(err) {
        return errors.New(fmt.Sprintf("无法查询文件信息: %v", err))
    }

    return u.init()
}

// init 初始化配置文件
func (u *userCmd) init() error {
    file, err := os.OpenFile(u.cmdInputConfig.Config, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
    if nil != err {
        return errors.New(fmt.Sprintf("创建配置文件失败: %v", err))
    }
    defer file.Close()

    n, err := fmt.Fprint(file, config.UserConfigTemplate)
    if nil != err || n != len(config.UserConfigTemplate) {
        return errors.New(fmt.Sprintf("创建配置文件失败: %v", err))
    }

    return nil
}

// load 加载配置文件
func (u *userCmd) load() (*config.UserConfig, error) {
    file, err := os.Open(u.cmdInputConfig.Config)
    if nil != err {
        return nil, errors.New(fmt.Sprintf("无法打开配置文件: %v", err))
    }
    defer file.Close()

    data, err := ioutil.ReadAll(file)
    if nil != err {
        return nil, errors.New(fmt.Sprintf("无法读取配置文件: %v", err))
    }

    userConfig := &config.UserConfig{}
    if err := json.Unmarshal(data, userConfig); nil != err {
        return nil, errors.New(fmt.Sprintf("解析配置文件失败: %v", err))
    }

    return userConfig, nil
}

// exec 执行操作
func (u *userCmd) exec() error {
    update := false

    switch {
    case u.cmdInputConfig.List:
        for key, value := range u.fileConfig.Users {
            fmt.Printf("%v : %v : %d\n", key, value.Expired, value.MaxRate)
        }
    case u.cmdInputConfig.Add:
        if _, ok := u.fileConfig.Users[u.cmdInputConfig.Username]; ok {
            return errors.New("用户名称已存在，无法重复添加")
        }

        password, err := u.password(u.cmdInputConfig.Password)
        if nil != err {
            return err
        }

        userInfo := &config.UserInfo{Password: password, Expired: u.cmdInputConfig.Expired, MaxRate: u.cmdInputConfig.MaxRate}
        u.fileConfig.Users[u.cmdInputConfig.Username] = userInfo

        update = true
    case u.cmdInputConfig.Upd:
        userInfo, ok := u.fileConfig.Users[u.cmdInputConfig.Username]
        if !ok {
            return errors.New("无法找到用户")
        }

        if "" != u.cmdInputConfig.Password {
            password, err := u.password(u.cmdInputConfig.Password)
            if nil != err {
                return errors.New("加密密码失败")
            }

            userInfo.Password = password
        }

        if "" != u.cmdInputConfig.Expired {
            userInfo.Expired = u.cmdInputConfig.Expired
        }

        if -1 != u.cmdInputConfig.MaxRate {
            userInfo.MaxRate = u.cmdInputConfig.MaxRate
        }

        update = true
    case u.cmdInputConfig.Del:
        _, ok := u.fileConfig.Users[u.cmdInputConfig.Username]
        if ok {
            u.fileConfig.Users[u.cmdInputConfig.Username] = nil

            update = true
        }
    }

    if update {
        return u.save()
    }

    return nil
}

// password 加密密码
func (u *userCmd) password(password string) (string, error) {
    password = utils.EncryptPassword(password)

    data, err := bcrypt.GenerateFromPassword([]byte(password), 4)
    if nil != err {
        return "", errors.New(fmt.Sprintf("生成密码失败: %v", err))
    }

    return string(data), nil
}

// save 保存到文件
func (u *userCmd) save() error {
    userConfig, err := u.load()
    if nil != err {
        return err
    }

    if u.fileConfig.Ver != userConfig.Ver {
        return errors.New("文件已被更改，本次更新无法提交，请重新尝试")
    }
    u.fileConfig.Ver += 1

    data, err := json.Marshal(u.fileConfig)
    if nil != err {
        return errors.New(fmt.Sprintf("序列化数据失败: %v", err))
    }

    file, err := os.OpenFile(u.cmdInputConfig.Config, os.O_RDWR|os.O_TRUNC, 0644)
    if nil != err {
        return errors.New(fmt.Sprintf("无法打开配置文件: %v", err))
    }
    defer file.Close()

    n, err := fmt.Fprint(file, string(data))
    if nil != err || n != len(string(data)) {
        return errors.New(fmt.Sprintf("保存数据失败: %v", err))
    }

    return nil
}
