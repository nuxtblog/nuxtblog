package service

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/verifycode/v1"
)

type IVerifycode interface {
	Send(ctx context.Context, req *v1.VerifycodeSendReq) (*v1.VerifycodeSendRes, error)
	// CheckCode validates a code submitted during registration.
	CheckCode(ctx context.Context, target, codeType, code string) error
	// GetRegisterVerifyMode reads auth_register_verify.mode from options ("none"|"email"|"sms").
	GetRegisterVerifyMode(ctx context.Context) string
}

var localVerifycode IVerifycode

func Verifycode() IVerifycode {
	if localVerifycode == nil {
		panic("implement not found for interface IVerifycode, forgot register?")
	}
	return localVerifycode
}

func RegisterVerifycode(i IVerifycode) {
	localVerifycode = i
}
