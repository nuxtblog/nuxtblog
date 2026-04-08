// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Options is the golang structure for table options.
type Options struct {
	Key       string      `json:"key"       orm:"key"        description:""` //
	Value     string      `json:"value"     orm:"value"      description:""` //
	Autoload  int         `json:"autoload"  orm:"autoload"   description:""` //
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" description:""` //
}
