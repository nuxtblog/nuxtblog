// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// ObjectTaxonomies is the golang structure for table object_taxonomies.
type ObjectTaxonomies struct {
	ObjectId   int    `json:"objectId"   orm:"object_id"   description:""` //
	ObjectType string `json:"objectType" orm:"object_type" description:""` //
	TaxonomyId int    `json:"taxonomyId" orm:"taxonomy_id" description:""` //
	SortOrder  int    `json:"sortOrder"  orm:"sort_order"  description:""` //
}
