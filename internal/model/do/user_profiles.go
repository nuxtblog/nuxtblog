// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// UserProfiles is the golang structure of table user_profiles for DAO operations like Where/Data.
type UserProfiles struct {
	g.Meta            `orm:"table:user_profiles, do:true"`
	UserId            any //
	Website           any //
	Twitter           any //
	Github            any //
	Location          any //
	SocialLinks       any //
	NotificationPrefs any //
}
