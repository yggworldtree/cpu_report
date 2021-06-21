package comm

import "github.com/yggworldtree/go-core/utils"

type ParamWarn struct {
	WarnVal  float64
	WarnTips string
	Tmr      *utils.Timer `json:"-"`
}
