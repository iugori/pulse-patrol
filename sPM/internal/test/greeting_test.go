package test

import (
	"net/http"
	"net/http/httptest"
	"sPM/internal/greeting"
	"strings"
	"testing"
)

func TestGreetEndpoint(t *testing.T) {
	req, _ := http.NewRequest("GET", "/hello?name=sPM", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(greeting.GreetHandler)

	handler.ServeHTTP(rr, req)

	expected := `{"message":"Hello sPM!"}`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("Unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
