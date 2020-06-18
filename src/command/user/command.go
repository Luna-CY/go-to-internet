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

// UserCmd 用户子命令结构
type Cmd struct {
    Config *Config

    userConfig *config.UserConfig
}

func (c *Cmd) Exec() error {
    if err := c.validateConfigFilePath(); nil != err {
        return err
    }

    userConfig, err := c.load()
    if nil != err {
        return err
    }

    c.userConfig = userConfig

    return c.exec()
}

// validateConfigFilePath 检查配置文件
func (c *Cmd) validateConfigFilePath() error {
    if "" == c.Config.Config {
        return errors.New("配置文件路径不能为空")
    }

    c.Config.Config = path.Clean(c.Config.Config)
    info, err := os.Stat(c.Config.Config)
    if nil == err {
        if info.IsDir() {
            return errors.New("用户配置文件路径必须是一个文件")
        }

        return nil
    }

    if !os.IsNotExist(err) {
        return errors.New(fmt.Sprintf("无法查询文件信息: %v", err))
    }

    return c.init()
}

// init 初始化配置文件
func (c *Cmd) init() error {
    if err := os.MkdirAll(path.Dir(c.Config.Config), os.FileMode(0644)); nil != err {
        return errors.New(fmt.Sprintf("创建配置目录失败: %v", err))
    }

    file, err := os.OpenFile(c.Config.Config, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
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
func (c *Cmd) load() (*config.UserConfig, error) {
    file, err := os.Open(c.Config.Config)
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
func (c *Cmd) exec() error {
    update := false

    switch {
    case c.Config.List:
        for key, value := range c.userConfig.Users {
            fmt.Printf("%v : %v : %d\n", key, value.Expired, value.MaxRate)
        }
    case c.Config.Add:
        if _, ok := c.userConfig.Users[c.Config.Username]; ok {
            return errors.New("用户名称已存在，无法重复添加")
        }

        password, err := c.password(c.Config.Password)
        if nil != err {
            return err
        }

        expired := "-"
        maxRate := 0
        maxConnection := 0

        if "" != c.Config.Expired {
            expired = c.Config.Expired
        }
        if 0 < c.Config.MaxRate {
            maxRate = c.Config.MaxRate
        }
        if 0 < c.Config.MaxConnection {
            maxConnection = c.Config.MaxConnection
        }

        userInfo := &config.UserInfo{Password: password, Expired: expired, MaxRate: maxRate, MaxConnection: maxConnection}
        c.userConfig.Users[c.Config.Username] = userInfo

        update = true
    case c.Config.Upd:
        userInfo, ok := c.userConfig.Users[c.Config.Username]
        if !ok {
            return errors.New("无法找到用户")
        }

        if "" != c.Config.Password {
            password, err := c.password(c.Config.Password)
            if nil != err {
                return errors.New("加密密码失败")
            }

            userInfo.Password = password
        }

        if "" != c.Config.Expired {
            userInfo.Expired = c.Config.Expired
        }

        if 0 <= c.Config.MaxRate {
            userInfo.MaxRate = c.Config.MaxRate
        }

        if 0 <= c.Config.MaxConnection {
            userInfo.MaxConnection = c.Config.MaxConnection
        }

        update = true
    case c.Config.Del:
        _, ok := c.userConfig.Users[c.Config.Username]
        if ok {
            delete(c.userConfig.Users, c.Config.Username)

            update = true
        }
    }

    if update {
        return c.save()
    }

    return nil
}

// password 加密密码
func (c *Cmd) password(password string) (string, error) {
    password = utils.EncryptPassword(password)

    data, err := bcrypt.GenerateFromPassword([]byte(password), 4)
    if nil != err {
        return "", errors.New(fmt.Sprintf("生成密码失败: %v", err))
    }

    return string(data), nil
}

// save 保存到文件
func (c *Cmd) save() error {
    userConfig, err := c.load()
    if nil != err {
        return err
    }

    if c.userConfig.Ver != userConfig.Ver {
        return errors.New("文件已被更改，本次更新无法提交，请重新尝试")
    }
    c.userConfig.Ver += 1

    data, err := json.Marshal(c.userConfig)
    if nil != err {
        return errors.New(fmt.Sprintf("序列化数据失败: %v", err))
    }

    file, err := os.OpenFile(c.Config.Config, os.O_RDWR|os.O_TRUNC, 0644)
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
