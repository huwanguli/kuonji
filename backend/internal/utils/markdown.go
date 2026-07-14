package utils

import (
	"bytes"
	"fmt"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

var md = goldmark.New(
	goldmark.WithExtensions(
		extension.Table,
	),
	goldmark.WithRendererOptions(
		html.WithUnsafe(),
	),
)

func RenderMarkdown(content string) (string, error) {
	var buf bytes.Buffer
	if err := md.Convert([]byte(content), &buf); err != nil {
		return "", fmt.Errorf("render markdown: %w", err)
	}
	return buf.String(), nil
}
