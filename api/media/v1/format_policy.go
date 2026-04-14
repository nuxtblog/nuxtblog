package v1

import "github.com/gogf/gf/v2/frame/g"

// ── Extension Groups ──────────────────────────────────────────────────────────

type ExtensionGroupItem struct {
	Name       string   `json:"name"`
	LabelZh    string   `json:"label_zh"`
	LabelEn    string   `json:"label_en"`
	Extensions []string `json:"extensions"`
	MaxSizeMB  float64  `json:"max_size_mb"`
}

type ExtensionGroupListReq struct {
	g.Meta `path:"/admin/media/extension-groups" method:"get" tags:"Media" summary:"List extension groups"`
}
type ExtensionGroupListRes struct {
	List []ExtensionGroupItem `json:"list"`
}

type ExtensionGroupSaveReq struct {
	g.Meta `path:"/admin/media/extension-groups" method:"put" tags:"Media" summary:"Save extension groups"`
	Groups []ExtensionGroupItem `v:"required" dc:"full list of extension groups"`
}
type ExtensionGroupSaveRes struct{}

// ── Format Policies ───────────────────────────────────────────────────────────

type FormatPolicyItem struct {
	Name     string   `json:"name"`
	LabelZh  string   `json:"label_zh"`
	LabelEn  string   `json:"label_en"`
	IsSystem bool     `json:"is_system"`
	Groups   []string `json:"groups"`
}

type FormatPolicyListReq struct {
	g.Meta `path:"/admin/media/format-policies" method:"get" tags:"Media" summary:"List format policies"`
}
type FormatPolicyListRes struct {
	List []FormatPolicyItem `json:"list"`
}

type FormatPolicyCreateReq struct {
	g.Meta  `path:"/admin/media/format-policies" method:"post" tags:"Media" summary:"Create format policy"`
	Name    string   `v:"required|max-length:64" dc:"unique policy name"`
	LabelZh string   `v:"max-length:128"         dc:"Chinese label"`
	LabelEn string   `v:"max-length:128"         dc:"English label"`
	Groups  []string `v:"required"               dc:"extension group names to include"`
}
type FormatPolicyCreateRes struct{}

type FormatPolicyUpdateReq struct {
	g.Meta  `path:"/admin/media/format-policies/{name}" method:"put" tags:"Media" summary:"Update format policy"`
	Name    string   `v:"required" dc:"policy name (path param)"`
	LabelZh string   `v:"max-length:128" dc:"Chinese label"`
	LabelEn string   `v:"max-length:128" dc:"English label"`
	Groups  []string `dc:"extension group names to include"`
}
type FormatPolicyUpdateRes struct{}

type FormatPolicyDeleteReq struct {
	g.Meta `path:"/admin/media/format-policies/{name}" method:"delete" tags:"Media" summary:"Delete format policy"`
	Name   string `v:"required" dc:"policy name (path param)"`
}
type FormatPolicyDeleteRes struct{}
