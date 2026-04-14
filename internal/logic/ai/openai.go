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

type oaiMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type oaiRequest struct {
	Model       string       `json:"model"`
	Messages    []oaiMessage `json:"messages"`
	MaxTokens   int          `json:"max_tokens,omitempty"`
	Temperature float64      `json:"temperature"`
}

type oaiResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func openaiGenerate(ctx context.Context, cfg v1.AIConfig, system, user string, timeout time.Duration) (string, error) {
	base := cfg.BaseURL
	if base == "" {
		return "", gerror.New("base_url is required — please set it in Admin → AI config")
	}
	base = strings.TrimRight(base, "/")
	endpoint := base + "/v1/chat/completions"

	msgs := []oaiMessage{}
	if system != "" {
		msgs = append(msgs, oaiMessage{Role: "system", Content: system})
	}
	msgs = append(msgs, oaiMessage{Role: "user", Content: user})

	body, _ := json.Marshal(oaiRequest{
		Model:       cfg.Model,
		Messages:    msgs,
		Temperature: 0.7,
	})

	httpCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(httpCtx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+cfg.APIKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var oaiRes oaiResponse
	if err := json.Unmarshal(data, &oaiRes); err != nil {
		return "", fmt.Errorf("parse response: %w", err)
	}
	if oaiRes.Error != nil {
		return "", gerror.New(oaiRes.Error.Message)
	}
	if len(oaiRes.Choices) == 0 {
		return "", gerror.New("AI returned no content")
	}
	return oaiRes.Choices[0].Message.Content, nil
}
