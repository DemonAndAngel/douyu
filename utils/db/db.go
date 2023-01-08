package db

import (
	"douyu/utils/helpers"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"time"
)

// 数据库配置结构体
type Config struct {
	Username        string `json:"username"`
	Host            string `json:"host"`
	Port            int    `json:"port"`
	Password        string `json:"password"`
	Database        string `json:"database"`
	Prefix          string `json:"prefix"` // 表名前缀
	ConnMaxLifetime int    `json:"connMaxLifetime"`
}

type Client struct {
	db     *gorm.DB
	config *Config
	driver string
}

var clients map[string]*Client

var defaultDriver string

func InitWithEnv() (err error) {
	clients = make(map[string]*Client)
	// 默认驱动
	defaultDriver = viper.GetString("driver.mysql")
	// 取出所有redis配置
	configs := viper.GetStringMap("mysql")
	for k, c := range configs {
		// 生成所有连接
		dc := &Config{}
		_ = json.Unmarshal([]byte(helpers.JsonMarshal(c)), &dc)
		db, _err := ConnectMysqlWithGorm(dc)
		if _err != nil {
			return _err
		}
		clients[k] = &Client{
			db:     db,
			config: dc,
			driver: k,
		}
	}
	return nil
}

func GetDefaultClient() *Client {
	return GetClient(defaultDriver)
}
func GetClient(driver string) *Client {
	if client, ok := clients[driver]; ok {
		return client
	} else {
		panic("client not found")
	}
}
func GetClients() map[string]*Client {
	return clients
}

func (c *Client) GetDB() *gorm.DB {
	if viper.GetString("appEnv") == "prod" {
		return c.db
	} else {
		return c.db.Debug()
	}
}
func (c *Client) GetConfig() *Config {
	return c.config
}
func (c *Client) GetDriver() string {
	return c.driver
}

// 获取数据库连接
func ConnectMysqlWithGorm(tmpConf *Config) (db *gorm.DB, err error) {
	str := "%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True"
	db, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                      fmt.Sprintf(str, tmpConf.Username, tmpConf.Password, tmpConf.Host, tmpConf.Port, tmpConf.Database) + "&loc=Asia%2FShanghai", // DSN data source name
		DisableDatetimePrecision: true,                                                                                                                        // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:   true,                                                                                                                        // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		// DontSupportRenameColumn: true, // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据版本自动配置
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: tmpConf.Prefix,
		},
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return
	}
	// 获取通用数据库对象 sql.DB ，然后使用其提供的功能
	sqlDB, err := db.DB()
	if err != nil {
		return
	}
	// SetMaxIdleConns 用于设置连接池中空闲连接的最大数量。
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(1000)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Duration(tmpConf.ConnMaxLifetime) * time.Second)
	return
}

func GetDb() (db *gorm.DB, err error) {
	str := viper.GetStringMap("mysql")
	mysqlCfg := str["default"]
	dc := &Config{}
	_ = json.Unmarshal([]byte(helpers.JsonMarshal(mysqlCfg)), &dc)
	gdb, err := ConnectMysqlWithGorm(dc)
	if err != nil {
		panic(err)
	}
	return gdb, err
}
