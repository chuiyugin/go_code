package auth

import (
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/redis"
)

func rtAllowKey(uid int64, jti string) string {
	return fmt.Sprintf("RT:allow:%d:%s", uid, jti)
}

// ttlSeconds = RefreshExpire
func RTAllow(rds *redis.Redis, uid int64, jti string, ttlSeconds int64) error {
	return rds.Setex(rtAllowKey(uid, jti), "1", int(ttlSeconds))
}

func RTIsAllowed(rds *redis.Redis, uid int64, jti string) (bool, error) {
	v, err := rds.Get(rtAllowKey(uid, jti))
	if err != nil && err != redis.Nil {
		return false, err
	}
	return v == "1", nil
}

func RTRevoke(rds *redis.Redis, uid int64, jti string) error {
	_, err := rds.Del(rtAllowKey(uid, jti))
	return err
}
