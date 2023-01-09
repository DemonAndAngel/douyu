package timer

import (
	"context"
	"douyu/models"
	"douyu/utils/email"
	"fmt"
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

func SendEmail(ctx context.Context, zb models.ZhuBo) {
	tos := viper.GetStringSlice("appConfig.spider.etos")
	if len(tos) <= 0 {
		tos = []string{"568089002@qq.com"}
	}
	subject, body := genHtml(zb)
	fmt.Println(subject, body)
	err := e.Send(ctx, tos, "html", subject, body)
	if err != nil {
		fmt.Println(fmt.Sprintf("发送邮件异常,信息:%s", err.Error()))
	}
}

// 生成邮件html
func genHtml(zb models.ZhuBo) (subject, body string) {
	subject = fmt.Sprintf("%s-%s", zb.Name, zb.SendMsg)
	body = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>%s</title>
</head>
<body>
    %s-%s
</body>
</html>
`
	body = fmt.Sprintf(body, zb.Name, zb.Name, zb.SendMsg)
	return
}
