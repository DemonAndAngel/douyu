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

func Timer() {
	s := helpers.GetScheduler(true)
	limit := viper.GetInt("appConfig.spider.limit")
	if limit <= 0 {
		limit = 20
	}
	fmt.Println("limit", limit)
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
		zb.LastSyncedAt = now
		resp, mode, err := spider.GetH5Play(zb.RoomId, zb.SignMode)
		if err != nil {
			zb.Status = models.StatusFail
			zb.ErrMsg = err.Error()
		} else {
			zb.SignMode = mode
			if resp.Error != 0 {
				zb.Status = models.StatusPending
				zb.ErrMsg = resp.Msg
				if !zb.IsSend {
					zb.IsSend = true
					zb.SendMsg = "停止直播"
					go spider.SendEmail(ctx, zb)
				}
			} else {
				zb.Status = models.StatusNormal
				zb.ErrMsg = ""
				// 发邮件
				// 判断是否需要发邮件
				if zb.IsSend {
					zb.IsSend = false
					zb.SendMsg = "正在直播"
					go spider.SendEmail(ctx, zb)
				}
			}
		}
		_ = zb.Save(ctx)
	}
	fmt.Println("---结束同步数据,当前时间:" + helpers.TimeFormat("Y-m-d H:i:s", now) + "---")
}
