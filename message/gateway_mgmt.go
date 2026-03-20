package message

import (
	"fmt"
	"strconv"
	"strings"
)

// WhatGatewayMgmt represents WHAT values for Gateway Management messages (WHO=13).
type WhatGatewayMgmt int

const (
	WhatGatewayMgmtBootMode      WhatGatewayMgmt = 12
	WhatGatewayMgmtResetDevice   WhatGatewayMgmt = 22
	WhatGatewayMgmtCreateNetwork WhatGatewayMgmt = 30
	WhatGatewayMgmtCloseNetwork  WhatGatewayMgmt = 31
	WhatGatewayMgmtOpenNetwork   WhatGatewayMgmt = 32
	WhatGatewayMgmtJoinNetwork   WhatGatewayMgmt = 33
	WhatGatewayMgmtLeaveNetwork  WhatGatewayMgmt = 34
	WhatGatewayMgmtKeepConnect   WhatGatewayMgmt = 60
	WhatGatewayMgmtScan          WhatGatewayMgmt = 65
	WhatGatewayMgmtSupervisor    WhatGatewayMgmt = 66
	WhatGatewayMgmtTest          WhatGatewayMgmt = 9999
)

func (w WhatGatewayMgmt) Value() int { return int(w) }

var whatGatewayMgmtValues = map[int]WhatGatewayMgmt{
	12: WhatGatewayMgmtBootMode, 22: WhatGatewayMgmtResetDevice,
	30: WhatGatewayMgmtCreateNetwork, 31: WhatGatewayMgmtCloseNetwork,
	32: WhatGatewayMgmtOpenNetwork, 33: WhatGatewayMgmtJoinNetwork,
	34: WhatGatewayMgmtLeaveNetwork, 60: WhatGatewayMgmtKeepConnect,
	65: WhatGatewayMgmtScan, 66: WhatGatewayMgmtSupervisor,
	9999: WhatGatewayMgmtTest,
}

func WhatGatewayMgmtFromValue(i int) What {
	v, ok := whatGatewayMgmtValues[i]
	if !ok {
		return nil
	}
	return v
}

// DimGatewayMgmt represents DIM values for Gateway Management messages.
type DimGatewayMgmt int

const (
	DimGatewayMgmtMACAddress      DimGatewayMgmt = 12
	DimGatewayMgmtModel           DimGatewayMgmt = 15
	DimGatewayMgmtFirmwareVersion DimGatewayMgmt = 16
	DimGatewayMgmtHardwareVersion DimGatewayMgmt = 17
	DimGatewayMgmtWhoImplemented  DimGatewayMgmt = 26
	DimGatewayMgmtProductInfo     DimGatewayMgmt = 66
	DimGatewayMgmtNbNetwProd      DimGatewayMgmt = 67
	DimGatewayMgmtIdentify        DimGatewayMgmt = 70
	DimGatewayMgmtZigBeeChannel   DimGatewayMgmt = 71
)

func (d DimGatewayMgmt) Value() int { return int(d) }

var dimGatewayMgmtValues = map[int]DimGatewayMgmt{
	12: DimGatewayMgmtMACAddress, 15: DimGatewayMgmtModel,
	16: DimGatewayMgmtFirmwareVersion, 17: DimGatewayMgmtHardwareVersion,
	26: DimGatewayMgmtWhoImplemented, 66: DimGatewayMgmtProductInfo,
	67: DimGatewayMgmtNbNetwProd, 70: DimGatewayMgmtIdentify,
	71: DimGatewayMgmtZigBeeChannel,
}

func DimGatewayMgmtFromValue(i int) Dim {
	v, ok := dimGatewayMgmtValues[i]
	if !ok {
		return nil
	}
	return v
}

func newGatewayMgmtMessage(frame string) *BaseOpenMessage {
	msg := NewBaseOpenMessage(frame)
	msg.WhoField = WhoGatewayManagement
	msg.WhatFromValue = WhatGatewayMgmtFromValue
	msg.DimFromValue = DimGatewayMgmtFromValue
	msg.ParseWhere = parseGatewayMgmtWhere
	msg.DetectDeviceTyp = func(m *BaseOpenMessage) (OpenDeviceType, error) { return 0, nil }
	return msg
}

func parseGatewayMgmtWhere(whereStr string) (Where, error) {
	if whereStr == "" {
		return nil, nil // GatewayMgmt WHERE can be empty
	}
	return NewWhereZigBee(whereStr)
}

// --- Gateway Management request constructors ---

// GatewayMgmtRequestSupervisor creates a request for supervisor mode (*13*66*##).
func GatewayMgmtRequestSupervisor() *BaseOpenMessage {
	frame := fmt.Sprintf(FormatRequest, WhoGatewayManagement.Value(), WhatGatewayMgmtSupervisor.Value(), "")
	return newGatewayMgmtMessage(frame)
}

// GatewayMgmtRequestKeepConnect creates a request keep connect (*13*60*##).
func GatewayMgmtRequestKeepConnect() *BaseOpenMessage {
	frame := fmt.Sprintf(FormatRequest, WhoGatewayManagement.Value(), WhatGatewayMgmtKeepConnect.Value(), "")
	return newGatewayMgmtMessage(frame)
}

// GatewayMgmtRequestMACAddress creates a request for gateway MAC address (*#13**12##).
func GatewayMgmtRequestMACAddress() *BaseOpenMessage {
	frame := fmt.Sprintf(FormatDimensionRequest, WhoGatewayManagement.Value(), "", DimGatewayMgmtMACAddress.Value())
	return newGatewayMgmtMessage(frame)
}

// GatewayMgmtParseMACAddress parses MAC address from a GatewayMgmt message.
func GatewayMgmtParseMACAddress(msg *BaseOpenMessage) ([]byte, error) {
	vals, err := msg.GetDimValues()
	if err != nil {
		return nil, err
	}
	mac := make([]byte, len(vals))
	for i, v := range vals {
		n, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("cannot parse MAC address byte: %v", err)
		}
		mac[i] = byte(n)
	}
	return mac, nil
}

// GatewayMgmtRequestModel creates a request for gateway model (*#13**15##).
func GatewayMgmtRequestModel() *BaseOpenMessage {
	frame := fmt.Sprintf(FormatDimensionRequest, WhoGatewayManagement.Value(), "", DimGatewayMgmtModel.Value())
	return newGatewayMgmtMessage(frame)
}

// GatewayMgmtRequestFirmwareVersion creates a request for firmware version (*#13**16##).
func GatewayMgmtRequestFirmwareVersion() *BaseOpenMessage {
	frame := fmt.Sprintf(FormatDimensionRequest, WhoGatewayManagement.Value(), "", DimGatewayMgmtFirmwareVersion.Value())
	return newGatewayMgmtMessage(frame)
}

// GatewayMgmtParseFirmwareVersion parses firmware version from a GatewayMgmt message.
func GatewayMgmtParseFirmwareVersion(msg *BaseOpenMessage) (string, error) {
	vals, err := msg.GetDimValues()
	if err != nil {
		return "", err
	}
	if len(vals) < 3 {
		return "", NewFrameError("Not enough values for firmware version")
	}
	return vals[0] + "." + vals[1] + "." + vals[2], nil
}

// GatewayMgmtRequestScanNetwork creates a request to scan network (*13*65*##).
func GatewayMgmtRequestScanNetwork() *BaseOpenMessage {
	frame := fmt.Sprintf(FormatRequest, WhoGatewayManagement.Value(), WhatGatewayMgmtScan.Value(), "")
	return newGatewayMgmtMessage(frame)
}

// GatewayMgmtRequestProductInfo creates a request for product information.
// NOTE: Uses '*' instead of '#' to separate index due to USB gateway bug.
func GatewayMgmtRequestProductInfo(index int) *BaseOpenMessage {
	frame := fmt.Sprintf(FormatDimensionRequest, WhoGatewayManagement.Value(), "", DimGatewayMgmtProductInfo.Value())
	frame = AddValues(frame, strconv.Itoa(index))
	return newGatewayMgmtMessage(frame)
}

// IsGatewayMgmtMessage checks if a BaseOpenMessage is a gateway management message.
func IsGatewayMgmtMessage(msg *BaseOpenMessage) bool {
	return msg.WhoField == WhoGatewayManagement
}

// GatewayMgmtFormatMACAddress formats a MAC address byte slice as a colon-separated hex string.
func GatewayMgmtFormatMACAddress(mac []byte) string {
	if mac == nil {
		return ""
	}
	parts := make([]string, len(mac))
	for i, b := range mac {
		parts[i] = fmt.Sprintf("%02x", b)
	}
	return strings.Join(parts, ":")
}
