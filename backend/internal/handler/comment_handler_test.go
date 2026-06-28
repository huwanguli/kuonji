package handler_test

import (
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"zblog-backend/internal/testutil"
)

func TestCreateComment(t *testing.T) {
	r := testutil.SetupTestRouter(t)
	token := loginAndGetToken(t, r)
	articleID := createArticle(t, r, token, "Test Article", 1)

	body := fmt.Sprintf(`{"article_id":%d,"author":"Visitor","content":"Great post!"}`, articleID)
	req := httptest.NewRequest("POST", "/api/comments", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	m := parseJSON(t, w.Body.Bytes())
	assert.Equal(t, float64(200), m["code"])
	data := m["data"].(map[string]interface{})
	assert.Equal(t, "Visitor", data["author"])
	assert.Equal(t, float64(1), data["status"])
}

func TestCreateCommentMissingFields(t *testing.T) {
	r := testutil.SetupTestRouter(t)

	body := `{"author":"Visitor"}`
	req := httptest.NewRequest("POST", "/api/comments", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	m := parseJSON(t, w.Body.Bytes())
	assert.Equal(t, float64(400), m["code"])
}

func TestCreateCommentNonexistentArticle(t *testing.T) {
	r := testutil.SetupTestRouter(t)

	body := `{"article_id":9999,"author":"Visitor","content":"Hey"}`
	req := httptest.NewRequest("POST", "/api/comments", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	m := parseJSON(t, w.Body.Bytes())
	assert.Equal(t, float64(400), m["code"])
}

func TestGetComments(t *testing.T) {
	r := testutil.SetupTestRouter(t)
	token := loginAndGetToken(t, r)
	articleID := createArticle(t, r, token, "Test Article", 1)

	body := fmt.Sprintf(`{"article_id":%d,"author":"User1","content":"First comment"}`, articleID)
	req := httptest.NewRequest("POST", "/api/comments", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	req = httptest.NewRequest("GET", fmt.Sprintf("/api/comments?article_id=%d", articleID), nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	m := parseJSON(t, w.Body.Bytes())
	assert.Equal(t, float64(200), m["code"])
	data := m["data"].(map[string]interface{})
	assert.GreaterOrEqual(t, data["total"].(float64), float64(1))
}

func TestAdminUpdateCommentStatus(t *testing.T) {
	r := testutil.SetupTestRouter(t)
	token := loginAndGetToken(t, r)
	articleID := createArticle(t, r, token, "Test Article", 1)

	body := fmt.Sprintf(`{"article_id":%d,"author":"User","content":"Test comment"}`, articleID)
	req := httptest.NewRequest("POST", "/api/comments", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	m := parseJSON(t, w.Body.Bytes())
	commentID := int(m["data"].(map[string]interface{})["id"].(float64))

	body = `{"status":0}`
	req = httptest.NewRequest("PUT", fmt.Sprintf("/api/admin/comments/%d", commentID), strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	m = parseJSON(t, w.Body.Bytes())
	assert.Equal(t, float64(200), m["code"])
}

func TestAdminDeleteComment(t *testing.T) {
	r := testutil.SetupTestRouter(t)
	token := loginAndGetToken(t, r)
	articleID := createArticle(t, r, token, "Test Article", 1)

	body := fmt.Sprintf(`{"article_id":%d,"author":"User","content":"Delete me"}`, articleID)
	req := httptest.NewRequest("POST", "/api/comments", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	m := parseJSON(t, w.Body.Bytes())
	commentID := int(m["data"].(map[string]interface{})["id"].(float64))

	req = httptest.NewRequest("DELETE", fmt.Sprintf("/api/admin/comments/%d", commentID), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	m = parseJSON(t, w.Body.Bytes())
	assert.Equal(t, float64(200), m["code"])
}
