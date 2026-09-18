package redis

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/inconshreveable/log15"
	discoverd "github.com/randy-girard/flynn/discoverd/client"
	"github.com/randy-girard/flynn/pkg/status"
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

type hbStub struct {
	closed bool
}

func (h *hbStub) SetMeta(map[string]string) error { return nil }
func (h *hbStub) Close() error                    { h.closed = true; return nil }
func (h *hbStub) Addr() string                    { return "127.0.0.1:1" }
func (h *hbStub) SetClient(*discoverd.Client)     {}

func TestProcessRestoreMissingDataDir(t *testing.T) {
	p := NewProcess()
	p.DataDir = filepath.Join(t.TempDir(), "missing")
	if err := p.Restore(strings.NewReader("RDB")); err == nil {
		t.Fatal("restore into a missing DataDir must fail")
	}
}

func TestHandlerStopClosesHeartbeater(t *testing.T) {
	h := NewHandler()
	h.Process = &Process{}
	h.Logger = log15.New()
	hb := &hbStub{}
	h.Heartbeater = hb

	req := httptest.NewRequest(http.MethodPost, "/stop", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != 200 || !hb.closed {
		t.Fatalf("stop status=%d closed=%v", rec.Code, hb.closed)
	}
}
