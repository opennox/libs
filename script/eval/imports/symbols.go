package imports

import (
	"reflect"
)

//go:generate yaegi extract github.com/opennox/libs/types
//go:generate yaegi extract github.com/opennox/libs/object
//go:generate yaegi extract github.com/opennox/libs/wall
//go:generate yaegi extract github.com/opennox/libs/player
//go:generate yaegi extract github.com/opennox/libs/script

//go:generate goimports -w .

var Symbols = make(map[string]map[string]reflect.Value)
