package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupTestContext() (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	return c, w
}

func TestSuccess(t *testing.T) {
	c, w := setupTestContext()
	Success(c, gin.H{"key": "value"})
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}

func TestBadRequest(t *testing.T) {
	c, w := setupTestContext()
	BadRequest(c, "invalid param")
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestUnauthorized(t *testing.T) {
	c, w := setupTestContext()
	Unauthorized(c, "no auth")
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", w.Code)
	}
}

func TestForbidden(t *testing.T) {
	c, w := setupTestContext()
	Forbidden(c, "no permission")
	if w.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d", w.Code)
	}
}

func TestNotFound(t *testing.T) {
	c, w := setupTestContext()
	NotFound(c, "not found")
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", w.Code)
	}
}

func TestInternalError(t *testing.T) {
	c, w := setupTestContext()
	InternalError(c, "boom")
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", w.Code)
	}
}

func TestPageSuccess(t *testing.T) {
	c, w := setupTestContext()
	PageSuccess(c, []string{"a", "b"}, 100, 1, 10)
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}

func TestSuccessWithMessage(t *testing.T) {
	c, w := setupTestContext()
	SuccessWithMessage(c, "ok", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}
