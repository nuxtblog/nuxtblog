package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	v1 "github.com/nuxtblog/nuxtblog/api/ai/v1"

	"github.com/gogf/gf/v2/errors/gerror"
)

type claudeRequest struct {
	Model     string       `json:"model"`
	MaxTokens int          `json:"max_tokens"`
	System    string       `json:"system,omitempty"`
	Messages  []oaiMessage `json:"messages"`
}

type claudeResponse struct {
	Content []struct {
		Text string `json:"text"`
	} `json:"content"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func claudeGenerate(ctx context.Context, cfg v1.AIConfig, system, user string, timeout time.Duration) (string, error) {
	base := cfg.BaseURL
	if base == "" {
		base = "https://api.anthropic.com"
	}
	base = strings.TrimRight(base, "/")
	endpoint := base + "/v1/messages"

	payload := claudeRequest{
		Model:     cfg.Model,
		MaxTokens: 4096,
		System:    system,
		Messages:  []oaiMessage{{Role: "user", Content: user}},
	}
	body, _ := json.Marshal(payload)

	httpCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(httpCtx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", cfg.APIKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var claudeRes claudeResponse
	if err := json.Unmarshal(data, &claudeRes); err != nil {
		return "", fmt.Errorf("parse response: %w", err)
	}
	if claudeRes.Error != nil {
		return "", gerror.New(claudeRes.Error.Message)
	}
	if len(claudeRes.Content) == 0 {
		return "", gerror.New("AI returned no content")
	}
	return claudeRes.Content[0].Text, nil
}
