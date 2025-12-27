package redis

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

// Cache 通用缓存服务
type Cache struct {
	client *redis.Client
}

// NewCache 创建通用缓存服务
func NewCache(client *redis.Client) *Cache {
	return &Cache{client: client}
}

// Set 设置缓存
func (c *Cache) Set(key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.client.Set(key, data, ttl).Err()
}

// Get 获取缓存
func (c *Cache) Get(key string) ([]byte, error) {
	data, err := c.client.Get(key).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	return data, err
}

// GetStruct 获取缓存并反序列化
func (c *Cache) GetStruct(key string, dest interface{}) error {
	data, err := c.Get(key)
	if err != nil || data == nil {
		return err
	}
	return json.Unmarshal(data, dest)
}

// Delete 删除缓存
func (c *Cache) Delete(keys ...string) error {
	if len(keys) == 0 {
		return nil
	}
	return c.client.Del(keys...).Err()
}

// SetFriendsCache 设置好友列表缓存
func (c *Cache) SetFriendsCache(userID int64, friends interface{}) error {
	key := fmt.Sprintf(friendsKey, userID)
	return c.Set(key, friends, friendsTTL)
}

// GetFriendsCache 获取好友列表缓存
func (c *Cache) GetFriendsCache(userID int64, dest interface{}) error {
	key := fmt.Sprintf(friendsKey, userID)
	return c.GetStruct(key, dest)
}

// DeleteFriendsCache 删除好友列表缓存
func (c *Cache) DeleteFriendsCache(userID int64) error {
	key := fmt.Sprintf(friendsKey, userID)
	return c.Delete(key)
}
