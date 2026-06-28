package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateSlugEnglish(t *testing.T) {
	slug := GenerateSlug("Hello World Blog Post")
	assert.Equal(t, "hello-world-blog-post", slug)
}

func TestGenerateSlugWithSpecialChars(t *testing.T) {
	slug := GenerateSlug("Hello & World!")
	assert.Equal(t, "hello-world", slug)
}

func TestGenerateSlugChinese(t *testing.T) {
	slug := GenerateSlug("你好世界")
	assert.NotEmpty(t, slug)
	assert.Contains(t, slug, "post-")
}

func TestGenerateSlugEmpty(t *testing.T) {
	slug := GenerateSlug("")
	assert.NotEmpty(t, slug)
}

func TestGenerateSlugLong(t *testing.T) {
	longTitle := ""
	for i := 0; i < 250; i++ {
		longTitle += "a"
	}
	slug := GenerateSlug(longTitle)
	assert.LessOrEqual(t, len(slug), 200)
}

func TestRenderMarkdown(t *testing.T) {
	html, err := RenderMarkdown("# Hello\n\nThis is **bold**")
	require.NoError(t, err)
	assert.Contains(t, html, "<h1>Hello</h1>")
	assert.Contains(t, html, "<strong>bold</strong>")
}

func TestRenderMarkdownCode(t *testing.T) {
	html, err := RenderMarkdown("```go\nfmt.Println()\n```")
	require.NoError(t, err)
	assert.Contains(t, html, "<code")
}

func TestGenerateAndParseToken(t *testing.T) {
	secret := "test-secret"
	token, err := GenerateToken(1, "admin", secret, 24*time.Hour)
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	claims, err := ParseToken(token, secret)
	require.NoError(t, err)
	assert.Equal(t, uint(1), claims.UserID)
	assert.Equal(t, "admin", claims.Username)
}

func TestParseTokenInvalid(t *testing.T) {
	_, err := ParseToken("invalid.token.here", "secret")
	assert.Error(t, err)
}

func TestParseTokenWrongSecret(t *testing.T) {
	secret := "test-secret"
	token, err := GenerateToken(1, "admin", secret, 24*time.Hour)
	require.NoError(t, err)

	_, err = ParseToken(token, "wrong-secret")
	assert.Error(t, err)
}
