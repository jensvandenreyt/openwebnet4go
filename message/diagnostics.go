package message

import (
	"fmt"
	"strings"
)

// --- Energy Management Diagnostic (WHO=1018) ---

func newEnergyMgmtDiagnosticMessage(frame string) *BaseOpenMessage {
	msg := NewBaseOpenMessage(frame)
	msg.WhoField = WhoEnergyManagementDiagnostic
	msg.WhatFromValue = func(i int) What { return nil }
	msg.DimFromValue = func(i int) Dim { return nil }
	msg.ParseWhere = parseEnergyMgmtDiagnosticWhere
	msg.DetectDeviceTyp = detectEnergyMgmtDiagnosticDeviceType
	return msg
}

func parseEnergyMgmtDiagnosticWhere(whereStr string) (Where, error) {
	if whereStr == "" {
		return nil, NewFrameError("Frame has no WHERE part")
	}
	if strings.HasSuffix(whereStr, ZBNetwork) {
		return NewWhereZigBee(whereStr)
	}
	return NewWhereEnergyManagement(whereStr)
}

func detectEnergyMgmtDiagnosticDeviceType(msg *BaseOpenMessage) (OpenDeviceType, error) {
	w := msg.GetWhere()
	if w != nil && strings.HasPrefix(w.Value(), "5") {
		return DeviceSCSEnergyMeter, nil
	}
	return 0, nil
}

// EnergyMgmtDiagnosticRequestDiagnostic creates a diagnostic request (*#1018*WHERE*7##).
func EnergyMgmtDiagnosticRequestDiagnostic(where string) *BaseOpenMessage {
	frame := fmt.Sprintf(FormatDimensionRequest, WhoEnergyManagementDiagnostic.Value(), where, 7)
	return newEnergyMgmtDiagnosticMessage(frame)
}

// --- Thermoregulation Diagnostic (WHO=1004) ---

func newThermoregulationDiagnosticMessage(frame string) *BaseOpenMessage {
	msg := NewBaseOpenMessage(frame)
	msg.WhoField = WhoThermoregulationDiagnostic
	msg.WhatFromValue = func(i int) What { return nil }
	msg.DimFromValue = func(i int) Dim { return nil }
	msg.ParseWhere = parseThermoregulationDiagnosticWhere
	msg.DetectDeviceTyp = detectThermoregulationDiagnosticDeviceType
	return msg
}

func parseThermoregulationDiagnosticWhere(whereStr string) (Where, error) {
	if whereStr == "" {
		return nil, NewFrameError("Frame has no WHERE part")
	}
	if strings.HasSuffix(whereStr, ZBNetwork) {
		return NewWhereZigBee(whereStr)
	}
	return NewWhereThermo(whereStr)
}

func detectThermoregulationDiagnosticDeviceType(msg *BaseOpenMessage) (OpenDeviceType, error) {
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

// ThermoregulationDiagnosticRequestDiagnostic creates a diagnostic request (*#1004*WHERE*7##).
func ThermoregulationDiagnosticRequestDiagnostic(where string) *BaseOpenMessage {
	frame := fmt.Sprintf(FormatDimensionRequest, WhoThermoregulationDiagnostic.Value(), where, 7)
	return newThermoregulationDiagnosticMessage(frame)
}
