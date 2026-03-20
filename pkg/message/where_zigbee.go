package message

import (
	"fmt"
	"strings"
)

const (
	ZBNetwork = "#9"
)

// WhereZigBee is WHERE for ZigBee Lighting and Automation frames.
type WhereZigBee struct {
	BaseWhere
	unit string // UNIT part of the WHERE address
	addr string // ADDR part of the WHERE address
}

const (
	ZBUnit01  = "01"
	ZBUnit02  = "02"
	ZBUnitAll = "00"
)

func NewWhereZigBee(w string) (*WhereZigBee, error) {
	bw, err := NewBaseWhere(w)
	if err != nil {
		return nil, err
	}
	wz := &WhereZigBee{BaseWhere: *bw}
	lastHash := strings.LastIndex(w, "#")
	if lastHash > 0 && len(w) >= 4 {
		wz.unit = w[len(w)-4 : len(w)-2]
		wz.addr = w[:len(w)-4]
	} else {
		return nil, fmt.Errorf("WHERE address is invalid")
	}
	return wz, nil
}

// ValueWithUnit returns the WHERE value using the provided string as UNIT.
func (wz *WhereZigBee) ValueWithUnit(u string) string {
	return wz.addr + u + ZBNetwork
}

// GetUnit returns the UNIT part (e.g. WHERE=123456702#9 -> UNIT=02).
func (wz *WhereZigBee) GetUnit() string {
	return wz.unit
}

// GetAddr returns the ADDR part by removing UNIT and network.
func (wz *WhereZigBee) GetAddr() string {
	return wz.addr
}
