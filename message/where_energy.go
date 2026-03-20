package message

import (
	"fmt"
	"strconv"
	"strings"
)

// WhereEnergyManagement is WHERE for Energy Management frames.
type WhereEnergyManagement struct {
	BaseWhere
}

var WhereEnergyManagementGeneral = &WhereEnergyManagement{BaseWhere{WhereStr: "0"}}

func NewWhereEnergyManagement(w string) (*WhereEnergyManagement, error) {
	bw, err := NewBaseWhere(w)
	if err != nil {
		return nil, err
	}
	we := &WhereEnergyManagement{BaseWhere: *bw}

	if len(w) == 0 {
		return nil, fmt.Errorf("WHERE address '%s' is invalid", w)
	}

	switch w[0] {
	case '0': // GENERAL
		// OK
	case '1':
		n, err := strconv.Atoi(w[1:])
		if err != nil || n < 1 || n > 127 {
			return nil, fmt.Errorf("WHERE address '%s' is invalid: not in range [1-127]", w)
		}
	case '5':
		n, err := strconv.Atoi(w[1:])
		if err != nil || n < 1 || n > 255 {
			return nil, fmt.Errorf("WHERE address '%s' is invalid: not in range [1-255]", w)
		}
	case '7':
		if !strings.HasSuffix(w, "#0") {
			return nil, fmt.Errorf("WHERE address '%s' is invalid: missing '#0' trailer", w)
		}
		n, err := strconv.Atoi(w[1 : len(w)-2])
		if err != nil || n < 1 || n > 255 {
			return nil, fmt.Errorf("WHERE address '%s' is invalid: not in range [1-255]", w)
		}
	default:
		return nil, fmt.Errorf("WHERE address '%s' is invalid", w)
	}

	return we, nil
}
