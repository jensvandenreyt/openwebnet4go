package message

import "fmt"

// CENPressure is a marker interface for button pressure types.
type CENPressure interface {
	isCENPressure()
}

// --- CEN Scenario (WHO=15) ---

// WhatCEN represents WHAT values for CEN Scenario messages.
type WhatCEN int

func (w WhatCEN) Value() int { return int(w) }

// WhatCEN button constants (0-31).
const (
	WhatCENButton00 WhatCEN = 0
	WhatCENButton01 WhatCEN = 1
	WhatCENButton02 WhatCEN = 2
	WhatCENButton03 WhatCEN = 3
	WhatCENButton04 WhatCEN = 4
	WhatCENButton05 WhatCEN = 5
	WhatCENButton06 WhatCEN = 6
	WhatCENButton07 WhatCEN = 7
	WhatCENButton08 WhatCEN = 8
	WhatCENButton09 WhatCEN = 9
	WhatCENButton10 WhatCEN = 10
	WhatCENButton11 WhatCEN = 11
	WhatCENButton12 WhatCEN = 12
	WhatCENButton13 WhatCEN = 13
	WhatCENButton14 WhatCEN = 14
	WhatCENButton15 WhatCEN = 15
	WhatCENButton16 WhatCEN = 16
	WhatCENButton17 WhatCEN = 17
	WhatCENButton18 WhatCEN = 18
	WhatCENButton19 WhatCEN = 19
	WhatCENButton20 WhatCEN = 20
	WhatCENButton21 WhatCEN = 21
	WhatCENButton22 WhatCEN = 22
	WhatCENButton23 WhatCEN = 23
	WhatCENButton24 WhatCEN = 24
	WhatCENButton25 WhatCEN = 25
	WhatCENButton26 WhatCEN = 26
	WhatCENButton27 WhatCEN = 27
	WhatCENButton28 WhatCEN = 28
	WhatCENButton29 WhatCEN = 29
	WhatCENButton30 WhatCEN = 30
	WhatCENButton31 WhatCEN = 31
)

var whatCENValues = map[int]WhatCEN{}

func init() {
	for i := 0; i <= 31; i++ {
		whatCENValues[i] = WhatCEN(i)
	}
}

func WhatCENFromValue(i int) What {
	v, ok := whatCENValues[i]
	if !ok {
		return nil
	}
	return v
}

// CENPressureType represents CEN button pressure types.
type CENPressureType int

const (
	CENPressureStart           CENPressureType = 0
	CENPressureReleaseShort    CENPressureType = 1
	CENPressureReleaseExtended CENPressureType = 2
	CENPressureExtended        CENPressureType = 3
)

func (p CENPressureType) isCENPressure() {}

func CENPressureFromValue(i int) (CENPressureType, bool) {
	switch i {
	case 0, 1, 2, 3:
		return CENPressureType(i), true
	}
	return 0, false
}

func newCENScenarioMessage(frame string) *BaseOpenMessage {
	msg := NewBaseOpenMessage(frame)
	msg.WhoField = WhoCENScenarioScheduler
	msg.WhatFromValue = WhatCENFromValue
	msg.DimFromValue = func(i int) Dim { return nil } // no dims for CEN
	msg.ParseWhere = parseCENWhere
	msg.DetectDeviceTyp = func(m *BaseOpenMessage) (OpenDeviceType, error) {
		if !m.IsCommand() {
			return 0, nil
		}
		return DeviceScenarioControl, nil
	}
	return msg
}

func parseCENWhere(whereStr string) (Where, error) {
	if whereStr == "" {
		return nil, NewFrameError("Frame has no WHERE part")
	}
	return NewWhereLightAutom(whereStr)
}

func cenWhatFromButton(button int) (string, error) {
	if button < 0 || button > 31 {
		return "", fmt.Errorf("button number must be between 0 and 31")
	}
	if button < 10 {
		return fmt.Sprintf("0%d", button), nil
	}
	return fmt.Sprintf("%d", button), nil
}

// CENScenarioVirtualStartPressure creates a Virtual Start Pressure message (*15*BUTTON*WHERE##).
func CENScenarioVirtualStartPressure(where string, buttonNumber int) (*BaseOpenMessage, error) {
	ws, err := cenWhatFromButton(buttonNumber)
	if err != nil {
		return nil, err
	}
	frame := fmt.Sprintf(FormatRequestWhatStr, WhoCENScenarioScheduler.Value(), ws, where)
	return newCENScenarioMessage(frame), nil
}

// CENScenarioVirtualReleaseShortPressure creates a Virtual Release after Short Pressure (*15*BUTTON#1*WHERE##).
func CENScenarioVirtualReleaseShortPressure(where string, buttonNumber int) (*BaseOpenMessage, error) {
	ws, err := cenWhatFromButton(buttonNumber)
	if err != nil {
		return nil, err
	}
	frame := fmt.Sprintf(FormatRequestParamStr, WhoCENScenarioScheduler.Value(), ws, int(CENPressureReleaseShort), where)
	return newCENScenarioMessage(frame), nil
}

// CENScenarioVirtualExtendedPressure creates a Virtual Extended Pressure (*15*BUTTON#3*WHERE##).
func CENScenarioVirtualExtendedPressure(where string, buttonNumber int) (*BaseOpenMessage, error) {
	ws, err := cenWhatFromButton(buttonNumber)
	if err != nil {
		return nil, err
	}
	frame := fmt.Sprintf(FormatRequestParamStr, WhoCENScenarioScheduler.Value(), ws, int(CENPressureExtended), where)
	return newCENScenarioMessage(frame), nil
}

// CENScenarioVirtualReleaseExtendedPressure creates a Virtual Release after Extended Pressure (*15*BUTTON#2*WHERE##).
func CENScenarioVirtualReleaseExtendedPressure(where string, buttonNumber int) (*BaseOpenMessage, error) {
	ws, err := cenWhatFromButton(buttonNumber)
	if err != nil {
		return nil, err
	}
	frame := fmt.Sprintf(FormatRequestParamStr, WhoCENScenarioScheduler.Value(), ws, int(CENPressureReleaseExtended), where)
	return newCENScenarioMessage(frame), nil
}

// CENScenarioGetButtonNumber returns the button number from a CEN message.
func CENScenarioGetButtonNumber(msg *BaseOpenMessage) (int, error) {
	w := msg.GetWhat()
	if w == nil {
		return -1, NewFrameError("invalid WHAT in frame")
	}
	return w.Value(), nil
}

// CENScenarioGetButtonPressure returns the button pressure from a CEN message.
func CENScenarioGetButtonPressure(msg *BaseOpenMessage) (CENPressureType, error) {
	params, err := msg.GetCommandParams()
	if err != nil {
		return 0, err
	}
	if len(params) == 0 {
		return CENPressureStart, nil
	}
	p, ok := CENPressureFromValue(params[0])
	if !ok {
		return 0, NewFrameError("unknown pressure type")
	}
	return p, nil
}
