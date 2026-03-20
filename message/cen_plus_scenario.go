package message

import "fmt"

// WhatCENPlus represents WHAT values for CEN+ messages (WHO=25).
type WhatCENPlus int

const (
	WhatCENPlusShortPressure      WhatCENPlus = 21
	WhatCENPlusStartExtPressure   WhatCENPlus = 22
	WhatCENPlusExtPressure        WhatCENPlus = 23
	WhatCENPlusReleaseExtPressure WhatCENPlus = 24
	WhatCENPlusOnIRDetection      WhatCENPlus = 31
	WhatCENPlusOffIRNoDetection   WhatCENPlus = 32
)

func (w WhatCENPlus) Value() int { return int(w) }

var whatCENPlusValues = map[int]WhatCENPlus{
	21: WhatCENPlusShortPressure, 22: WhatCENPlusStartExtPressure,
	23: WhatCENPlusExtPressure, 24: WhatCENPlusReleaseExtPressure,
	31: WhatCENPlusOnIRDetection, 32: WhatCENPlusOffIRNoDetection,
}

func WhatCENPlusFromValue(i int) What {
	v, ok := whatCENPlusValues[i]
	if !ok {
		return nil
	}
	return v
}

// CENPlusPressureType represents CEN+ button pressure types.
type CENPlusPressureType int

const (
	CENPlusPressureShort           CENPlusPressureType = 21
	CENPlusPressureStartExtended   CENPlusPressureType = 22
	CENPlusPressureExtended        CENPlusPressureType = 23
	CENPlusPressureReleaseExtended CENPlusPressureType = 24
)

func (p CENPlusPressureType) isCENPressure() {}

func newCENPlusScenarioMessage(frame string) *BaseOpenMessage {
	msg := NewBaseOpenMessage(frame)
	msg.WhoField = WhoCENPlusScenarioScheduler
	msg.WhatFromValue = WhatCENPlusFromValue
	msg.DimFromValue = func(i int) Dim { return nil } // no dims for CEN+
	msg.ParseWhere = parseCENPlusWhere
	msg.DetectDeviceTyp = detectCENPlusDeviceType
	return msg
}

func parseCENPlusWhere(whereStr string) (Where, error) {
	if whereStr == "" {
		return nil, NewFrameError("Frame has no WHERE part")
	}
	return NewWhereLightAutom(whereStr)
}

func detectCENPlusDeviceType(msg *BaseOpenMessage) (OpenDeviceType, error) {
	if !msg.IsCommand() {
		return 0, nil
	}
	if CENPlusIsDryContactIR(msg) {
		return DeviceSCSDryContactIR, nil
	}
	return DeviceMultifunctionScenarioControl, nil
}

// CENPlusIsDryContactIR returns true if the message is a dry contact/IR sensor message.
func CENPlusIsDryContactIR(msg *BaseOpenMessage) bool {
	w := msg.GetWhat()
	if w == nil {
		return false
	}
	return w.Value() == WhatCENPlusOnIRDetection.Value() || w.Value() == WhatCENPlusOffIRNoDetection.Value()
}

// CENPlusIsOn returns true if the CEN+ message is ON (WHAT=31).
func CENPlusIsOn(msg *BaseOpenMessage) (bool, error) {
	w := msg.GetWhat()
	if w == nil {
		return false, NewFrameError("invalid WHAT in frame")
	}
	return w.Value() == WhatCENPlusOnIRDetection.Value(), nil
}

// CENPlusIsOff returns true if the CEN+ message is OFF (WHAT=32).
func CENPlusIsOff(msg *BaseOpenMessage) (bool, error) {
	w := msg.GetWhat()
	if w == nil {
		return false, NewFrameError("invalid WHAT in frame")
	}
	return w.Value() == WhatCENPlusOffIRNoDetection.Value(), nil
}

// CENPlusRequestStatus creates a request for CEN+ status.
func CENPlusRequestStatus(where string) *BaseOpenMessage {
	frame := fmt.Sprintf(FormatStatus, WhoCENPlusScenarioScheduler.Value(), where)
	return newCENPlusScenarioMessage(frame)
}

// CENPlusVirtualShortPressure creates a Virtual Short Pressure message (*25*21#BUTTON*WHERE##).
func CENPlusVirtualShortPressure(where string, buttonNumber int) *BaseOpenMessage {
	frame := fmt.Sprintf(FormatRequestParamStr, WhoCENPlusScenarioScheduler.Value(),
		fmt.Sprintf("%d", WhatCENPlusShortPressure.Value()), buttonNumber, where)
	return newCENPlusScenarioMessage(frame)
}

// CENPlusVirtualStartExtendedPressure creates a Virtual Start Extended Pressure (*25*22#BUTTON*WHERE##).
func CENPlusVirtualStartExtendedPressure(where string, buttonNumber int) *BaseOpenMessage {
	frame := fmt.Sprintf(FormatRequestParamStr, WhoCENPlusScenarioScheduler.Value(),
		fmt.Sprintf("%d", WhatCENPlusStartExtPressure.Value()), buttonNumber, where)
	return newCENPlusScenarioMessage(frame)
}

// CENPlusVirtualExtendedPressure creates a Virtual Extended Pressure (*25*23#BUTTON*WHERE##).
func CENPlusVirtualExtendedPressure(where string, buttonNumber int) *BaseOpenMessage {
	frame := fmt.Sprintf(FormatRequestParamStr, WhoCENPlusScenarioScheduler.Value(),
		fmt.Sprintf("%d", WhatCENPlusExtPressure.Value()), buttonNumber, where)
	return newCENPlusScenarioMessage(frame)
}

// CENPlusVirtualReleaseExtendedPressure creates a Virtual Release after Extended Pressure (*25*24#BUTTON*WHERE##).
func CENPlusVirtualReleaseExtendedPressure(where string, buttonNumber int) *BaseOpenMessage {
	frame := fmt.Sprintf(FormatRequestParamStr, WhoCENPlusScenarioScheduler.Value(),
		fmt.Sprintf("%d", WhatCENPlusReleaseExtPressure.Value()), buttonNumber, where)
	return newCENPlusScenarioMessage(frame)
}

// CENPlusGetButtonNumber returns the button number from a CEN+ message.
func CENPlusGetButtonNumber(msg *BaseOpenMessage) (int, error) {
	w := msg.GetWhat()
	if w == nil {
		return -1, NewFrameError("invalid WHAT in frame")
	}
	if w.Value() == WhatCENPlusOffIRNoDetection.Value() || w.Value() == WhatCENPlusOnIRDetection.Value() {
		return -1, nil // dry contact/IR messages don't have button numbers
	}
	params, err := msg.GetCommandParams()
	if err != nil {
		return -1, err
	}
	if len(params) > 0 {
		return params[0], nil
	}
	return -1, nil
}

// CENPlusGetButtonPressure returns the button pressure from a CEN+ message.
func CENPlusGetButtonPressure(msg *BaseOpenMessage) (CENPlusPressureType, bool) {
	w := msg.GetWhat()
	if w == nil {
		return 0, false
	}
	if w.Value() == WhatCENPlusOffIRNoDetection.Value() || w.Value() == WhatCENPlusOnIRDetection.Value() {
		return 0, false
	}
	switch w.Value() {
	case WhatCENPlusShortPressure.Value():
		return CENPlusPressureShort, true
	case WhatCENPlusStartExtPressure.Value():
		return CENPlusPressureStartExtended, true
	case WhatCENPlusExtPressure.Value():
		return CENPlusPressureExtended, true
	case WhatCENPlusReleaseExtPressure.Value():
		return CENPlusPressureReleaseExtended, true
	}
	return 0, false
}
