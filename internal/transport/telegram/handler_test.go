package telegram

import (
	"testing"

	"github.com/antonhancharyk/crypto-knight-tg-bot/internal/config"
	"github.com/stretchr/testify/require"
)

func TestHandler_isAdmin(t *testing.T) {
	cfg := &config.Config{UserIDs: []int64{123, 456}}
	h := &Handler{cfg: cfg, states: make(map[int64]*userFlowState)}

	require.True(t, h.isAdmin(123))
	require.True(t, h.isAdmin(456))
	require.False(t, h.isAdmin(999))
	require.False(t, h.isAdmin(0))
}

func TestHandler_getState(t *testing.T) {
	h := &Handler{states: make(map[int64]*userFlowState)}

	st1 := h.getState(100)
	st2 := h.getState(100)
	require.Same(t, st1, st2, "same user should get same state")

	st3 := h.getState(200)
	require.NotSame(t, st1, st3)
}
