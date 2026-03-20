package message

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// ThermoFunction represents heating/cooling/generic function.
type ThermoFunction int

const (
	ThermoFunctionHeating ThermoFunction = 1
	ThermoFunctionCooling ThermoFunction = 2
	ThermoFunctionGeneric ThermoFunction = 3
)

func (f ThermoFunction) Value() int { return int(f) }

func ThermoFunctionFromValue(i int) (ThermoFunction, bool) {
	switch i {
	case 1:
		return ThermoFunctionHeating, true
	case 2:
		return ThermoFunctionCooling, true
	case 3:
		return ThermoFunctionGeneric, true
	}
	return 0, false
}

// ThermoOperationMode represents manual/protection/off modes.
type ThermoOperationMode int

const (
	ThermoOperationManual     ThermoOperationMode = 1
	ThermoOperationProtection ThermoOperationMode = 2
	ThermoOperationOff        ThermoOperationMode = 3
)

func (m ThermoOperationMode) Value() int { return int(m) }

func ThermoOperationModeFromValue(i int) (ThermoOperationMode, bool) {
	switch i {
	case 1:
		return ThermoOperationManual, true
	case 2:
		return ThermoOperationProtection, true
	case 3:
		return ThermoOperationOff, true
	}
	return 0, false
}

// WhatThermo represents the WHAT values for Thermoregulation messages (WHO=4).
type WhatThermo int

const (
	WhatThermoConditioning           WhatThermo = 0
	WhatThermoHeating                WhatThermo = 1
	WhatThermoGeneric                WhatThermo = 3
	WhatThermoProtectionHeating      WhatThermo = 102
	WhatThermoProtectionConditioning WhatThermo = 202
	WhatThermoProtectionGeneric      WhatThermo = 302
	WhatThermoOffHeating             WhatThermo = 103
	WhatThermoOffConditioning        WhatThermo = 203
	WhatThermoOffGeneric             WhatThermo = 303
	WhatThermoManualHeating          WhatThermo = 110
	WhatThermoManualConditioning     WhatThermo = 210
	WhatThermoManualGeneric          WhatThermo = 310
	WhatThermoProgramHeating         WhatThermo = 111
	WhatThermoProgramConditioning    WhatThermo = 211
	WhatThermoProgramGeneric         WhatThermo = 311
	WhatThermoHolidayHeating         WhatThermo = 115
	WhatThermoHolidayConditioning    WhatThermo = 215
	WhatThermoHolidayGeneric         WhatThermo = 315
)

func (w WhatThermo) Value() int { return int(w) }

var whatThermoValues = map[int]WhatThermo{
	0: WhatThermoConditioning, 1: WhatThermoHeating, 3: WhatThermoGeneric,
	102: WhatThermoProtectionHeating, 202: WhatThermoProtectionConditioning, 302: WhatThermoProtectionGeneric,
	103: WhatThermoOffHeating, 203: WhatThermoOffConditioning, 303: WhatThermoOffGeneric,
	110: WhatThermoManualHeating, 210: WhatThermoManualConditioning, 310: WhatThermoManualGeneric,
	111: WhatThermoProgramHeating, 211: WhatThermoProgramConditioning, 311: WhatThermoProgramGeneric,
	115: WhatThermoHolidayHeating, 215: WhatThermoHolidayConditioning, 315: WhatThermoHolidayGeneric,
}

func WhatThermoFromValue(i int) What {
	v, ok := whatThermoValues[i]
	if !ok {
		return nil
	}
	return v
}

// ThermoLocalOffset represents local offset values.
type ThermoLocalOffset struct {
	OffsetValue string
	Label       string
}

var (
	ThermoLocalOffsetPlus3      = ThermoLocalOffset{"03", "+3"}
	ThermoLocalOffsetPlus2      = ThermoLocalOffset{"02", "+2"}
	ThermoLocalOffsetPlus1      = ThermoLocalOffset{"01", "+1"}
	ThermoLocalOffsetNormal     = ThermoLocalOffset{"00", "NORMAL"}
	ThermoLocalOffsetMinus1     = ThermoLocalOffset{"11", "-1"}
	ThermoLocalOffsetMinus2     = ThermoLocalOffset{"12", "-2"}
	ThermoLocalOffsetMinus3     = ThermoLocalOffset{"13", "-3"}
	ThermoLocalOffsetOff        = ThermoLocalOffset{"4", "OFF"}
	ThermoLocalOffsetProtection = ThermoLocalOffset{"5", "PROTECTION"}
)

var thermoLocalOffsets = []ThermoLocalOffset{
	ThermoLocalOffsetPlus3, ThermoLocalOffsetPlus2, ThermoLocalOffsetPlus1,
	ThermoLocalOffsetNormal, ThermoLocalOffsetMinus1, ThermoLocalOffsetMinus2,
	ThermoLocalOffsetMinus3, ThermoLocalOffsetOff, ThermoLocalOffsetProtection,
}

// ThermoLocalOffsetFromValue returns the local offset for the given string value.
func ThermoLocalOffsetFromValue(s string) *ThermoLocalOffset {
	for _, lo := range thermoLocalOffsets {
		if lo.OffsetValue == s {
			return &lo
		}
	}
	return nil
}

// ThermoFanCoilSpeed represents fan coil speed values.
type ThermoFanCoilSpeed int

const (
	ThermoFanCoilSpeedAuto   ThermoFanCoilSpeed = 0
	ThermoFanCoilSpeedSpeed1 ThermoFanCoilSpeed = 1
	ThermoFanCoilSpeedSpeed2 ThermoFanCoilSpeed = 2
	ThermoFanCoilSpeedSpeed3 ThermoFanCoilSpeed = 3
)

func ThermoFanCoilSpeedFromValue(i int) (ThermoFanCoilSpeed, bool) {
	switch i {
	case 0:
		return ThermoFanCoilSpeedAuto, true
	case 1:
		return ThermoFanCoilSpeedSpeed1, true
	case 2:
		return ThermoFanCoilSpeedSpeed2, true
	case 3:
		return ThermoFanCoilSpeedSpeed3, true
	}
	return 0, false
}

func (f ThermoFanCoilSpeed) Value() int { return int(f) }

// ThermoValveOrActuatorStatus represents valve/actuator status values.
type ThermoValveOrActuatorStatus int

const (
	ThermoValveOff       ThermoValveOrActuatorStatus = 0
	ThermoValveOn        ThermoValveOrActuatorStatus = 1
	ThermoValveOpened    ThermoValveOrActuatorStatus = 2
	ThermoValveClosed    ThermoValveOrActuatorStatus = 3
	ThermoValveStop      ThermoValveOrActuatorStatus = 4
	ThermoValveOffFC     ThermoValveOrActuatorStatus = 5
	ThermoValveOnSpeed1  ThermoValveOrActuatorStatus = 6
	ThermoValveOnSpeed2  ThermoValveOrActuatorStatus = 7
	ThermoValveOnSpeed3  ThermoValveOrActuatorStatus = 8
	ThermoValveOnFC      ThermoValveOrActuatorStatus = 9
	ThermoValveOffSpeed1 ThermoValveOrActuatorStatus = 14
	ThermoValveOffSpeed2 ThermoValveOrActuatorStatus = 15
	ThermoValveOffSpeed3 ThermoValveOrActuatorStatus = 16
)

func ThermoValveOrActuatorStatusFromValue(i int) (ThermoValveOrActuatorStatus, bool) {
	switch i {
	case 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 14, 15, 16:
		return ThermoValveOrActuatorStatus(i), true
	}
	return 0, false
}

// DimThermo represents DIM values for Thermoregulation messages.
type DimThermo int

const (
	DimThermoTemperature         DimThermo = 0
	DimThermoFanCoilSpeed        DimThermo = 11
	DimThermoCompleteProbeStatus DimThermo = 12
	DimThermoOffset              DimThermo = 13
	DimThermoTempSetpoint        DimThermo = 14
	DimThermoProbeTemperature    DimThermo = 15
	DimThermoValvesStatus        DimThermo = 19
	DimThermoActuatorStatus      DimThermo = 20
)

func (d DimThermo) Value() int { return int(d) }

var dimThermoValues = map[int]DimThermo{
	0: DimThermoTemperature, 11: DimThermoFanCoilSpeed,
	12: DimThermoCompleteProbeStatus, 13: DimThermoOffset,
	14: DimThermoTempSetpoint, 15: DimThermoProbeTemperature,
	19: DimThermoValvesStatus, 20: DimThermoActuatorStatus,
}

func DimThermoFromValue(i int) Dim {
	v, ok := dimThermoValues[i]
	if !ok {
		return nil
	}
	return v
}

func newThermoregulationMessage(frame string) *BaseOpenMessage {
	msg := NewBaseOpenMessage(frame)
	msg.WhoField = WhoThermoregulation
	msg.WhatFromValue = WhatThermoFromValue
	msg.DimFromValue = DimThermoFromValue
	msg.ParseWhere = parseThermoWhere
	msg.DetectDeviceTyp = detectThermoDeviceType
	return msg
}

func parseThermoWhere(whereStr string) (Where, error) {
	if whereStr == "" {
		return nil, NewFrameError("Frame has no WHERE part")
	}
	return NewWhereThermo(whereStr)
}

func detectThermoDeviceType(msg *BaseOpenMessage) (OpenDeviceType, error) {
	w := msg.GetWhere()
	if w == nil {
		return 0, nil
	}
	if strings.HasPrefix(w.Value(), "0") {
		return 0, nil
	}
	wt, ok := w.(*WhereThermo)
	if !ok {
		return 0, nil
	}
	if wt.IsProbe() {
		return DeviceSCSThermoSensor, nil
	} else if wt.IsCentralUnit() {
		return DeviceSCSThermoCentralUnit, nil
	}
	return DeviceSCSThermoZone, nil
}

// --- Thermoregulation request constructors ---

// ThermoRequestWriteSetpointTemperature creates a message to set setpoint temperature.
func ThermoRequestWriteSetpointTemperature(where string, temp float64, function ThermoFunction) (*BaseOpenMessage, error) {
	if temp < 5 || temp > 40 {
		return nil, NewMalformedFrameError("Set Point Temperature should be between 5° and 40° Celsius.")
	}
	rounded := math.Round(temp*2) / 2
	frame := fmt.Sprintf(FormatDimensionWriting2V,
		WhoThermoregulation.Value(), where, DimThermoTempSetpoint.Value(),
		ThermoEncodeTemperature(rounded), fmt.Sprintf("%d", function.Value()))
	return newThermoregulationMessage(frame), nil
}

// ThermoRequestWriteFanCoilSpeed creates a message to set fan coil speed.
func ThermoRequestWriteFanCoilSpeed(where string, speed ThermoFanCoilSpeed) *BaseOpenMessage {
	frame := fmt.Sprintf(FormatDimensionWriting1V,
		WhoThermoregulation.Value(), where, DimThermoFanCoilSpeed.Value(),
		fmt.Sprintf("%d", speed.Value()))
	return newThermoregulationMessage(frame)
}

// ThermoRequestFanCoilSpeed creates a message to request fan coil speed.
func ThermoRequestFanCoilSpeed(where string) *BaseOpenMessage {
	frame := fmt.Sprintf(FormatDimensionRequest,
		WhoThermoregulation.Value(), where, DimThermoFanCoilSpeed.Value())
	return newThermoregulationMessage(frame)
}

// ThermoRequestWriteFunction creates a message to set the function.
func ThermoRequestWriteFunction(where string, function ThermoFunction) *BaseOpenMessage {
	switch function {
	case ThermoFunctionHeating:
		frame := fmt.Sprintf(FormatRequest, WhoThermoregulation.Value(), WhatThermoProtectionHeating.Value(), where)
		return newThermoregulationMessage(frame)
	case ThermoFunctionCooling:
		frame := fmt.Sprintf(FormatRequest, WhoThermoregulation.Value(), WhatThermoProtectionConditioning.Value(), where)
		return newThermoregulationMessage(frame)
	case ThermoFunctionGeneric:
		frame := fmt.Sprintf(FormatRequest, WhoThermoregulation.Value(), WhatThermoProtectionGeneric.Value(), where)
		return newThermoregulationMessage(frame)
	}
	return nil
}

// ThermoRequestWriteMode creates a message to set the operation mode.
func ThermoRequestWriteMode(where string, mode ThermoOperationMode, currentFunc ThermoFunction, setpointTemp float64) *BaseOpenMessage {
	switch mode {
	case ThermoOperationManual:
		msg, err := ThermoRequestWriteSetpointTemperature(where, setpointTemp, currentFunc)
		if err != nil {
			return nil
		}
		return msg
	case ThermoOperationProtection:
		switch currentFunc {
		case ThermoFunctionHeating:
			return newThermoregulationMessage(fmt.Sprintf(FormatRequest, WhoThermoregulation.Value(), WhatThermoProtectionHeating.Value(), where))
		case ThermoFunctionCooling:
			return newThermoregulationMessage(fmt.Sprintf(FormatRequest, WhoThermoregulation.Value(), WhatThermoProtectionConditioning.Value(), where))
		case ThermoFunctionGeneric:
			return newThermoregulationMessage(fmt.Sprintf(FormatRequest, WhoThermoregulation.Value(), WhatThermoProtectionGeneric.Value(), where))
		}
	case ThermoOperationOff:
		switch currentFunc {
		case ThermoFunctionHeating:
			return newThermoregulationMessage(fmt.Sprintf(FormatRequest, WhoThermoregulation.Value(), WhatThermoOffHeating.Value(), where))
		case ThermoFunctionCooling:
			return newThermoregulationMessage(fmt.Sprintf(FormatRequest, WhoThermoregulation.Value(), WhatThermoOffConditioning.Value(), where))
		case ThermoFunctionGeneric:
			return newThermoregulationMessage(fmt.Sprintf(FormatRequest, WhoThermoregulation.Value(), WhatThermoOffGeneric.Value(), where))
		}
	}
	return nil
}

// ThermoRequestMode creates a request for the set point temperature with local offset and operation mode.
func ThermoRequestMode(where string) *BaseOpenMessage {
	frame := fmt.Sprintf(FormatDimensionRequest, WhoThermoregulation.Value(), where, DimThermoCompleteProbeStatus.Value())
	return newThermoregulationMessage(frame)
}

// ThermoRequestValvesStatus creates a request for valves status.
func ThermoRequestValvesStatus(where string) *BaseOpenMessage {
	frame := fmt.Sprintf(FormatDimensionRequest, WhoThermoregulation.Value(), where, DimThermoValvesStatus.Value())
	return newThermoregulationMessage(frame)
}

// ThermoRequestWriteSetMode creates a message to set the zone mode.
func ThermoRequestWriteSetMode(where string, newMode WhatThermo) *BaseOpenMessage {
	frame := fmt.Sprintf(FormatRequest, WhoThermoregulation.Value(), newMode.Value(), where)
	return newThermoregulationMessage(frame)
}

// ThermoRequestTurnOff creates a message to turn off the zone.
func ThermoRequestTurnOff(where string) *BaseOpenMessage {
	frame := fmt.Sprintf(FormatRequest, WhoThermoregulation.Value(), WhatThermoOffGeneric.Value(), where)
	return newThermoregulationMessage(frame)
}

// ThermoRequestTemperature creates a request for current sensed temperature.
func ThermoRequestTemperature(where string) *BaseOpenMessage {
	frame := fmt.Sprintf(FormatDimensionRequest, WhoThermoregulation.Value(), where, DimThermoTemperature.Value())
	return newThermoregulationMessage(frame)
}

// ThermoRequestSetPointTemperature creates a request for current set point temperature.
func ThermoRequestSetPointTemperature(where string) *BaseOpenMessage {
	frame := fmt.Sprintf(FormatDimensionRequest, WhoThermoregulation.Value(), where, DimThermoTempSetpoint.Value())
	return newThermoregulationMessage(frame)
}

// ThermoRequestStatus creates a request for zone status.
func ThermoRequestStatus(where string) *BaseOpenMessage {
	frame := fmt.Sprintf(FormatStatus, WhoThermoregulation.Value(), where)
	return newThermoregulationMessage(frame)
}

// ThermoRequestActuatorsStatus creates a request for actuators status.
func ThermoRequestActuatorsStatus(where string) *BaseOpenMessage {
	frame := fmt.Sprintf(FormatDimensionRequest, WhoThermoregulation.Value(), where, DimThermoActuatorStatus.Value())
	return newThermoregulationMessage(frame)
}

// --- Thermoregulation parse helpers ---

// ThermoParseTemperature parses temperature from a Thermoregulation message.
func ThermoParseTemperature(msg *BaseOpenMessage) (float64, error) {
	vals, err := msg.GetDimValues()
	if err != nil {
		return 0, err
	}
	dim := msg.GetDim()
	if dim == nil {
		return 0, fmt.Errorf("could not parse temperature: no DIM")
	}
	switch dim.Value() {
	case DimThermoTemperature.Value(), DimThermoTempSetpoint.Value(), DimThermoCompleteProbeStatus.Value():
		if len(vals) > 0 {
			return ThermoDecodeTemperature(vals[0])
		}
	case DimThermoProbeTemperature.Value():
		if len(vals) > 1 {
			return ThermoDecodeTemperature(vals[1])
		}
	}
	return 0, fmt.Errorf("could not parse temperature from: %s", msg.GetFrameValue())
}

// ThermoDecodeTemperature converts temperature from BTicino format to number.
// Example: 0235 --> +23.5 (°C) and 1048 --> -4.8 (°C)
func ThermoDecodeTemperature(temperature string) (float64, error) {
	t := temperature
	if len(t) > 0 && t[0] == '#' {
		t = t[1:]
	}
	sign := 1.0
	var tempInt int
	var err error
	if len(t) == 4 {
		if t[0] == '1' {
			sign = -1.0
		}
		tempInt, err = strconv.Atoi(t[1:])
		if err != nil {
			return 0, fmt.Errorf("unrecognized temperature format: %s", t)
		}
	} else if len(t) == 3 {
		tempInt, err = strconv.Atoi(t)
		if err != nil {
			return 0, fmt.Errorf("unrecognized temperature format: %s", t)
		}
	} else {
		return 0, fmt.Errorf("unrecognized temperature format: %s", t)
	}
	tempDouble := sign * float64(tempInt) / 10.0
	return math.Round(tempDouble*100.0) / 100.0, nil
}

// ThermoEncodeTemperature encodes temperature to BTicino format.
// Example: +23.51 °C --> '0235'; -4.86 °C --> '1049'
func ThermoEncodeTemperature(temp float64) string {
	sign := '0'
	if temp < 0 {
		sign = '1'
	}
	absTemp := int(math.Abs(math.Round(temp * 10)))
	digits := ""
	if absTemp < 100 {
		digits += "0"
	}
	if absTemp < 10 {
		digits += "0"
	}
	digits += strconv.Itoa(absTemp)
	if digits == "000" {
		sign = '0'
	}
	return string(sign) + digits
}

// ThermoParseFanCoilSpeed parses fan coil speed from a Thermoregulation message.
func ThermoParseFanCoilSpeed(msg *BaseOpenMessage) (ThermoFanCoilSpeed, error) {
	vals, err := msg.GetDimValues()
	if err != nil {
		return 0, err
	}
	dim := msg.GetDim()
	if dim != nil && dim.Value() == DimThermoFanCoilSpeed.Value() && len(vals) > 0 {
		v, err := strconv.Atoi(vals[0])
		if err != nil {
			return 0, fmt.Errorf("could not parse fan coil speed: %v", err)
		}
		fcs, ok := ThermoFanCoilSpeedFromValue(v)
		if ok {
			return fcs, nil
		}
	}
	return 0, fmt.Errorf("could not parse fan coil speed from: %s", msg.GetFrameValue())
}

// ThermoParseValveStatus parses valve status from a Thermoregulation message (dimension: 19).
func ThermoParseValveStatus(msg *BaseOpenMessage, what WhatThermo) (ThermoValveOrActuatorStatus, error) {
	if what != WhatThermoConditioning && what != WhatThermoHeating {
		return 0, NewFrameError("Only CONDITIONING and HEATING are allowed as what input parameter.")
	}
	vals, err := msg.GetDimValues()
	if err != nil {
		return 0, err
	}
	dim := msg.GetDim()
	if dim != nil && dim.Value() == DimThermoValvesStatus.Value() && len(vals) >= 2 {
		idx := 0
		if what == WhatThermoHeating {
			idx = 1
		}
		v, err := strconv.Atoi(vals[idx])
		if err != nil {
			return 0, fmt.Errorf("could not parse valve status: %v", err)
		}
		vs, ok := ThermoValveOrActuatorStatusFromValue(v)
		if ok {
			return vs, nil
		}
	}
	return 0, fmt.Errorf("could not parse valve status from: %s", msg.GetFrameValue())
}

// ThermoParseActuatorStatus parses actuator status from a Thermoregulation message (dimension: 20).
func ThermoParseActuatorStatus(msg *BaseOpenMessage) (ThermoValveOrActuatorStatus, error) {
	vals, err := msg.GetDimValues()
	if err != nil {
		return 0, err
	}
	dim := msg.GetDim()
	if dim != nil && dim.Value() == DimThermoActuatorStatus.Value() && len(vals) > 0 {
		v, err := strconv.Atoi(vals[0])
		if err != nil {
			return 0, fmt.Errorf("could not parse actuator status: %v", err)
		}
		vs, ok := ThermoValveOrActuatorStatusFromValue(v)
		if ok {
			return vs, nil
		}
	}
	return 0, fmt.Errorf("could not parse actuator status from: %s", msg.GetFrameValue())
}

// ThermoGetLocalOffset extracts the Local Offset value.
func ThermoGetLocalOffset(msg *BaseOpenMessage) (*ThermoLocalOffset, error) {
	vals, err := msg.GetDimValues()
	if err != nil {
		return nil, err
	}
	if len(vals) > 0 {
		lo := ThermoLocalOffsetFromValue(vals[0])
		return lo, nil
	}
	return nil, NewFrameError("Could not parse local offset")
}
