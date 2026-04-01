package message

import (
	"fmt"
	"math"
	"strings"
)

// WhatLighting represents the WHAT values for Lighting messages (WHO=1).
type WhatLighting int

const (
	WhatLightingOff                 WhatLighting = 0
	WhatLightingOn                  WhatLighting = 1
	WhatLightingDimmerLevel2        WhatLighting = 2
	WhatLightingDimmerLevel3        WhatLighting = 3
	WhatLightingDimmerLevel4        WhatLighting = 4
	WhatLightingDimmerLevel5        WhatLighting = 5
	WhatLightingDimmerLevel6        WhatLighting = 6
	WhatLightingDimmerLevel7        WhatLighting = 7
	WhatLightingDimmerLevel8        WhatLighting = 8
	WhatLightingDimmerLevel9        WhatLighting = 9
	WhatLightingDimmerLevel10       WhatLighting = 10
	WhatLightingDimmerLevelUp       WhatLighting = 30
	WhatLightingDimmerLevelDown     WhatLighting = 31
	WhatLightingDimmerToggle        WhatLighting = 32
	WhatLightingMovementDetected    WhatLighting = 34
	WhatLightingEndMovementDetected WhatLighting = 39
)

func (w WhatLighting) Value() int {
	return int(w)
}

func (w WhatLighting) String() string {
	return whatLightingAsString(w)
}

// WhatLightingFromValue returns the WhatLighting for a given int.
func WhatLightingFromValue(i int) What {
	switch i {
	case WhatLightingOff.Value():
		return WhatLightingOff
	case WhatLightingOn.Value():
		return WhatLightingOn
	case WhatLightingDimmerLevel2.Value():
		return WhatLightingDimmerLevel2
	case WhatLightingDimmerLevel3.Value():
		return WhatLightingDimmerLevel3
	case WhatLightingDimmerLevel4.Value():
		return WhatLightingDimmerLevel4
	case WhatLightingDimmerLevel5.Value():
		return WhatLightingDimmerLevel5
	case WhatLightingDimmerLevel6.Value():
		return WhatLightingDimmerLevel6
	case WhatLightingDimmerLevel7.Value():
		return WhatLightingDimmerLevel7
	case WhatLightingDimmerLevel8.Value():
		return WhatLightingDimmerLevel8
	case WhatLightingDimmerLevel9.Value():
		return WhatLightingDimmerLevel9
	case WhatLightingDimmerLevel10.Value():
		return WhatLightingDimmerLevel10
	case WhatLightingDimmerLevelUp.Value():
		return WhatLightingDimmerLevelUp
	case WhatLightingDimmerLevelDown.Value():
		return WhatLightingDimmerLevelDown
	case WhatLightingDimmerToggle.Value():
		return WhatLightingDimmerToggle
	case WhatLightingMovementDetected.Value():
		return WhatLightingMovementDetected
	case WhatLightingEndMovementDetected.Value():
		return WhatLightingEndMovementDetected
	default:
		return nil
	}
}

func whatLightingAsString(wl WhatLighting) string {
	switch wl {
	case WhatLightingOff:
		return "LightingOff"
	case WhatLightingOn:
		return "LightingOn"
	case WhatLightingDimmerLevel2:
		return "LightingDimmerLevel2"
	case WhatLightingDimmerLevel3:
		return "LightingDimmerLevel3"
	case WhatLightingDimmerLevel4:
		return "LightingDimmerLevel4"
	case WhatLightingDimmerLevel5:
		return "LightingDimmerLevel5"
	case WhatLightingDimmerLevel6:
		return "LightingDimmerLevel6"
	case WhatLightingDimmerLevel7:
		return "LightingDimmerLevel7"
	case WhatLightingDimmerLevel8:
		return "LightingDimmerLevel8"
	case WhatLightingDimmerLevel9:
		return "LightingDimmerLevel9"
	case WhatLightingDimmerLevel10:
		return "LightingDimmerLevel10"
	case WhatLightingDimmerLevelUp:
		return "LightingDimmerLevelUp"
	case WhatLightingDimmerLevelDown:
		return "LightingDimmerLevelDown"
	case WhatLightingDimmerToggle:
		return "LightingDimmerToggle"
	case WhatLightingMovementDetected:
		return "LightingMovementDetected"
	case WhatLightingEndMovementDetected:
		return "LightingEndMovementDetected"
	default:
		return ""
	}
}

// DimLighting represents DIM values for Lighting messages.
type DimLighting int

const (
	DimLightingDimmerLevel100 DimLighting = 1
)

func (d DimLighting) Value() int { return int(d) }

var dimLightingValues = map[int]DimLighting{
	1: DimLightingDimmerLevel100,
}

func DimLightingFromValue(i int) Dim {
	v, ok := dimLightingValues[i]
	if !ok {
		return nil
	}
	return v
}

const (
	DimmerLevel100Off = 100
	DimmerLevel100Max = 200
)

func newLightingMessage(frame string) *BaseOpenMessage {
	msg := NewBaseOpenMessage(frame)
	msg.WhoField = WhoLighting
	msg.WhatFromValue = WhatLightingFromValue
	msg.DimFromValue = DimLightingFromValue
	msg.ParseWhere = parseLightingWhere
	msg.DetectDeviceTyp = detectLightingDeviceType
	return msg
}

func parseLightingWhere(whereStr string) (Where, error) {
	if whereStr == "" {
		return nil, NewFrameError("Lighting frame has no WHERE part")
	}
	if strings.HasSuffix(whereStr, ZBNetwork) {
		return NewWhereZigBee(whereStr)
	}
	return NewWhereLightAutom(whereStr)
}

func detectLightingDeviceType(msg *BaseOpenMessage) (OpenDeviceType, error) {
	if msg.IsCommand() {
		w := msg.GetWhat()
		if w != nil {
			wl, ok := w.(WhatLighting)
			if ok {
				switch wl {
				case WhatLightingOff, WhatLightingOn, WhatLightingMovementDetected, WhatLightingEndMovementDetected:
					return DeviceSCSOnOffSwitch, nil
				default:
					if wl.Value() >= 2 && wl <= 10 {
						return DeviceSCSDimmerSwitch, nil
					}
				}
			}
		}
	}
	return 0, nil
}

// --- Lighting request constructors ---

// LightingRequestTurnOn creates a request to turn light ON (*1*1*WHERE##).
func LightingRequestTurnOn(where string) *BaseOpenMessage {
	frame := fmt.Sprintf(FormatRequest, WhoLighting.Value(), WhatLightingOn.Value(), where)
	return newLightingMessage(frame)
}

// LightingRequestTurnOff creates a request to turn light OFF (*1*0*WHERE##).
func LightingRequestTurnOff(where string) *BaseOpenMessage {
	frame := fmt.Sprintf(FormatRequest, WhoLighting.Value(), WhatLightingOff.Value(), where)
	return newLightingMessage(frame)
}

// LightingRequestDimTo creates a request to dim light to level.
func LightingRequestDimTo(where string, level What) *BaseOpenMessage {
	frame := fmt.Sprintf(FormatRequest, WhoLighting.Value(), level.Value(), where)
	return newLightingMessage(frame)
}

// LightingRequestStatus creates a request for light status (*#1*WHERE##).
func LightingRequestStatus(where string) *BaseOpenMessage {
	frame := fmt.Sprintf(FormatStatus, WhoLighting.Value(), where)
	return newLightingMessage(frame)
}

// LightingIsOn returns true if the lighting message is ON.
func LightingIsOn(msg *BaseOpenMessage) bool {
	w := msg.GetWhat()
	if w == nil {
		return false
	}
	return w.Value() == WhatLightingOn.Value()
}

// LightingIsOff returns true if the lighting message is OFF.
func LightingIsOff(msg *BaseOpenMessage) bool {
	w := msg.GetWhat()
	if w == nil {
		return false
	}
	return w.Value() == WhatLightingOff.Value()
}

// LightingParseDimmerLevel100 parses dimmerLevel100 (DIM: 1) and returns percentage (0-100).
func LightingParseDimmerLevel100(msg *BaseOpenMessage) (int, error) {
	dim := msg.GetDim()
	if dim == nil {
		return 0, NewFrameError("Could not parse dimmerLevel100")
	}
	if dim.Value() != DimLightingDimmerLevel100.Value() {
		return 0, NewFrameError("Could not parse dimmerLevel100")
	}
	vals, err := msg.GetDimValues()
	if err != nil || len(vals) == 0 {
		return 0, NewFrameError("Could not parse dimmerLevel100")
	}
	level100 := 0
	fmt.Sscanf(vals[0], "%d", &level100)
	if level100 >= DimmerLevel100Off && level100 <= DimmerLevel100Max {
		return level100 - 100, nil
	}
	return 0, NewFrameError("Value for dimmerLevel100 out of range")
}

// LightingLevelToPercent transforms a 0-10 level to percent (0-100).
func LightingLevelToPercent(level int) (int, error) {
	if level >= 0 && level <= 10 {
		return level * 10, nil
	}
	return 0, fmt.Errorf("level must be between 0 and 10")
}

// LightingPercentToWhat returns the What corresponding to the brightness percent.
func LightingPercentToWhat(percent int) (What, error) {
	if percent < 0 || percent > 100 {
		return nil, fmt.Errorf("percent must be between 0 and 100")
	}
	level := percentToWhatLevel(percent)
	w := WhatLightingFromValue(level)
	return w, nil
}

func percentToWhatLevel(percent int) int {
	if percent == 0 {
		return 0
	}
	if percent > 0 && percent < 10 {
		return 2
	}
	level := int(math.Floor(float64(percent) / 10.0))
	if level == 1 {
		level++ // level 1 is not allowed -> move to 2
	}
	return level
}
