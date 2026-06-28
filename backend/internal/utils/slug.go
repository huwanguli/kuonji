package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"regexp"
	"strings"
	"unicode"
)

var nonAlphaNumRegex = regexp.MustCompile(`[^a-zA-Z0-9\-]+`)

func GenerateSlug(title string) string {
	slug := strings.ToLower(title)
	slug = strings.ReplaceAll(slug, " ", "-")

	slug = nonAlphaNumRegex.ReplaceAllString(slug, "")

	var result strings.Builder
	lastHyphen := false
	for _, r := range slug {
		if r == '-' {
			if !lastHyphen && result.Len() > 0 {
				result.WriteRune(r)
				lastHyphen = true
			}
			continue
		}
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result.WriteRune(r)
			lastHyphen = false
		}
	}

	slug = result.String()
	slug = strings.Trim(slug, "-")

	asciiLen := 0
	for _, r := range slug {
		if r <= unicode.MaxASCII {
			asciiLen++
		}
	}

	if len(slug) < 2 || asciiLen < 2 {
		slug = fmt.Sprintf("post-%s", randomString(6))
	}

	if len(slug) > 200 {
		slug = slug[:200]
	}

	return slug
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, n)
	for i := range result {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		result[i] = letters[num.Int64()]
	}
	return string(result)
}
