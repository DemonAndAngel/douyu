package models

import (
	"douyu/utils/db"
)

// Migrate 数据表迁移
func Migrate() {
	if err := db.GetDefaultClient().GetDB().Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").
		AutoMigrate(ZhuBo{}); err != nil {
		panic(err.Error())
	}
}
