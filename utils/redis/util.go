package redis

import (
	"context"
	"fmt"
	"time"
)

const (
	LIMITER_PREFIX = `seller_limiters`
)

// RedisAppIdLimiter redis分布式限流器
// @param prefix 用于区分请求方
// @param limit 每秒最大请求数
func RedisAppIdLimiter(prefix string, limit uint64) (bool, error) {
	// 获取redis连接
	nowTime := time.Now().Unix()
	redisKey := LIMITER_PREFIX + ":app_id-limiter:" + fmt.Sprintf("%s_%d", prefix, nowTime)
	num, err := GetDefaultClient().GetClient().Incr(context.Background(), redisKey).Uint64()
	if err != nil {
		return false, err
	}
	if num == 1 {
		// 首次需要设置过期时间
		GetDefaultClient().GetClient().Expire(context.Background(), redisKey, time.Second*5)
	}
	// 当前请求次数大于限制次数
	if num > limit {
		// 这里回滚请求次数
		GetDefaultClient().GetClient().Decr(context.Background(), redisKey)
		return false, nil
	}
	return true, nil
}
