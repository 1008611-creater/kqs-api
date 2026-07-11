package handler

import (
	"time"

	"github.com/Wei-Shaw/sub2api/internal/handler/admin"
	"github.com/Wei-Shaw/sub2api/internal/handler/dto"
	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

// ChannelMonitorUserHandler 渠道监控用户只读 handler。
type ChannelMonitorUserHandler struct {
	monitorService      *service.ChannelMonitorService
	publicStatusService *service.PublicChannelStatusService
	settingService      *service.SettingService
}

// NewChannelMonitorUserHandler 创建 handler。
// settingService 用于每次请求前读取功能开关；关闭时 List/GetStatus 直接返回空/404。
func NewChannelMonitorUserHandler(
	monitorService *service.ChannelMonitorService,
	publicStatusService *service.PublicChannelStatusService,
	settingService *service.SettingService,
) *ChannelMonitorUserHandler {
	return &ChannelMonitorUserHandler{
		monitorService:      monitorService,
		publicStatusService: publicStatusService,
		settingService:      settingService,
	}
}

// featureEnabled 返回当前渠道监控功能是否开启。
// settingService 为 nil（测试场景）视为启用。
func (h *ChannelMonitorUserHandler) featureEnabled(c *gin.Context) bool {
	if h.settingService == nil {
		return true
	}
	return h.settingService.GetChannelMonitorRuntime(c.Request.Context()).Enabled
}

// --- Response ---

type channelMonitorUserListItem struct {
	ID                   int64                                `json:"id"`
	Name                 string                               `json:"name"`
	Provider             string                               `json:"provider"`
	GroupName            string                               `json:"group_name"`
	PrimaryModel         string                               `json:"primary_model"`
	PrimaryStatus        string                               `json:"primary_status"`
	PrimaryLatencyMs     *int                                 `json:"primary_latency_ms"`
	PrimaryPingLatencyMs *int                                 `json:"primary_ping_latency_ms"`
	Availability7d       float64                              `json:"availability_7d"`
	ExtraModels          []dto.ChannelMonitorExtraModelStatus `json:"extra_models"`
	Timeline             []channelMonitorUserTimelinePoint    `json:"timeline"`
}

// channelMonitorUserTimelinePoint 主模型最近一次检测的 timeline 点。
// 仅用于用户视图 list 响应，admin 视图不使用。
type channelMonitorUserTimelinePoint struct {
	Status        string `json:"status"`
	LatencyMs     *int   `json:"latency_ms"`
	PingLatencyMs *int   `json:"ping_latency_ms"`
	CheckedAt     string `json:"checked_at"`
}

type channelMonitorUserDetailResponse struct {
	ID        int64                         `json:"id"`
	Name      string                        `json:"name"`
	Provider  string                        `json:"provider"`
	GroupName string                        `json:"group_name"`
	Models    []channelMonitorUserModelStat `json:"models"`
}

type channelMonitorUserModelStat struct {
	Model           string  `json:"model"`
	LatestStatus    string  `json:"latest_status"`
	LatestLatencyMs *int    `json:"latest_latency_ms"`
	Availability7d  float64 `json:"availability_7d"`
	Availability15d float64 `json:"availability_15d"`
	Availability30d float64 `json:"availability_30d"`
	AvgLatency7dMs  *int    `json:"avg_latency_7d_ms"`
}

type publicChannelStatusResponse struct {
	GroupName        string                             `json:"group_name"`
	ModelID          string                             `json:"model_id"`
	Status           string                             `json:"status"`
	ActiveChannels   int                                `json:"active_channels"`
	CheckedChannels  int                                `json:"checked_channels"`
	SuccessRate1h    float64                            `json:"success_rate_1h"`
	SuccessRate24h   float64                            `json:"success_rate_24h"`
	AverageLatencyMs *int                               `json:"average_latency_ms"`
	P95LatencyMs     *int                               `json:"p95_latency_ms"`
	LastCheckedAt    *string                            `json:"last_checked_at"`
	GeneratedAt      string                             `json:"generated_at"`
	Timeline         []publicChannelStatusTimelinePoint `json:"timeline"`
}

type publicChannelStatusTimelinePoint struct {
	Status    string `json:"status"`
	CheckedAt string `json:"checked_at"`
}

func userMonitorViewToItem(v *service.UserMonitorView) channelMonitorUserListItem {
	extras := make([]dto.ChannelMonitorExtraModelStatus, 0, len(v.ExtraModels))
	for _, e := range v.ExtraModels {
		extras = append(extras, dto.ChannelMonitorExtraModelStatus{
			Model:     e.Model,
			Status:    e.Status,
			LatencyMs: e.LatencyMs,
		})
	}
	timeline := make([]channelMonitorUserTimelinePoint, 0, len(v.Timeline))
	for _, p := range v.Timeline {
		timeline = append(timeline, channelMonitorUserTimelinePoint{
			Status:        p.Status,
			LatencyMs:     p.LatencyMs,
			PingLatencyMs: p.PingLatencyMs,
			CheckedAt:     p.CheckedAt.UTC().Format(time.RFC3339),
		})
	}
	return channelMonitorUserListItem{
		ID:                   v.ID,
		Name:                 v.Name,
		Provider:             v.Provider,
		GroupName:            v.GroupName,
		PrimaryModel:         v.PrimaryModel,
		PrimaryStatus:        v.PrimaryStatus,
		PrimaryLatencyMs:     v.PrimaryLatencyMs,
		PrimaryPingLatencyMs: v.PrimaryPingLatencyMs,
		Availability7d:       v.Availability7d,
		ExtraModels:          extras,
		Timeline:             timeline,
	}
}

func userMonitorDetailToResponse(d *service.UserMonitorDetail) *channelMonitorUserDetailResponse {
	models := make([]channelMonitorUserModelStat, 0, len(d.Models))
	for _, m := range d.Models {
		models = append(models, channelMonitorUserModelStat{
			Model:           m.Model,
			LatestStatus:    m.LatestStatus,
			LatestLatencyMs: m.LatestLatencyMs,
			Availability7d:  m.Availability7d,
			Availability15d: m.Availability15d,
			Availability30d: m.Availability30d,
			AvgLatency7dMs:  m.AvgLatency7dMs,
		})
	}
	return &channelMonitorUserDetailResponse{
		ID:        d.ID,
		Name:      d.Name,
		Provider:  d.Provider,
		GroupName: d.GroupName,
		Models:    models,
	}
}

func publicChannelStatusToResponse(v *service.PublicChannelStatusSnapshot) publicChannelStatusResponse {
	var lastCheckedAt *string
	if v.LastCheckedAt != nil {
		formatted := v.LastCheckedAt.UTC().Format(time.RFC3339)
		lastCheckedAt = &formatted
	}
	timeline := make([]publicChannelStatusTimelinePoint, 0, len(v.Timeline))
	for _, p := range v.Timeline {
		timeline = append(timeline, publicChannelStatusTimelinePoint{
			Status:    p.Status,
			CheckedAt: p.CheckedAt.UTC().Format(time.RFC3339),
		})
	}
	return publicChannelStatusResponse{
		GroupName:        v.GroupName,
		ModelID:          v.ModelID,
		Status:           v.Status,
		ActiveChannels:   v.ActiveChannels,
		CheckedChannels:  v.CheckedChannels,
		SuccessRate1h:    v.SuccessRate1h,
		SuccessRate24h:   v.SuccessRate24h,
		AverageLatencyMs: v.AverageLatencyMs,
		P95LatencyMs:     v.P95LatencyMs,
		LastCheckedAt:    lastCheckedAt,
		GeneratedAt:      v.GeneratedAt.UTC().Format(time.RFC3339),
		Timeline:         timeline,
	}
}

// --- Handlers ---

// List GET /api/v1/channel-monitors
func (h *ChannelMonitorUserHandler) List(c *gin.Context) {
	if !h.featureEnabled(c) {
		response.Success(c, gin.H{"items": []channelMonitorUserListItem{}})
		return
	}
	views, err := h.monitorService.ListUserView(c.Request.Context())
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	items := make([]channelMonitorUserListItem, 0, len(views))
	for _, v := range views {
		items = append(items, userMonitorViewToItem(v))
	}
	response.Success(c, gin.H{"items": items})
}

// PublicStatus GET /api/v1/channel-monitors/public-status
// Optional query: ?profile=pro4. Empty profile keeps the default gptplus2.5倍 status.
func (h *ChannelMonitorUserHandler) PublicStatus(c *gin.Context) {
	if h.publicStatusService == nil {
		groupName := service.PublicChannelStatusDisplayName
		modelID := service.PublicChannelStatusModelID
		switch c.Query("profile") {
		case "pro4":
			groupName = service.PublicChannelStatusPro4DisplayName
		case "gpt56":
			groupName = service.PublicChannelStatusGPT56DisplayName
			modelID = service.PublicChannelStatusGPT56ModelID
		}
		response.Success(c, publicChannelStatusResponse{
			GroupName:   groupName,
			ModelID:     modelID,
			Status:      "unknown",
			GeneratedAt: time.Now().UTC().Format(time.RFC3339),
			Timeline:    []publicChannelStatusTimelinePoint{},
		})
		return
	}
	snapshot, err := h.publicStatusService.GetSnapshotByProfile(c.Request.Context(), c.Query("profile"))
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, publicChannelStatusToResponse(snapshot))
}

// GetStatus GET /api/v1/channel-monitors/:id/status
func (h *ChannelMonitorUserHandler) GetStatus(c *gin.Context) {
	if !h.featureEnabled(c) {
		response.ErrorFrom(c, service.ErrChannelMonitorNotFound)
		return
	}
	// 复用 admin.ParseChannelMonitorID 保持错误码与日志一致。
	id, ok := admin.ParseChannelMonitorID(c)
	if !ok {
		return
	}
	detail, err := h.monitorService.GetUserDetail(c.Request.Context(), id)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, userMonitorDetailToResponse(detail))
}
