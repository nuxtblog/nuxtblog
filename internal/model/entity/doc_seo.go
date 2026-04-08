package entity

type DocSeo struct {
	DocId          int    `json:"docId"          orm:"doc_id"`
	MetaTitle      string `json:"metaTitle"      orm:"meta_title"`
	MetaDesc       string `json:"metaDesc"       orm:"meta_desc"`
	OgTitle        string `json:"ogTitle"        orm:"og_title"`
	OgImage        string `json:"ogImage"        orm:"og_image"`
	CanonicalUrl   string `json:"canonicalUrl"   orm:"canonical_url"`
	Robots         string `json:"robots"         orm:"robots"`
	StructuredData string `json:"structuredData" orm:"structured_data"`
}
