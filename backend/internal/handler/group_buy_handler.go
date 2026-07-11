package handler

import (
	"strconv"

	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	middleware2 "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
)

type GroupBuyHandler struct {
	groupBuyService *service.GroupBuyService
}

type CreateGroupBuyRoomRequest struct {
	PlanID      int64 `json:"plan_id" binding:"required"`
	TargetCount int   `json:"target_count" binding:"required"`
}

func NewGroupBuyHandler(groupBuyService *service.GroupBuyService) *GroupBuyHandler {
	return &GroupBuyHandler{groupBuyService: groupBuyService}
}

func (h *GroupBuyHandler) List(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not found in context")
		return
	}
	hall, err := h.groupBuyService.ListHall(c.Request.Context(), subject.UserID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, hall)
}

func (h *GroupBuyHandler) Create(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not found in context")
		return
	}
	var req CreateGroupBuyRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}
	room, err := h.groupBuyService.CreateRoom(c.Request.Context(), subject.UserID, req.PlanID, req.TargetCount)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, room)
}

func (h *GroupBuyHandler) Join(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not found in context")
		return
	}
	roomID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || roomID <= 0 {
		response.BadRequest(c, "Invalid room ID")
		return
	}
	room, err := h.groupBuyService.JoinRoom(c.Request.Context(), subject.UserID, roomID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, room)
}
