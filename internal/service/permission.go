package service

import "context"

type IPermission interface {
	// Can reports whether the given role has the named capability.
	// role is the raw JWT value; custom roles (>4) are resolved internally.
	Can(ctx context.Context, role int, cap string) bool
}

var localPermission IPermission

func Permission() IPermission {
	if localPermission == nil {
		panic("implement not found for interface IPermission, forgot register?")
	}
	return localPermission
}

func RegisterPermission(i IPermission) {
	localPermission = i
}
