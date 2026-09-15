package redis

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/flynn/flynn/pkg/status"
	"github.com/inconshreveable/log15"
)

func TestHandlerStatusWhenStopped(t *testing.T) {
	h := NewHandler()
	h.Process = &Process{}
	h.Logger = log15.New()

	req := httptest.NewRequest(http.MethodGet, "/status", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatalf("status=%d", rec.Code)
	}
	var body Status
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.Process == nil || body.Process.Running {
		t.Fatalf("stopped process must report running=false: %+v", body.Process)
	}
	if h.healthStatus() != status.Unhealthy {
		t.Fatal("stopped redis must be unhealthy")
	}
}
