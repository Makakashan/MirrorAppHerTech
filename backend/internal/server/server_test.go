package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestChatEndpointReturnsReflectionAndUpdatesMessages(t *testing.T) {
	t.Setenv("MIRROR_DEMO_DATA", "false")
	app := New()
	body := bytes.NewBufferString(`{"message":"I felt anxious and could not focus"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/chat", body)
	rec := httptest.NewRecorder()

	app.Routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var res chatResponse
	if err := json.NewDecoder(rec.Body).Decode(&res); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if res.Message.Role != "ai" {
		t.Fatalf("expected ai message, got %q", res.Message.Role)
	}
	if res.Metrics.ChartData[6].Value >= 80 {
		t.Fatalf("expected last chart point to reflect anxious tone, got %d", res.Metrics.ChartData[6].Value)
	}

	msgReq := httptest.NewRequest(http.MethodGet, "/api/messages", nil)
	msgRec := httptest.NewRecorder()
	app.Routes().ServeHTTP(msgRec, msgReq)

	var messages []Message
	if err := json.NewDecoder(msgRec.Body).Decode(&messages); err != nil {
		t.Fatalf("decode messages: %v", err)
	}
	if len(messages) != 3 {
		t.Fatalf("expected welcome, user, and ai messages, got %d", len(messages))
	}
}

func TestChatEndpointRejectsEmptyMessage(t *testing.T) {
	t.Setenv("MIRROR_DEMO_DATA", "false")
	app := New()
	req := httptest.NewRequest(http.MethodPost, "/api/chat", bytes.NewBufferString(`{"message":"   "}`))
	rec := httptest.NewRecorder()

	app.Routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
}
