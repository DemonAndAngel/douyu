package models

import (
	"context"
	"douyu/utils/db"
	"gorm.io/gorm"
	"time"
)

const (
	StatusSubmitted = `submitted`
	StatusPending   = `pending`
	StatusNormal    = `normal`
	StatusFail      = `fail`
)

type ZhuBo struct {
	Id           uint64    `json:"id" gorm:"primary_key;autoIncrement:false"`
	RoomId       string    `json:"roomId" gorm:"index;not null;size:50;comment:斗鱼房间号"`
	Name         string    `json:"name" grom:"not null;size:200;comment:主播名"`
	Alias        string    `json:"alias" gorm:"not null;size:200;comment:别名"`
	Title        string    `json:"title" gorm:"not null;size:500;comment:标题"`
	SignMode     string    `json:"signMode" gorm:"not null;size:30;comment:签名模式"`
	Status       string    `json:"status" gorm:"not null;size:30;comment:状态"`
	ErrMsg       string    `json:"errMsg" gorm:"not null;size:1024;comment:异常信息"`
	SendMsg      string    `json:"sendMsg" gorm:"not null;size:1024;comment:发送的消息"`
	IsSend       bool      `json:"isSend" gorm:"not null;comment:是否需要通知"`
	LastSyncedAt time.Time `json:"lastSyncedAt" gorm:"not null;comment:最后一次同步时间"`
	CreatedAt    time.Time `json:"createdAt" gorm:"not null"`
	UpdatedAt    time.Time `json:"updatedAt" gorm:"not null"`
}

func ZhuBoGetToSession(ctx context.Context, search map[string]interface{}) *gorm.DB {
	session := db.GetDefaultClient().GetDB().WithContext(ctx).Model(&ZhuBo{})
	for k, v := range search {
		switch k {
		case "id":
			if v.(uint64) != 0 {
				session = session.Where("id = ?", v)
			}
			break
		case "roomId":
			if v.(string) != "" {
				session = session.Where("room_id = ?", v)
			}
			break
		}
	}
	return session
}

func ZhuBoGetToList(ctx context.Context, search map[string]interface{}) (zbs []ZhuBo) {
	session := ZhuBoGetToSession(ctx, search)
	session.Order("id desc").Find(&zbs)
	return
}

func (m *ZhuBo) Save(ctx context.Context) error {
	session := db.GetDefaultClient().GetDB().WithContext(ctx)
	return session.Save(m).Error
}
