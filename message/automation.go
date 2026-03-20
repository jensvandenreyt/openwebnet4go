package message

import (
	"fmt"
	"strings"
)

// WhatAutomation represents the WHAT values for Automation messages (WHO=2).
type WhatAutomation int

const (
	WhatAutomationStop WhatAutomation = 0
	WhatAutomationUp   WhatAutomation = 1
	WhatAutomationDown WhatAutomation = 2
)

func (w WhatAutomation) Value() int { return int(w) }

var whatAutomationValues = map[int]WhatAutomation{
	0: WhatAutomationStop,
	1: WhatAutomationUp,
	2: WhatAutomationDown,
}

func WhatAutomationFromValue(i int) What {
	v, ok := whatAutomationValues[i]
	if !ok {
		return nil
	}
	return v
}

// DimAutomation represents DIM values for Automation messages.
type DimAutomation int

const (
	DimAutomationShutterStatus DimAutomation = 10
	DimAutomationGotoLevel     DimAutomation = 11
)

func (d DimAutomation) Value() int { return int(d) }

var dimAutomationValues = map[int]DimAutomation{
	10: DimAutomationShutterStatus,
	11: DimAutomationGotoLevel,
}

func DimAutomationFromValue(i int) Dim {
	v, ok := dimAutomationValues[i]
	if !ok {
		return nil
	}
	return v
}

func newAutomationMessage(frame string) *BaseOpenMessage {
	msg := NewBaseOpenMessage(frame)
	msg.WhoField = WhoAutomation
	msg.WhatFromValue = WhatAutomationFromValue
	msg.DimFromValue = DimAutomationFromValue
	msg.ParseWhere = parseAutomationWhere
	msg.DetectDeviceTyp = detectAutomationDeviceType
	return msg
}

func parseAutomationWhere(whereStr string) (Where, error) {
	if whereStr == "" {
		return nil, NewFrameError("Frame has no WHERE part")
	}
	if strings.HasSuffix(whereStr, ZBNetwork) {
		return NewWhereZigBee(whereStr)
	}
	return NewWhereLightAutom(whereStr)
}

func detectAutomationDeviceType(msg *BaseOpenMessage) (OpenDeviceType, error) {
	if msg.IsCommand() {
		return DeviceSCSShutterControl, nil
	}
	return 0, nil
}

// --- Automation request constructors ---

// AutomationRequestStop creates a request to send STOP (*2*0*WHERE##).
func AutomationRequestStop(where string) *BaseOpenMessage {
	frame := fmt.Sprintf(FormatRequest, WhoAutomation.Value(), WhatAutomationStop.Value(), where)
	return newAutomationMessage(frame)
}

// AutomationRequestMoveUp creates a request to send UP (*2*1*WHERE##).
func AutomationRequestMoveUp(where string) *BaseOpenMessage {
	frame := fmt.Sprintf(FormatRequest, WhoAutomation.Value(), WhatAutomationUp.Value(), where)
	return newAutomationMessage(frame)
}

// AutomationRequestMoveDown creates a request to send DOWN (*2*2*WHERE##).
func AutomationRequestMoveDown(where string) *BaseOpenMessage {
	frame := fmt.Sprintf(FormatRequest, WhoAutomation.Value(), WhatAutomationDown.Value(), where)
	return newAutomationMessage(frame)
}

// AutomationRequestStatus creates a request for automation status (*#2*WHERE##).
func AutomationRequestStatus(where string) *BaseOpenMessage {
	frame := fmt.Sprintf(FormatStatus, WhoAutomation.Value(), where)
	return newAutomationMessage(frame)
}

// AutomationIsStop returns true if the automation message is STOP.
func AutomationIsStop(msg *BaseOpenMessage) bool {
	w := msg.GetWhat()
	if w == nil {
		return false
	}
	return w.Value() == WhatAutomationStop.Value()
}

// AutomationIsUp returns true if the automation message is UP.
func AutomationIsUp(msg *BaseOpenMessage) bool {
	w := msg.GetWhat()
	if w == nil {
		return false
	}
	return w.Value() == WhatAutomationUp.Value()
}

// AutomationIsDown returns true if the automation message is DOWN.
func AutomationIsDown(msg *BaseOpenMessage) bool {
	w := msg.GetWhat()
	if w == nil {
		return false
	}
	return w.Value() == WhatAutomationDown.Value()
}

// AutomationConvertUpDown converts an Automation message UP<->DOWN.
func AutomationConvertUpDown(msg *BaseOpenMessage) (*BaseOpenMessage, error) {
	if AutomationIsUp(msg) {
		newFrame := strings.Replace(msg.GetFrameValue(), "*2*1", "*2*2", 1)
		parsed, err := Parse(newFrame)
		if err != nil {
			return nil, err
		}
		if bm, ok := parsed.(*BaseOpenMessage); ok {
			return bm, nil
		}
		return nil, NewFrameError("unexpected message type after conversion")
	} else if AutomationIsDown(msg) {
		newFrame := strings.Replace(msg.GetFrameValue(), "*2*2", "*2*1", 1)
		parsed, err := Parse(newFrame)
		if err != nil {
			return nil, err
		}
		if bm, ok := parsed.(*BaseOpenMessage); ok {
			return bm, nil
		}
		return nil, NewFrameError("unexpected message type after conversion")
	}
	return msg, nil
}
