package redis

import (
	"app/config"
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

type IVerify interface {
	SetCaptchaKey(ctx context.Context, key string, value string, expire time.Duration) error
	GetCaptchaKey(ctx context.Context, key string) (string, error)
	SetCaptchaTicket(ctx context.Context, key string, value string, expire time.Duration) error
	GetCaptchaTicket(ctx context.Context, key string) (string, error)
	SetToken(ctx context.Context, token string, userID int64, expire time.Duration) error
	GetToken(ctx context.Context, token string) (int64, error)
	DelToken(ctx context.Context, token string) error
}

type Verify struct {
	redis *redis.Client
}

func NewVerify(redisClient *redis.Client) *Verify {
	return &Verify{
		redis: redisClient,
	}
}

func fmtVerifyCaptchaKey(key string) string {
	return fmt.Sprintf("%s:captcha:%s", config.ServerFullName, key)
}

func fmtVerifyCaptchaTicket(key string) string {
	return fmt.Sprintf("%s:captcha:ticket:%s", config.ServerFullName, key)
}
func (v *Verify) SetCaptchaKey(ctx context.Context, key string, value string, expire time.Duration) error {
	redisKey := fmtVerifyCaptchaKey(key)
	return v.redis.Set(redisKey, value, expire).Err()
}
func (v *Verify) GetCaptchaKey(ctx context.Context, key string) (string, error) {
	redisKey := fmtVerifyCaptchaKey(key)
	get, err := v.redis.Get(redisKey).Result()
	if err != nil {
		return "", err
	}
	v.redis.Del(redisKey)
	return get, nil
}
func (v *Verify) SetCaptchaTicket(ctx context.Context, key string, value string, expire time.Duration) error {
	redisKey := fmtVerifyCaptchaTicket(key)
	return v.redis.Set(redisKey, value, expire).Err()
}
func (v *Verify) GetCaptchaTicket(ctx context.Context, key string) (string, error) {
	redisKey := fmtVerifyCaptchaTicket(key)
	get, err := v.redis.Get(redisKey).Result()
	if err != nil {
		return "", err
	}
	v.redis.Del(redisKey)
	return get, nil
}

func fmtTokenKey(token string) string {
	return fmt.Sprintf("%s:token:%s", config.ServerFullName, token)
}

func (v *Verify) SetToken(ctx context.Context, token string, userID int64, expire time.Duration) error {
	redisKey := fmtTokenKey(token)
	return v.redis.Set(redisKey, userID, expire).Err()
}

func (v *Verify) GetToken(ctx context.Context, token string) (int64, error) {
	redisKey := fmtTokenKey(token)
	return v.redis.Get(redisKey).Int64()
}

func (v *Verify) DelToken(ctx context.Context, token string) error {
	redisKey := fmtTokenKey(token)
	return v.redis.Del(redisKey).Err()
}
