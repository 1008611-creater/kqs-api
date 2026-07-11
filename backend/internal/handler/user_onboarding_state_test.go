//go:build unit

package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	middleware2 "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type userOnboardingAttributeServiceStub struct {
	def    *service.UserAttributeDefinition
	values map[int64]string
}

func (s *userOnboardingAttributeServiceStub) GetDefinitionByKey(_ context.Context, key string) (*service.UserAttributeDefinition, error) {
	if s.def == nil || s.def.Key != key {
		return nil, service.ErrAttributeDefinitionNotFound
	}
	return s.def, nil
}

func (s *userOnboardingAttributeServiceStub) CreateDefinition(_ context.Context, input service.CreateAttributeDefinitionInput) (*service.UserAttributeDefinition, error) {
	if s.def != nil && s.def.Key == input.Key {
		return nil, service.ErrAttributeKeyExists
	}
	s.def = &service.UserAttributeDefinition{
		ID:          99,
		Key:         input.Key,
		Name:        input.Name,
		Description: input.Description,
		Type:        input.Type,
		Enabled:     input.Enabled,
	}
	return s.def, nil
}

func (s *userOnboardingAttributeServiceStub) GetUserAttributes(_ context.Context, userID int64) ([]service.UserAttributeValue, error) {
	if s.values == nil || s.def == nil {
		return nil, nil
	}
	value, ok := s.values[userID]
	if !ok {
		return nil, nil
	}
	return []service.UserAttributeValue{
		{UserID: userID, AttributeID: s.def.ID, Value: value},
	}, nil
}

func (s *userOnboardingAttributeServiceStub) UpdateUserAttributes(_ context.Context, userID int64, inputs []service.UpdateUserAttributeInput) error {
	if s.values == nil {
		s.values = map[int64]string{}
	}
	for _, input := range inputs {
		if s.def != nil && input.AttributeID == s.def.ID {
			s.values[userID] = input.Value
		}
	}
	return nil
}

func performUserOnboardingRequest(handlerFunc gin.HandlerFunc, method, path string, body []byte) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(method, path, bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set(string(middleware2.ContextKeyUser), middleware2.AuthSubject{UserID: 11})

	handlerFunc(c)
	return w
}

func decodeOnboardingResponse(t *testing.T, w *httptest.ResponseRecorder) userOnboardingStateResponse {
	t.Helper()
	var envelope struct {
		Code int                         `json:"code"`
		Data userOnboardingStateResponse `json:"data"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &envelope))
	require.Equal(t, 0, envelope.Code)
	return envelope.Data
}

func TestUserOnboardingStateDefaultsAndEvents(t *testing.T) {
	gin.SetMode(gin.TestMode)
	attrSvc := &userOnboardingAttributeServiceStub{}
	h := NewUserHandler(nil, nil, nil, nil, nil, attrSvc)

	w := performUserOnboardingRequest(h.GetOnboardingState, http.MethodGet, "/api/v1/user/onboarding-state", nil)
	require.Equal(t, http.StatusOK, w.Code)
	state := decodeOnboardingResponse(t, w)
	require.Equal(t, userOnboardingStepPurchase, state.CurrentStep)
	require.Empty(t, state.CompletedSteps)
	require.Equal(t, 0, state.ProgressPercent)

	w = performUserOnboardingRequest(
		h.UpdateOnboardingState,
		http.MethodPatch,
		"/api/v1/user/onboarding-state",
		[]byte(`{"event":"redeem_success"}`),
	)
	require.Equal(t, http.StatusOK, w.Code)
	state = decodeOnboardingResponse(t, w)
	require.Equal(t, userOnboardingStepKey, state.CurrentStep)
	require.Equal(t, []userOnboardingStep{userOnboardingStepPurchase, userOnboardingStepRedeem}, state.CompletedSteps)
	require.Equal(t, 50, state.ProgressPercent)
	require.NotNil(t, state.RedeemCompletedAt)

	w = performUserOnboardingRequest(
		h.UpdateOnboardingState,
		http.MethodPatch,
		"/api/v1/user/onboarding-state",
		[]byte(`{"event":"key_created"}`),
	)
	require.Equal(t, http.StatusOK, w.Code)
	state = decodeOnboardingResponse(t, w)
	require.Equal(t, userOnboardingStepConfig, state.CurrentStep)
	require.Contains(t, state.CompletedSteps, userOnboardingStepKey)
	require.Equal(t, 75, state.ProgressPercent)
	require.NotNil(t, state.KeyCreatedAt)

	w = performUserOnboardingRequest(
		h.UpdateOnboardingState,
		http.MethodPatch,
		"/api/v1/user/onboarding-state",
		[]byte(`{"event":"first_request"}`),
	)
	require.Equal(t, http.StatusOK, w.Code)
	state = decodeOnboardingResponse(t, w)
	require.Equal(t, userOnboardingStepConfig, state.CurrentStep)
	require.True(t, state.Completed)
	require.Equal(t, 100, state.ProgressPercent)
	require.NotNil(t, state.FirstRequestAt)
}

func TestUserOnboardingStateRejectsUnknownEvent(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewUserHandler(nil, nil, nil, nil, nil, &userOnboardingAttributeServiceStub{})

	w := performUserOnboardingRequest(
		h.UpdateOnboardingState,
		http.MethodPatch,
		"/api/v1/user/onboarding-state",
		[]byte(`{"event":"unknown"}`),
	)
	require.Equal(t, http.StatusBadRequest, w.Code)
}
