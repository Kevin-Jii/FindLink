package redis

import (
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis"

	"app/service/dto"
)

const (
	// PubSub channels
	locationUpdateChannel = "channel:location:update"
	wsMessageChannel      = "channel:ws:message"
)

// PubSub 发布订阅服务
type PubSub struct {
	client *redis.Client
}

// NewPubSub 创建发布订阅服务
func NewPubSub(client *redis.Client) *PubSub {
	return &PubSub{client: client}
}

// PublishLocationUpdate 发布位置更新
func (p *PubSub) PublishLocationUpdate(entityType string, entityID string, loc *dto.LocationResp) error {
	msg := map[string]interface{}{
		"type":        "location_update",
		"entity_type": entityType,
		"entity_id":   entityID,
		"location":    loc,
	}
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return p.client.Publish(locationUpdateChannel, data).Err()
}

// PublishWSMessage 发布WebSocket消息
func (p *PubSub) PublishWSMessage(userID int64, msg interface{}) error {
	key := fmt.Sprintf("%s:%d", wsMessageChannel, userID)
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return p.client.Publish(key, data).Err()
}

// Subscribe 订阅频道
func (p *PubSub) Subscribe(channels ...string) *redis.PubSub {
	return p.client.Subscribe(channels...)
}

// Publish 通用发布
func (p *PubSub) Publish(channel string, message interface{}) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}
	return p.client.Publish(channel, data).Err()
}
