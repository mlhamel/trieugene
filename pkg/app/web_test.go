package app

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mlhamel/trieugene/pkg/config"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	cfg := config.NewConfig()

	req, _ := http.NewRequest("GET", "/ping", nil)
	rr := httptest.NewRecorder()

	web := NewWeb(cfg)

	runner, _ := web.(*Web)

	handler := http.HandlerFunc(runner.Ping)
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	require.Equal(t, "pong", rr.Body.String())
}
