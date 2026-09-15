package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestDeleteClusterRejectsInvalidID(t *testing.T) {
	h := NewHandler()
	for _, id := range []string{"../etc/passwd", "id/with/slash", "id\nnewline", ""} {
		req := httptest.NewRequest(http.MethodDelete, "/clusters?"+url.Values{"id": {id}}.Encode(), nil)
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		if rec.Code < 400 {
			t.Fatalf("id %q must be rejected, status=%d body=%s", id, rec.Code, rec.Body.Bytes())
		}
	}
}

func TestGetPing(t *testing.T) {
	h := NewHandler()
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatalf("ping status=%d", rec.Code)
	}
}
