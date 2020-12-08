package ext

import (
	"github.com/vdongchina/ratgo/utils/types"
)

type Abstract struct {
	config types.AnyMap
}

// Init config.
func (a *Abstract) Init(config map[string]interface{}) {
	a.config = config
}