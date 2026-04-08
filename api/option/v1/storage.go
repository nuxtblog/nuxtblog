package v1

import "github.com/gogf/gf/v2/frame/g"

// StorageBackendInfo describes one configured storage backend.
type StorageBackendInfo struct {
	Name        string `json:"name"`         // config key, used as ID in routing rules
	DisplayName string `json:"display_name"` // human-readable label
	Type        string `json:"type"`
	Enabled     bool   `json:"enabled"`
}

type StorageBackendsReq struct {
	g.Meta `path:"/admin/storage/backends" method:"get" tags:"Admin" summary:"List configured storage backends"`
}

type StorageBackendsRes struct {
	Backends []StorageBackendInfo `json:"backends"`
	Default  string               `json:"default"`
}
