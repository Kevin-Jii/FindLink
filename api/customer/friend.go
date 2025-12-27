package customer

import (
	"github.com/gin-gonic/gin"

	"app/api"
	"app/common"
	"app/service/dto"
	"app/service/friend"
)

// FriendCtrl 好友控制器 - 嵌入到 Ctrl 中
type FriendCtrl struct {
	Friend *friend.FriendService
}

// NewFriendCtrl 创建好友控制器
func NewFriendCtrl(friend *friend.FriendService) *FriendCtrl {
	return &FriendCtrl{Friend: friend}
}

// @Summary 发送好友请求
// @Description 发送好友请求给指定用户
// @Tags friend
// @Accept json
// @Produce json
// @Param Authorization header string true "Token"
// @Param req body dto.FriendRequestReq true "好友请求"
// @Success 200 {object} api.Resp
// @Router /api/app/customer/v1/friend/request [post]
func (c *Ctrl) SendFriendRequest(ctx *gin.Context) {
	req := &dto.FriendRequestReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, common.ParamErr.WithErr(err))
		return
	}

	userID := getUserID(ctx)
	if err := c.Friend.SendFriendRequest(ctx.Request.Context(), userID, req.ToUserID, req.Message); err != nil {
		api.WriteResp(ctx, nil, err.(common.Errno))
		return
	}

	api.WriteResp(ctx, nil, common.OK)
}

// @Summary 获取好友请求列表
// @Description 获取待处理的好友请求
// @Tags friend
// @Produce json
// @Param Authorization header string true "Token"
// @Success 200 {object} api.Resp{data=[]dto.FriendRequestResp}
// @Router /api/app/customer/v1/friend/requests [get]
func (c *Ctrl) GetFriendRequests(ctx *gin.Context) {
	userID := getUserID(ctx)

	requests, err := c.Friend.GetFriendRequests(ctx.Request.Context(), userID)
	if err != nil {
		api.WriteResp(ctx, nil, err.(common.Errno))
		return
	}

	api.WriteResp(ctx, requests, common.OK)
}

// @Summary 接受好友请求
// @Description 接受指定的好友请求
// @Tags friend
// @Accept json
// @Produce json
// @Param Authorization header string true "Token"
// @Param req body dto.FriendRequestActionReq true "请求ID"
// @Success 200 {object} api.Resp
// @Router /api/app/customer/v1/friend/accept [post]
func (c *Ctrl) AcceptFriendRequest(ctx *gin.Context) {
	req := &dto.FriendRequestActionReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, common.ParamErr.WithErr(err))
		return
	}

	userID := getUserID(ctx)
	if err := c.Friend.AcceptFriendRequest(ctx.Request.Context(), userID, req.RequestID); err != nil {
		api.WriteResp(ctx, nil, err.(common.Errno))
		return
	}

	api.WriteResp(ctx, nil, common.OK)
}

// @Summary 拒绝好友请求
// @Description 拒绝指定的好友请求
// @Tags friend
// @Accept json
// @Produce json
// @Param Authorization header string true "Token"
// @Param req body dto.FriendRequestActionReq true "请求ID"
// @Success 200 {object} api.Resp
// @Router /api/app/customer/v1/friend/reject [post]
func (c *Ctrl) RejectFriendRequest(ctx *gin.Context) {
	req := &dto.FriendRequestActionReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, common.ParamErr.WithErr(err))
		return
	}

	userID := getUserID(ctx)
	if err := c.Friend.RejectFriendRequest(ctx.Request.Context(), userID, req.RequestID); err != nil {
		api.WriteResp(ctx, nil, err.(common.Errno))
		return
	}

	api.WriteResp(ctx, nil, common.OK)
}

// @Summary 获取好友列表
// @Description 获取已接受的好友列表
// @Tags friend
// @Produce json
// @Param Authorization header string true "Token"
// @Success 200 {object} api.Resp{data=[]dto.FriendResp}
// @Router /api/app/customer/v1/friend/list [get]
func (c *Ctrl) GetFriendList(ctx *gin.Context) {
	userID := getUserID(ctx)

	friends, err := c.Friend.GetFriendList(ctx.Request.Context(), userID)
	if err != nil {
		api.WriteResp(ctx, nil, err.(common.Errno))
		return
	}

	api.WriteResp(ctx, friends, common.OK)
}

// @Summary 删除好友
// @Description 删除指定好友
// @Tags friend
// @Produce json
// @Param Authorization header string true "Token"
// @Param friend_id path int true "好友ID"
// @Success 200 {object} api.Resp
// @Router /api/app/customer/v1/friend/{friend_id} [delete]
func (c *Ctrl) RemoveFriend(ctx *gin.Context) {
	userID := getUserID(ctx)
	friendID := parseInt64(ctx.Param("friend_id"))

	if err := c.Friend.RemoveFriend(ctx.Request.Context(), userID, friendID); err != nil {
		api.WriteResp(ctx, nil, err.(common.Errno))
		return
	}

	api.WriteResp(ctx, nil, common.OK)
}

// @Summary 搜索用户
// @Description 搜索用户
// @Tags friend
// @Produce json
// @Param Authorization header string true "Token"
// @Param keyword query string true "搜索关键词"
// @Success 200 {object} api.Resp{data=[]dto.UserSearchResp}
// @Router /api/app/customer/v1/friend/search [get]
func (c *Ctrl) SearchUsers(ctx *gin.Context) {
	userID := getUserID(ctx)
	keyword := ctx.Query("keyword")

	users, err := c.Friend.SearchUsers(ctx.Request.Context(), userID, keyword)
	if err != nil {
		api.WriteResp(ctx, nil, err.(common.Errno))
		return
	}

	api.WriteResp(ctx, users, common.OK)
}
