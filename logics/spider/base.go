package spider

import (
	"context"
	"douyu/models"
	"douyu/utils/email"
	"fmt"
	"github.com/spf13/viper"
)

var E *email.Email

func init() {
	E = email.NewEmail("斗鱼推送", "568089002@qq.com", "mjqdlojunzvnbfii", "smtp.qq.com:465")
}

func SendEmail(ctx context.Context, zb models.ZhuBo) {
	tos := viper.GetStringSlice("appConfig.spider.etos")
	if len(tos) <= 0 {
		tos = []string{"568089002@qq.com"}
	}
	subject, body := genHtml(zb)
	fmt.Println(subject, body)
	err := E.Send(ctx, tos, "html", subject, body)
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
