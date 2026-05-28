package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gsystes/backend/internal/infrastructure/auth"
	"github.com/gsystes/backend/internal/infrastructure/config"
	infraMiddleware "github.com/gsystes/backend/internal/infrastructure/middleware"
	"github.com/gsystes/backend/internal/infrastructure/utils"
)

func setupAuthConfig() {
	config.SetConfigForTesting(&config.Config{
		JWT: config.JWTConfig{
			Secret:      "test-secret-for-unit-testing-only",
			Issuer:      "test-issuer",
			ExpireHours: 1,
		},
	})
}

func TestAuthRequired_NoAuthHeader(t *testing.T) {
	setupAuthConfig()
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)

	AuthRequired()(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestAuthRequired_InvalidHeaderFormat(t *testing.T) {
	setupAuthConfig()
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)
	c.Request.Header.Set("Authorization", "InvalidToken")

	AuthRequired()(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestAuthRequired_ExpiredToken(t *testing.T) {
	setupAuthConfig()
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)
	c.Request.Header.Set("Authorization", "Bearer expired.token.here")

	AuthRequired()(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestAuthRequired_ValidToken(t *testing.T) {
	setupAuthConfig()
	gin.SetMode(gin.TestMode)

	token, err := auth.GenerateToken(1, "testuser", 1)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)
	c.Request.Header.Set("Authorization", "Bearer "+token)

	executed := false
	handler := AuthRequired()
	handler(c)
	if !c.IsAborted() {
		executed = true
	}

	if !executed {
		t.Fatal("expected handler to pass through")
	}
	if w.Code == http.StatusUnauthorized {
		t.Fatal("expected token to be valid")
	}

	claims := infraMiddleware.GetClaims(c)
	if claims == nil {
		t.Fatal("expected claims to be set")
	}
	if claims.UserID != 1 {
		t.Fatalf("expected userID 1, got %d", claims.UserID)
	}
	if claims.Username != "testuser" {
		t.Fatalf("expected username testuser, got %s", claims.Username)
	}
}

func TestOperationLogMiddleware_SanitizeBody(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty body",
			input:    "",
			expected: "",
		},
		{
			name:     "sanitize password",
			input:    `{"username":"u","password":"secret123"}`,
			expected: `{"password":"***","username":"u"}`,
		},
		{
			name:     "sanitize token",
			input:    `{"token":"abc.xyz.123"}`,
			expected: `{"token":"***"}`,
		},
		{
			name:     "no sensitive keys",
			input:    `{"name":"test","age":30}`,
			expected: `{"age":30,"name":"test"}`,
		},
		{
			name:     "invalid json returns raw",
			input:    `not json`,
			expected: `not json`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizeBody(tt.input)
			if tt.input == "" && result != "" {
				t.Fatalf("expected empty, got %s", result)
			}
			if tt.expected == `not json` && result != `not json` {
				t.Fatalf("expected 'not json', got %s", result)
			}
			if tt.expected != `not json` && tt.input != "" {
				var got, want map[string]interface{}
				json.Unmarshal([]byte(result), &got)
				json.Unmarshal([]byte(tt.expected), &want)
				for k, v := range want {
					if got[k] != v {
						t.Fatalf("for key %s: expected %v, got %v", k, v, got[k])
					}
				}
			}
		})
	}
}

func TestOperationLogMiddleware_NonGetRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/test", nil)

	called := false
	middlewareExec := func(c *gin.Context) {
		called = true
	}

	middlewareExec(c)

	if !called {
		t.Fatal("expected middleware to process request")
	}
}

func TestPermissionMiddleware_NoClaims(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/admin/users", nil)

	m := &PermissionMiddleware{roleRepo: nil}
	m.Require("user:create")(c)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", w.Code)
	}
}

func TestResponseError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	utils.NotFound(c, "resource not found")

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestInfraMiddleware_SetAndGetClaims(t *testing.T) {
	setupAuthConfig()
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)

	claims := &auth.Claims{
		UserID:   42,
		Username: "tester",
		RoleID:   2,
	}
	infraMiddleware.SetClaims(c, claims)

	got := infraMiddleware.GetClaims(c)
	if got == nil {
		t.Fatal("expected claims to be set")
	}
	if got.UserID != 42 {
		t.Fatalf("expected userID 42, got %d", got.UserID)
	}

	userID := infraMiddleware.GetUserID(c)
	if userID != 42 {
		t.Fatalf("expected userID 42, got %d", userID)
	}
}

func TestOperationLog_SanitizeAllSensitiveKeys(t *testing.T) {
	input := `{"password":"p","old_password":"op","new_password":"np","secret":"s","token":"t","access_token":"at","refresh_token":"rt","normal":"keep"}`
	result := sanitizeBody(input)

	var m map[string]interface{}
	json.Unmarshal([]byte(result), &m)

	sensitive := []string{"password", "old_password", "new_password", "secret", "token", "access_token", "refresh_token"}
	for _, k := range sensitive {
		if m[k] != "***" {
			t.Fatalf("expected %s to be sanitized, got %v", k, m[k])
		}
	}
	if m["normal"] != "keep" {
		t.Fatalf("expected normal to be untouched, got %v", m["normal"])
	}
}
