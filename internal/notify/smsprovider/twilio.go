package smsprovider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func init() {
	Register(&twilioProvider{})
}

type twilioProvider struct{}

func (t *twilioProvider) Name() string { return "twilio" }

func (t *twilioProvider) Send(ctx context.Context, cfg Config, phone, content string) error {
	accountSid := cfg.AccessKeyID
	authToken := cfg.AccessKeySecret
	from := cfg.SignName

	endpoint := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", accountSid)

	formData := url.Values{}
	formData.Set("To", phone)
	formData.Set("From", from)
	formData.Set("Body", content)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(formData.Encode()))
	if err != nil {
		return fmt.Errorf("twilio sms: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(accountSid, authToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("twilio sms: http: %w", err)
	}
	defer resp.Body.Close()

	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var errResp struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		}
		if json.Unmarshal(raw, &errResp) == nil && errResp.Message != "" {
			return fmt.Errorf("twilio sms: %d: %s", errResp.Code, errResp.Message)
		}
		return fmt.Errorf("twilio sms: HTTP %d", resp.StatusCode)
	}
	return nil
}
