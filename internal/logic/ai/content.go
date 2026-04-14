package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	v1 "github.com/nuxtblog/nuxtblog/api/ai/v1"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

func (s *sAI) FromURL(ctx context.Context, url, style string) (*v1.AIFromURLRes, error) {
	cfg, err := getActiveConfig(ctx)
	if err != nil {
		return nil, err
	}
	// Fetch the page content
	pageText, fetchErr := fetchPageText(ctx, url)
	if fetchErr != nil {
		return nil, gerror.NewCode(gcode.CodeInvalidOperation,
			fmt.Sprintf("fetch URL failed: %v", fetchErr))
	}
	// Truncate to 8000 chars to fit in context
	if len(pageText) > 8000 {
		pageText = pageText[:8000] + "..."
	}
	stylePart := ""
	switch style {
	case "formal":
		stylePart = "语气正式专业。"
	case "casual":
		stylePart = "语气轻松口语。"
	case "concise":
		stylePart = "简洁为主。"
	}
	system := fmt.Sprintf(`你是一位博客写手。请根据以下网页内容改写成一篇原创博客文章。%s
请以如下 JSON 格式输出（不要添加 markdown 代码块）：
{"title":"...","content":"...","excerpt":"..."}
其中 content 为 HTML 格式正文。`, stylePart)

	raw, err := generateText(ctx, *cfg, system, pageText)
	if err != nil {
		return nil, err
	}
	// Parse JSON response
	raw = strings.TrimSpace(raw)
	// Strip markdown code fences if model added them
	raw = strings.TrimPrefix(raw, "```json")
	raw = strings.TrimPrefix(raw, "```")
	raw = strings.TrimSuffix(raw, "```")
	raw = strings.TrimSpace(raw)

	var res v1.AIFromURLRes
	if err := json.Unmarshal([]byte(raw), &res); err != nil {
		// Fallback: return raw as content
		res.Title = "从 URL 生成的文章"
		res.Content = raw
	}
	return &res, nil
}

// fetchPageText fetches a URL and extracts plain text from the HTML.
func fetchPageText(ctx context.Context, url string) (string, error) {
	client := &http.Client{Timeout: 15 * time.Second}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; NuxtBlog/2.0)")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 512*1024)) // max 512 KB
	if err != nil {
		return "", err
	}
	// Very basic HTML → text: strip tags
	text := string(body)
	// Remove script/style blocks
	for _, tag := range []string{"script", "style", "head"} {
		for {
			start := strings.Index(strings.ToLower(text), "<"+tag)
			end := strings.Index(strings.ToLower(text), "</"+tag+">")
			if start < 0 || end < 0 || end < start {
				break
			}
			text = text[:start] + text[end+len("</"+tag+">"):]
		}
	}
	// Strip remaining HTML tags
	var b strings.Builder
	inTag := false
	for _, c := range text {
		if c == '<' {
			inTag = true
		} else if c == '>' {
			inTag = false
			b.WriteRune(' ')
		} else if !inTag {
			b.WriteRune(c)
		}
	}
	// Collapse whitespace
	result := b.String()
	lines := strings.Split(result, "\n")
	var clean []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			clean = append(clean, line)
		}
	}
	return strings.Join(clean, "\n"), nil
}
