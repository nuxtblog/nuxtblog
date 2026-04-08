// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// Taxonomies is the golang structure for table taxonomies.
type Taxonomies struct {
	Id          int    `json:"id"          orm:"id"          description:""` //
	TermId      int    `json:"termId"      orm:"term_id"     description:""` //
	Taxonomy    string `json:"taxonomy"    orm:"taxonomy"    description:""` //
	Description string `json:"description" orm:"description" description:""` //
	ParentId    int    `json:"parentId"    orm:"parent_id"   description:""` //
	PostCount   int    `json:"postCount"   orm:"post_count"  description:""` //
	Extra       string `json:"extra"       orm:"extra"       description:""` //
}
