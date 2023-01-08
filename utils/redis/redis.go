package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"douyu/utils/helpers"
)

type Config struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database int    `json:"database"`
}
type Client struct {
	client *redis.Client
	config *Config
	driver string
}

type ClientMap map[string]*Client

var clients ClientMap

var defaultDriver string // 默认驱动

func InitWithEnv() (err error) {
	clients = make(map[string]*Client)
	// 默认驱动
	defaultDriver = viper.GetString("driver.redis")
	ctx := context.Background()
	// 取出所有redis配置
	configs := viper.GetStringMap("redis")
	for k, c := range configs {
		// 生成所有连接
		rc := &Config{}
		_ = json.Unmarshal([]byte(helpers.JsonMarshal(c)), &rc)
		client, _err := ConnectRedisBase(rc)
		if _err != nil {
			return _err
		}
		_err = client.Ping(ctx).Err()
		if _err != nil {
			return _err
		}
		clients[k] = &Client{
			client: client,
			config: rc,
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
func GetClients() ClientMap {
	return clients
}

func (c *Client) GetClient() *redis.Client {
	return c.client
}
func (c *Client) GetConfig() *Config {
	return c.config
}
func (c *Client) GetDriver() string {
	return c.driver
}

func (m ClientMap) CloseAll() error {
	for k, _ := range m {
		err := m.Close(k)
		if err != nil {
			return err
		}
	}
	return nil
}
func (m ClientMap) Close(driver string) error {
	if c, ok := m[driver]; ok {
		err := c.GetClient().Close()
		if err != nil {
			return err
		}
		delete(m, driver)
	}
	return nil
}

func ConnectRedisBase(tmpConf *Config) (*redis.Client, error) {
	if tmpConf.Host == "" {
		return nil, errors.New("host is not nil")
	}
	c := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", tmpConf.Host, tmpConf.Port),
		Username:     tmpConf.Username,
		Password:     tmpConf.Password,
		DB:           tmpConf.Database,
		MinIdleConns: 5,
		OnConnect: func(ctx context.Context, cn *redis.Conn) error {
			fmt.Println(fmt.Sprintf("连接redis:%v", cn))
			return nil
		},
	})
	err := c.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}
	return c, nil
}
