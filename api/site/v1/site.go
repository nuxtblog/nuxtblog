package v1

import "github.com/gogf/gf/v2/frame/g"

type Language struct {
	Code  string `json:"code"`
	Name  string `json:"name"`
	Label string `json:"label"`
}

type LanguagesReq struct {
	g.Meta `path:"/languages" method:"get" tags:"Site" summary:"Get supported languages"`
}

type LanguagesRes struct {
	List []Language `json:"list"`
}
