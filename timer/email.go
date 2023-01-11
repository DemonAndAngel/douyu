package timer

import (
	"douyu/utils/email"
	"github.com/spf13/viper"
)

var e *email.Email

// InitEmail 初始化邮件驱动
func InitEmail() {
	sendUser := viper.GetString("appConfig.spider.email.sendUser")
	sendEmailUser := viper.GetString("appConfig.spider.email.sendEmailUser")
	sendEmailPass := viper.GetString("appConfig.spider.email.sendEmailPass")
	sendHost := viper.GetString("appConfig.spider.email.sendHost")
	if sendUser == "" || sendHost == "" || sendEmailUser == "" || sendEmailPass == "" {
		panic("邮件配置异常,请检查")
	}
	e = email.NewEmail(sendUser, sendEmailUser, sendEmailPass, sendHost)
}
