package v1

import "github.com/gogf/gf/v2/frame/g"

// VerifycodeSendReq — POST /auth/send-code
type VerifycodeSendReq struct {
	g.Meta `path:"/auth/send-code" method:"post" tags:"Verifycode" summary:"Send registration verification code"`
	Type   string `p:"type"   v:"required|in:email,sms" dc:"email or sms"`
	Target string `p:"target" v:"required"              dc:"email address or phone number"`
}

type VerifycodeSendRes struct {
	ExpiresIn int `json:"expires_in" dc:"code validity in seconds"`
}
