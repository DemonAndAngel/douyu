package timer

import (
	"context"
	"douyu/logics/spider"
	"douyu/models"
	"douyu/utils/helpers"
	"fmt"
	"github.com/spf13/viper"
	"time"
)

const SendEmailModeTitle = "title"
const SendEmailModeLive = "live"

func Timer() {
	// 初始化邮件驱动
	InitEmail()
	s := helpers.GetScheduler(true)
	limit := viper.GetInt("appConfig.spider.limit")
	if limit <= 0 {
		limit = 20
	}
	_, err := s.Every(limit).Second().Do(spiderHandler)
	if err != nil {
		panic("调度器启动失败;信息:" + err.Error())
	}
}

func spiderHandler() {
	ctx := context.Background()
	now := time.Now()
	fmt.Println("---开始同步数据,当前时间:" + helpers.TimeFormat("Y-m-d H:i:s", now) + "---")
	// 找出所有数据
	zbs := models.ZhuBoGetToList(ctx, map[string]interface{}{})
	for _, zb := range zbs {
		_ = spiderHandlerRun(ctx, zb)
	}
	fmt.Println("---结束同步数据,当前时间:" + helpers.TimeFormat("Y-m-d H:i:s", now) + "---")
}

func spiderHandlerRun(ctx context.Context, zb models.ZhuBo) models.ZhuBo {
	zb.LastSyncedAt = time.Now()
	defer func() {
		_ = zb.Save(ctx)
	}()
	iResp, err := spider.GetIndex(zb.RoomId)
	if err != nil {
		zb.Status = models.StatusFail
		zb.ErrMsg = err.Error()
		return zb
	}
	if zb.Name != iResp.Name {
		zb.Name = iResp.Name
	}
	if zb.Title != iResp.Title {
		// 更换标题 发送邮件
		zb.Title = iResp.Title
		go sendEmail(ctx, zb, SendEmailModeTitle)
	}
	hResp, mode, err := spider.GetH5Play(zb.RoomId, zb.SignMode)
	if err != nil {
		zb.Status = models.StatusFail
		zb.ErrMsg = err.Error()
		return zb
	} else {
		zb.SignMode = mode
		if hResp.Error != 0 {
			zb.Status = models.StatusPending
			zb.ErrMsg = hResp.Msg
			if !zb.IsSend {
				zb.IsSend = true
				zb.SendMsg = "停止直播"
				go sendEmail(ctx, zb, SendEmailModeLive)
			}
		} else {
			zb.Status = models.StatusNormal
			zb.ErrMsg = ""
			// 发邮件
			// 判断是否需要发邮件
			if zb.IsSend {
				zb.IsSend = false
				zb.SendMsg = "正在直播"
				go sendEmail(ctx, zb, SendEmailModeLive)
			}
		}
		return zb
	}
}

func sendEmail(ctx context.Context, zb models.ZhuBo, mode string) {
	tos := viper.GetStringSlice("appConfig.spider.etos")
	if len(tos) <= 0 {
		tos = []string{"568089002@qq.com"}
	}
	subject := ""
	body := ""
	switch mode {
	case SendEmailModeLive:
		subject = fmt.Sprintf("%s-%s", zb.Name, zb.SendMsg)
		body = genBodyHtml(zb.Name, subject)
		break
	case SendEmailModeTitle:
		subject = fmt.Sprintf("%s-%s-%s", zb.Name, "更换标题", zb.Title)
		body = genBodyHtml(zb.Name, subject)
		break
	}
	fmt.Println(subject, body)
	err := e.Send(ctx, tos, "html", subject, body)
	if err != nil {
		fmt.Println(fmt.Sprintf("发送邮件异常,信息:%s", err.Error()))
	}
}

// 生成邮件html
func genBodyHtml(title, body string) (html string) {
	//subject = fmt.Sprintf("%s-%s", zb.Name, zb.SendMsg)
	html = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>%s</title>
</head>
<body>
    %s
</body>
</html>
`
	html = fmt.Sprintf(html, title, body)
	return
}
