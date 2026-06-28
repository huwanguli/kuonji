package handler_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"zblog-backend/internal/testutil"
)

func createArticle(t *testing.T, r http.Handler, token string, title string, status int) uint {
	body := fmt.Sprintf(`{"title":"%s","content_md":"# %s","status":%d}`, title, title, status)
	req := httptest.NewRequest("POST", "/api/admin/articles", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var resp struct {
		Data struct {
			ID uint `json:"id"`
		} `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	return resp.Data.ID
}

func TestCreateArticle(t *testing.T) {
	r := testutil.SetupTestRouter(t)
	token := loginAndGetToken(t, r)

	body := `{"title":"Hello World","content_md":"# Hello\n\n**Bold**","status":1}`
	req := httptest.NewRequest("POST", "/api/admin/articles", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	m := parseJSON(t, w.Body.Bytes())
	assert.Equal(t, float64(200), m["code"])
	data := m["data"].(map[string]interface{})
	assert.NotEmpty(t, data["slug"])
	assert.NotEmpty(t, data["content_html"])
	assert.Contains(t, data["content_html"], "<strong>")
}

func TestCreateArticleAsDraft(t *testing.T) {
	r := testutil.SetupTestRouter(t)
	token := loginAndGetToken(t, r)

	body := `{"title":"Draft Post","content_md":"draft content","status":0}`
	req := httptest.NewRequest("POST", "/api/admin/articles", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	m := parseJSON(t, w.Body.Bytes())
	assert.Equal(t, float64(200), m["code"])
	data := m["data"].(map[string]interface{})
	assert.Equal(t, float64(0), data["status"])
}

func TestGetPublicArticles(t *testing.T) {
	r := testutil.SetupTestRouter(t)
	token := loginAndGetToken(t, r)

	createArticle(t, r, token, "Published Post", 1)
	createArticle(t, r, token, "Draft Post", 0)

	req := httptest.NewRequest("GET", "/api/articles", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	m := parseJSON(t, w.Body.Bytes())
	assert.Equal(t, float64(200), m["code"])
	data := m["data"].(map[string]interface{})
	assert.GreaterOrEqual(t, data["total"].(float64), float64(1))
}

func TestGetArticleBySlug(t *testing.T) {
	r := testutil.SetupTestRouter(t)
	token := loginAndGetToken(t, r)

	body := `{"title":"My Great Article","content_md":"# Hello","status":1}`
	req := httptest.NewRequest("POST", "/api/admin/articles", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	m := parseJSON(t, w.Body.Bytes())
	slug := m["data"].(map[string]interface{})["slug"].(string)

	req = httptest.NewRequest("GET", "/api/articles/"+slug, nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	m = parseJSON(t, w.Body.Bytes())
	assert.Equal(t, float64(200), m["code"])
	data := m["data"].(map[string]interface{})
	articleData := data["article"].(map[string]interface{})
	assert.Equal(t, "My Great Article", articleData["title"])
	assert.Greater(t, articleData["view_count"].(float64), float64(0))
}

func TestGetArticleBySlugNotFound(t *testing.T) {
	r := testutil.SetupTestRouter(t)

	req := httptest.NewRequest("GET", "/api/articles/nonexistent", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	m := parseJSON(t, w.Body.Bytes())
	assert.Equal(t, float64(404), m["code"])
}

func TestUpdateArticle(t *testing.T) {
	r := testutil.SetupTestRouter(t)
	token := loginAndGetToken(t, r)

	body := `{"title":"Original","content_md":"# Original","status":1}`
	req := httptest.NewRequest("POST", "/api/admin/articles", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	m := parseJSON(t, w.Body.Bytes())
	id := int(m["data"].(map[string]interface{})["id"].(float64))

	updateBody := `{"title":"Updated","content_md":"# Updated","status":1}`
	req = httptest.NewRequest("PUT", fmt.Sprintf("/api/admin/articles/%d", id), strings.NewReader(updateBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	m = parseJSON(t, w.Body.Bytes())
	assert.Equal(t, float64(200), m["code"])
	data := m["data"].(map[string]interface{})
	assert.Equal(t, "Updated", data["title"])
}

func TestDeleteArticle(t *testing.T) {
	r := testutil.SetupTestRouter(t)
	token := loginAndGetToken(t, r)

	body := `{"title":"ToDelete","content_md":"...","status":1}`
	req := httptest.NewRequest("POST", "/api/admin/articles", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	m := parseJSON(t, w.Body.Bytes())
	id := int(m["data"].(map[string]interface{})["id"].(float64))

	req = httptest.NewRequest("DELETE", fmt.Sprintf("/api/admin/articles/%d", id), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	m = parseJSON(t, w.Body.Bytes())
	assert.Equal(t, float64(200), m["code"])
}

func TestCreateArticleWithoutAuth(t *testing.T) {
	r := testutil.SetupTestRouter(t)

	body := `{"title":"NoAuth","content_md":"...","status":1}`
	req := httptest.NewRequest("POST", "/api/admin/articles", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	m := parseJSON(t, w.Body.Bytes())
	assert.Equal(t, float64(401), m["code"])
}

func TestArticlePagination(t *testing.T) {
	r := testutil.SetupTestRouter(t)
	token := loginAndGetToken(t, r)

	for i := 0; i < 15; i++ {
		createArticle(t, r, token, fmt.Sprintf("Post %d", i), 1)
	}

	req := httptest.NewRequest("GET", "/api/articles?page=1&page_size=5", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	m := parseJSON(t, w.Body.Bytes())
	assert.Equal(t, float64(200), m["code"])
	data := m["data"].(map[string]interface{})
	list := data["list"].([]interface{})
	assert.Len(t, list, 5)
	assert.GreaterOrEqual(t, data["total"].(float64), float64(15))
}

func TestArticleSeries(t *testing.T) {
	r := testutil.SetupTestRouter(t)
	token := loginAndGetToken(t, r)

	body := fmt.Sprintf(`{"title":"Part 1","content_md":"one","status":1,"series":"my-series","series_order":1}`)
	req := httptest.NewRequest("POST", "/api/admin/articles", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	body2 := fmt.Sprintf(`{"title":"Part 2","content_md":"two","status":1,"series":"my-series","series_order":2}`)
	req2 := httptest.NewRequest("POST", "/api/admin/articles", strings.NewReader(body2))
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("Authorization", "Bearer "+token)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)

	// Get series list
	req3 := httptest.NewRequest("GET", "/api/series", nil)
	w3 := httptest.NewRecorder()
	r.ServeHTTP(w3, req3)
	m3 := parseJSON(t, w3.Body.Bytes())
	assert.Equal(t, float64(200), m3["code"])
	seriesList := m3["data"].([]interface{})
	assert.Len(t, seriesList, 1)

	// Get Part 2 and verify prev link
	req4 := httptest.NewRequest("GET", "/api/articles/part-2", nil)
	w4 := httptest.NewRecorder()
	r.ServeHTTP(w4, req4)
	m4 := parseJSON(t, w4.Body.Bytes())
	assert.Equal(t, float64(200), m4["code"])
	data := m4["data"].(map[string]interface{})
	prev := data["prev_in_series"].(map[string]interface{})
	assert.Equal(t, "part-1", prev["slug"])
	assert.Equal(t, "Part 1", prev["title"])
	assert.Nil(t, data["next_in_series"])
}
