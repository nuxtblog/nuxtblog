package entity

type NavMenuItems struct {
	Id         int    `json:"id"         orm:"id"`
	MenuId     int    `json:"menuId"     orm:"menu_id"`
	ParentId   int    `json:"parentId"   orm:"parent_id"`
	ObjectType string `json:"objectType" orm:"object_type"`
	ObjectId   int    `json:"objectId"   orm:"object_id"`
	Label      string `json:"label"      orm:"label"`
	Url        string `json:"url"        orm:"url"`
	Target     string `json:"target"     orm:"target"`
	CssClasses string `json:"cssClasses" orm:"css_classes"`
	SortOrder  int    `json:"sortOrder"  orm:"sort_order"`
}
