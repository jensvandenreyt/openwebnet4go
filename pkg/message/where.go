package message

import "fmt"

// Where is the base type for the WHERE part of an OpenWebNet frame.
type Where interface {
	Value() string
	String() string
}

// BaseWhere provides the base implementation for WHERE.
type BaseWhere struct {
	WhereStr string
}

// NewBaseWhere creates a new BaseWhere. Returns error if w is empty.
func NewBaseWhere(w string) (*BaseWhere, error) {
	if w == "" {
		return nil, fmt.Errorf("WHERE value is empty")
	}
	return &BaseWhere{WhereStr: w}, nil
}

func (bw *BaseWhere) Value() string {
	return bw.WhereStr
}

func (bw *BaseWhere) String() string {
	return "w:" + bw.WhereStr
}

// WhereLightAutom is WHERE for Lighting and Automation frames.
type WhereLightAutom struct {
	BaseWhere
}

var WhereLightAutomGeneral = &WhereLightAutom{BaseWhere{WhereStr: "0"}}

func NewWhereLightAutom(w string) (*WhereLightAutom, error) {
	bw, err := NewBaseWhere(w)
	if err != nil {
		return nil, err
	}
	return &WhereLightAutom{BaseWhere: *bw}, nil
}

// WhereCEN is WHERE for CEN scenario frames.
type WhereCEN struct {
	BaseWhere
}

func NewWhereCEN(w string) (*WhereCEN, error) {
	bw, err := NewBaseWhere(w)
	if err != nil {
		return nil, err
	}
	return &WhereCEN{BaseWhere: *bw}, nil
}
