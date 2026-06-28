package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"zblog-backend/internal/testutil"
)

func parseJSON(t *testing.T, body []byte) map[string]interface{} {
	t.Helper()
	var m map[string]interface{}
	err := json.Unmarshal(body, &m)
	require.NoError(t, err)
	return m
}

func loginAndGetToken(t *testing.T, r http.Handler) string {
	body := `{"username":"admin","password":"admin123"}`
	req := httptest.NewRequest("POST", "/api/admin/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var resp struct {
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	return resp.Data.Token
}

func TestLogin(t *testing.T) {
	r := testutil.SetupTestRouter(t)

	body := `{"username":"admin","password":"admin123"}`
	req := httptest.NewRequest("POST", "/api/admin/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.Code)
	assert.NotEmpty(t, resp.Data.Token)
}

func TestLoginInvalidPassword(t *testing.T) {
	r := testutil.SetupTestRouter(t)

	body := `{"username":"admin","password":"wrong"}`
	req := httptest.NewRequest("POST", "/api/admin/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	m := parseJSON(t, w.Body.Bytes())
	assert.Equal(t, float64(400), m["code"])
}

func TestLoginMissingFields(t *testing.T) {
	r := testutil.SetupTestRouter(t)

	body := `{"username":"admin"}`
	req := httptest.NewRequest("POST", "/api/admin/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	m := parseJSON(t, w.Body.Bytes())
	assert.Equal(t, float64(400), m["code"])
}

func TestProfile(t *testing.T) {
	r := testutil.SetupTestRouter(t)
	token := loginAndGetToken(t, r)

	req := httptest.NewRequest("GET", "/api/admin/profile", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	m := parseJSON(t, w.Body.Bytes())
	assert.Equal(t, float64(200), m["code"])
	data := m["data"].(map[string]interface{})
	assert.Equal(t, "admin", data["username"])
}

func TestProfileWithoutAuth(t *testing.T) {
	r := testutil.SetupTestRouter(t)

	req := httptest.NewRequest("GET", "/api/admin/profile", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	m := parseJSON(t, w.Body.Bytes())
	assert.Equal(t, float64(401), m["code"])
}
