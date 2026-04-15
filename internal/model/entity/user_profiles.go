// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// UserProfiles is the golang structure for table user_profiles.
type UserProfiles struct {
	UserId            int    `json:"userId"            orm:"user_id"            description:""` //
	Website           string `json:"website"           orm:"website"            description:""` //
	Twitter           string `json:"twitter"           orm:"twitter"            description:""` //
	Github            string `json:"github"            orm:"github"             description:""` //
	Location          string `json:"location"          orm:"location"           description:""` //
	SocialLinks       string `json:"socialLinks"       orm:"social_links"       description:""` //
	NotificationPrefs string `json:"notificationPrefs" orm:"notification_prefs" description:""` //
}
