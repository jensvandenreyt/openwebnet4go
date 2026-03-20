package message

import (
	"fmt"
	"strings"
)

// WhatEnergyMgmt represents the WHAT values for Energy Management messages (WHO=18).
type WhatEnergyMgmt int

const (
	WhatEnergyMgmtAutoResetOn  WhatEnergyMgmt = 26
	WhatEnergyMgmtAutoResetOff WhatEnergyMgmt = 27
)

func (w WhatEnergyMgmt) Value() int { return int(w) }

var whatEnergyMgmtValues = map[int]WhatEnergyMgmt{
	26: WhatEnergyMgmtAutoResetOn,
	27: WhatEnergyMgmtAutoResetOff,
}

func WhatEnergyMgmtFromValue(i int) What {
	v, ok := whatEnergyMgmtValues[i]
	if !ok {
		return nil
	}
	return v
}

// DimEnergyMgmt represents DIM values for Energy Management messages.
type DimEnergyMgmt int

const (
	DimEnergyMgmtActivePower          DimEnergyMgmt = 113
	DimEnergyMgmtActivePowerNotifTime DimEnergyMgmt = 1200
)

func (d DimEnergyMgmt) Value() int { return int(d) }

var dimEnergyMgmtValues = map[int]DimEnergyMgmt{
	113:  DimEnergyMgmtActivePower,
	1200: DimEnergyMgmtActivePowerNotifTime,
}

func DimEnergyMgmtFromValue(i int) Dim {
	v, ok := dimEnergyMgmtValues[i]
	if !ok {
		return nil
	}
	return v
}

func newEnergyManagementMessage(frame string) *BaseOpenMessage {
	msg := NewBaseOpenMessage(frame)
	msg.WhoField = WhoEnergyManagement
	msg.WhatFromValue = WhatEnergyMgmtFromValue
	msg.DimFromValue = DimEnergyMgmtFromValue
	msg.ParseWhere = parseEnergyManagementWhere
	msg.DetectDeviceTyp = detectEnergyManagementDeviceType
	return msg
}

func parseEnergyManagementWhere(whereStr string) (Where, error) {
	if whereStr == "" {
		return nil, NewFrameError("Frame has no WHERE part")
	}
	if strings.HasSuffix(whereStr, ZBNetwork) {
		return NewWhereZigBee(whereStr)
	}
	return NewWhereEnergyManagement(whereStr)
}

func detectEnergyManagementDeviceType(msg *BaseOpenMessage) (OpenDeviceType, error) {
	w := msg.GetWhere()
	if w != nil && strings.HasPrefix(w.Value(), "5") {
		return DeviceSCSEnergyMeter, nil
	}
	return 0, nil
}

// --- Energy Management request constructors ---

// EnergyMgmtRequestActivePower creates a request for active power (*#18*WHERE*113##).
func EnergyMgmtRequestActivePower(where string) *BaseOpenMessage {
	frame := fmt.Sprintf(FormatDimensionRequest, WhoEnergyManagement.Value(), where, DimEnergyMgmtActivePower.Value())
	return newEnergyManagementMessage(frame)
}

// EnergyMgmtSetActivePowerNotificationsTime creates a message to set active power notification time.
func EnergyMgmtSetActivePowerNotificationsTime(where string, time int) *BaseOpenMessage {
	t := time
	if t < 0 || t > 255 {
		t = 0
	}
	frame := fmt.Sprintf(FormatDimensionWriting1P1V,
		WhoEnergyManagement.Value(), where,
		DimEnergyMgmtActivePowerNotifTime.Value(), "1", fmt.Sprintf("%d", t))
	return newEnergyManagementMessage(frame)
}
