// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// PostSeo is the golang structure for table post_seo.
type PostSeo struct {
	PostId         int    `json:"postId"         orm:"post_id"         description:""` //
	MetaTitle      string `json:"metaTitle"      orm:"meta_title"      description:""` //
	MetaDesc       string `json:"metaDesc"       orm:"meta_desc"       description:""` //
	OgTitle        string `json:"ogTitle"        orm:"og_title"        description:""` //
	OgImage        string `json:"ogImage"        orm:"og_image"        description:""` //
	CanonicalUrl   string `json:"canonicalUrl"   orm:"canonical_url"   description:""` //
	Robots         string `json:"robots"         orm:"robots"          description:""` //
	StructuredData string `json:"structuredData" orm:"structured_data" description:""` //
}
