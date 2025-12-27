package customer

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	gorillawebsocket "github.com/gorilla/websocket"

	"app/api"
	"app/common"
	ws "app/service/websocket"
)

var upgrader = gorillawebsocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // 生产环境应该检查Origin
	},
}

// @Summary WebSocket连接
// @Description 建立WebSocket实时连接
// @Tags websocket
// @Produce json
// @Param Authorization header string true "Token"
// @Success 101 {object} api.Resp "Switching Protocols"
// @Router /api/app/customer/v1/ws [get]
func (c *Ctrl) WebSocketConnect(ctx *gin.Context) {
	// 从Context获取用户ID（通过JWT中间件设置）
	userID := getUserID(ctx)
	if userID == 0 {
		api.WriteResp(ctx, nil, common.WSAuthFailedErr)
		return
	}

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return
	}

	client := ws.NewClient(
		time.Now().Format("20060102150405.000000"),
		userID,
		c.Hub,
		conn,
	)

	c.Hub.Register(client)

	go client.WritePump()
	go client.ReadPump(handleWSMessage)
}

func handleWSMessage(msg *ws.Message) error {
	switch msg.Type {
	case "subscribe":
		// 处理订阅逻辑
		var payload struct {
			Entities []struct {
				Type string `json:"type"`
				ID   string `json:"id"`
			} `json:"entities"`
		}
		if err := json.Unmarshal(msg.Payload, &payload); err != nil {
			return err
		}
		// TODO: 实现订阅管理

	case "unsubscribe":
		// 处理取消订阅逻辑
		var payload struct {
			Entities []struct {
				Type string `json:"type"`
				ID   string `json:"id"`
			} `json:"entities"`
		}
		if err := json.Unmarshal(msg.Payload, &payload); err != nil {
			return err
		}
		// TODO: 实现取消订阅管理

	case "ping":
		// 心跳处理
	}

	return nil
}
