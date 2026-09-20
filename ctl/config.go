package main

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

const configTOML = "config.toml"

type config struct {
	Username  string
	Password  string
	Cookie    string
	CSRFtoken string
}

// String 故意不打印密码和 Cookie。原版会把明文密码格式化进去，
// 哪天有人 fmt.Println(cfg) 就会把凭据写进终端或 CI 日志。
func (c config) String() string {
	return fmt.Sprintf("Username: %s, Password: <redacted>, Cookie: <redacted>", c.Username)
}

// getConfig 读取 ctl/config.toml。
//
// 这个文件是可选的：没有它也能跑 build readme，只是 LeetCode 接口会以未登录身份返回数据，
// README 里的"个人数据"表格全是 0，"已 AC 但未收录"的列表也会是空的。
// 想要这两块有内容，按 ctl/README.md 的说明配上 Cookie。
// 原版在文件缺失时直接 log.Panic，这里改成返回空配置并提示一句。
func getConfig() *config {
	cfg := new(config)
	if _, err := os.Stat(configTOML); os.IsNotExist(err) {
		fmt.Println("未找到 ctl/config.toml，以未登录身份请求 LeetCode（个人数据统计会是 0）")
		return cfg
	}
	if _, err := toml.DecodeFile(configTOML, cfg); err != nil {
		fmt.Printf("解析 %v 失败: %v，改为以未登录身份请求 LeetCode\n", configTOML, err)
		return new(config)
	}
	return cfg
}
