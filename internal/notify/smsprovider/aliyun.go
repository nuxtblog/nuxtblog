package smsprovider

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

func init() {
	Register(&aliyunProvider{})
}

type aliyunProvider struct{}

func (a *aliyunProvider) Name() string { return "aliyun" }

func (a *aliyunProvider) Send(ctx context.Context, cfg Config, phone, content string) error {
	nonce := randomNonce()
	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05Z")

	templateParam := fmt.Sprintf(`{"code":"%s"}`, content)

	params := map[string]string{
		"Action":           "SendSms",
		"Version":          "2017-05-25",
		"Format":           "JSON",
		"AccessKeyId":      cfg.AccessKeyID,
		"SignatureMethod":  "HMAC-SHA256",
		"SignatureVersion": "1.0",
		"SignatureNonce":   nonce,
		"Timestamp":        timestamp,
		"PhoneNumbers":     phone,
		"SignName":         cfg.SignName,
		"TemplateCode":     cfg.TemplateCode,
		"TemplateParam":    templateParam,
	}

	// Build sorted query string for signing
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, url.QueryEscape(k)+"="+url.QueryEscape(params[k]))
	}
	canonicalQuery := strings.Join(parts, "&")

	stringToSign := "POST&%2F&" + url.QueryEscape(canonicalQuery)

	mac := hmac.New(sha256.New, []byte(cfg.AccessKeySecret+"&"))
	mac.Write([]byte(stringToSign))
	sig := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	params["Signature"] = sig

	formParts := make([]string, 0, len(params))
	for k, v := range params {
		formParts = append(formParts, url.QueryEscape(k)+"="+url.QueryEscape(v))
	}
	body := strings.Join(formParts, "&")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://dysmsapi.aliyuncs.com/", strings.NewReader(body))
	if err != nil {
		return fmt.Errorf("aliyun sms: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("aliyun sms: http: %w", err)
	}
	defer resp.Body.Close()

	raw, _ := io.ReadAll(resp.Body)
	var result struct {
		Code    string `json:"Code"`
		Message string `json:"Message"`
	}
	if err := json.Unmarshal(raw, &result); err != nil {
		return fmt.Errorf("aliyun sms: parse response: %w", err)
	}
	if result.Code != "OK" {
		return fmt.Errorf("aliyun sms: %s: %s", result.Code, result.Message)
	}
	return nil
}

func randomNonce() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
