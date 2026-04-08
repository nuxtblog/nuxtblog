package entity

type MomentMedia struct {
	MomentId  int `json:"momentId"  orm:"moment_id"`
	MediaId   int `json:"mediaId"   orm:"media_id"`
	SortOrder int `json:"sortOrder" orm:"sort_order"`
}
