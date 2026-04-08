package smsprovider

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func init() {
	Register(&tencentProvider{})
}

type tencentProvider struct{}

func (t *tencentProvider) Name() string { return "tencent" }

func (t *tencentProvider) Send(ctx context.Context, cfg Config, phone, content string) error {
	const (
		host    = "sms.tencentcloudapi.com"
		service = "sms"
		action  = "SendSms"
		version = "2021-01-11"
		region  = "ap-guangzhou"
	)

	now := time.Now().UTC()
	timestamp := now.Unix()
	date := now.Format("2006-01-02")

	phoneNumbers := []string{phone}
	if !strings.HasPrefix(phone, "+") {
		phoneNumbers = []string{"+" + phone}
	}

	// Use SignName as SmsSdkAppId is not available; TemplateCode as template ID.
	reqBody := map[string]interface{}{
		"SmsSdkAppId":    cfg.SignName,
		"SignName":       cfg.SignName,
		"TemplateId":     cfg.TemplateCode,
		"TemplateParamSet": []string{content},
		"PhoneNumberSet": phoneNumbers,
	}
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("tencent sms: marshal body: %w", err)
	}

	// TC3-HMAC-SHA256 signing
	httpRequestMethod := "POST"
	canonicalURI := "/"
	canonicalQueryString := ""
	canonicalHeaders := "content-type:application/json\nhost:" + host + "\nx-tc-action:" + strings.ToLower(action) + "\n"
	signedHeaders := "content-type;host;x-tc-action"

	hashedPayload := hashSHA256(bodyBytes)
	canonicalRequest := strings.Join([]string{
		httpRequestMethod,
		canonicalURI,
		canonicalQueryString,
		canonicalHeaders,
		signedHeaders,
		hashedPayload,
	}, "\n")

	credentialScope := fmt.Sprintf("%s/%s/tc3_request", date, service)
	stringToSign := strings.Join([]string{
		"TC3-HMAC-SHA256",
		fmt.Sprintf("%d", timestamp),
		credentialScope,
		hashSHA256([]byte(canonicalRequest)),
	}, "\n")

	secretDate := hmacSHA256([]byte("TC3"+cfg.AccessKeySecret), []byte(date))
	secretService := hmacSHA256(secretDate, []byte(service))
	secretSigning := hmacSHA256(secretService, []byte("tc3_request"))
	sig := hex.EncodeToString(hmacSHA256(secretSigning, []byte(stringToSign)))

	authorization := fmt.Sprintf(
		"TC3-HMAC-SHA256 Credential=%s/%s, SignedHeaders=%s, Signature=%s",
		cfg.AccessKeyID, credentialScope, signedHeaders, sig,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://"+host+"/", bytes.NewReader(bodyBytes))
	if err != nil {
		return fmt.Errorf("tencent sms: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Host", host)
	req.Header.Set("X-TC-Action", action)
	req.Header.Set("X-TC-Version", version)
	req.Header.Set("X-TC-Timestamp", fmt.Sprintf("%d", timestamp))
	req.Header.Set("X-TC-Region", region)
	req.Header.Set("Authorization", authorization)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("tencent sms: http: %w", err)
	}
	defer resp.Body.Close()

	raw, _ := io.ReadAll(resp.Body)
	var result struct {
		Response struct {
			SendStatusSet []struct {
				Code    string `json:"Code"`
				Message string `json:"Message"`
			} `json:"SendStatusSet"`
			Error *struct {
				Code    string `json:"Code"`
				Message string `json:"Message"`
			} `json:"Error"`
		} `json:"Response"`
	}
	if err := json.Unmarshal(raw, &result); err != nil {
		return fmt.Errorf("tencent sms: parse response: %w", err)
	}
	if result.Response.Error != nil {
		return fmt.Errorf("tencent sms: %s: %s", result.Response.Error.Code, result.Response.Error.Message)
	}
	for _, s := range result.Response.SendStatusSet {
		if s.Code != "Ok" {
			return fmt.Errorf("tencent sms: %s: %s", s.Code, s.Message)
		}
	}
	return nil
}

func hashSHA256(data []byte) string {
	h := sha256.Sum256(data)
	return hex.EncodeToString(h[:])
}

func hmacSHA256(key, data []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(data)
	return mac.Sum(nil)
}
