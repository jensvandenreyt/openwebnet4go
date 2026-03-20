package message

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// MaxFrameLength is the maximum OWN frame length.
const MaxFrameLength = 1024

// Format strings for building OpenWebNet frames.
const (
	FormatDimensionRequest     = "*#%d*%s*%d##"
	FormatDimensionWriting1V   = "*#%d*%s*#%d*%s##"
	FormatDimensionWriting2V   = "*#%d*%s*#%d*%s*%s##"
	FormatDimensionWriting1P1V = "*#%d*%s*#%d#%s*%s##"
	FormatRequest              = "*%d*%d*%s##"
	FormatRequestWhatStr       = "*%d*%s*%s##"
	FormatRequestParamStr      = "*%d*%s#%d*%s##"
	FormatStatus               = "*#%d*%s##"
)

// WhatFromValueFunc is a function type that returns a What from an int value.
type WhatFromValueFunc func(int) What

// DimFromValueFunc is a function type that returns a Dim from an int value.
type DimFromValueFunc func(int) Dim

// ParseWhereFunc is a function type that parses a WHERE string and returns a Where.
type ParseWhereFunc func(whereStr string) (Where, error)

// DetectDeviceTypeFunc is a function type that detects device type from a message.
type DetectDeviceTypeFunc func(msg *BaseOpenMessage) (OpenDeviceType, error)

// BaseOpenMessage is the base type for all OpenWebNet message types.
type BaseOpenMessage struct {
	BaseMessage

	WhatStr  string
	WhereStr string
	DimStr   string

	WhoField   Who
	WhatField  What
	WhereField Where
	DimField   Dim

	IsCmd               *bool
	IsCommandTransField *bool

	IsDimWritingField *bool
	DimParams         []int
	DimValues         []string

	CommandParams []int

	// Function hooks set by concrete message types
	WhatFromValue   WhatFromValueFunc
	DimFromValue    DimFromValueFunc
	ParseWhere      ParseWhereFunc
	DetectDeviceTyp DetectDeviceTypeFunc
}

// NewBaseOpenMessage creates a new BaseOpenMessage from a frame string.
func NewBaseOpenMessage(frame string) *BaseOpenMessage {
	return &BaseOpenMessage{
		BaseMessage: BaseMessage{FrameValue: frame},
	}
}

// IsCommand returns true if this is a command frame.
func (m *BaseOpenMessage) IsCommand() bool {
	if m.IsCmd == nil {
		isCmd := !strings.HasPrefix(m.FrameValue, FrameStartDim)
		m.IsCmd = &isCmd
	}
	return *m.IsCmd
}

// ToStringVerbose returns a verbose representation.
func (m *BaseOpenMessage) ToStringVerbose() string {
	verbose := fmt.Sprintf("<%s>{%s", m.FrameValue, m.WhoField.String())
	if m.IsCommand() {
		what := m.GetWhat()
		if what != nil {
			verbose += fmt.Sprintf("-%v", what)
		}
		isCT, _ := m.IsCommandTranslation()
		if isCT {
			verbose += "(translation)"
		}
		cp, _ := m.GetCommandParams()
		if len(cp) > 0 {
			verbose += fmt.Sprintf(",cmdParams=%v", cp)
		}
		where := m.GetWhere()
		if where != nil {
			verbose += fmt.Sprintf(",%s", where)
		}
	} else {
		dim := m.GetDim()
		if dim != nil {
			verbose += fmt.Sprintf("-%v", dim)
		}
		if m.IsDimWriting() {
			verbose += " (writing)"
		}
		where := m.GetWhere()
		if where != nil {
			verbose += fmt.Sprintf(",%s", where)
		}
		dp, _ := m.GetDimParams()
		if len(dp) > 0 {
			verbose += fmt.Sprintf(",dimParams=%v", dp)
		}
		dv, _ := m.GetDimValues()
		if len(dv) > 0 {
			verbose += fmt.Sprintf(",dimValues=%v", dv)
		}
	}
	return verbose + "}"
}

// GetWho returns the message WHO.
func (m *BaseOpenMessage) GetWho() Who {
	return m.WhoField
}

// GetWhat returns the message WHAT, or nil if message has no valid WHAT part.
func (m *BaseOpenMessage) GetWhat() What {
	if m.WhatField == nil {
		if m.WhatStr == "" {
			m.parseParts()
		}
		m.parseWhat()
	}
	return m.WhatField
}

// GetWhere returns the message WHERE, or nil if message has no valid WHERE part.
func (m *BaseOpenMessage) GetWhere() Where {
	if m.WhereField == nil {
		if m.WhereStr == "" {
			m.parseParts()
		}
		if m.ParseWhere != nil && m.WhereStr != "" {
			where, err := m.ParseWhere(m.WhereStr)
			if err == nil {
				m.WhereField = where
			}
		}
	}
	return m.WhereField
}

// GetDim returns the message DIM (dimension), or nil if no DIM is present.
func (m *BaseOpenMessage) GetDim() Dim {
	if m.DimField == nil {
		if m.DimStr == "" {
			m.parseParts()
		}
		m.parseDim()
	}
	return m.DimField
}

// IsDimWriting returns true if this is a dimension writing message.
func (m *BaseOpenMessage) IsDimWriting() bool {
	if m.IsDimWritingField == nil {
		m.GetDim()
	}
	if m.IsDimWritingField == nil {
		return false
	}
	return *m.IsDimWritingField
}

// IsCommandTranslation checks if message is a command translation (*WHO*1000#WHAT*...##).
func (m *BaseOpenMessage) IsCommandTranslation() (bool, error) {
	if m.IsCommand() {
		if m.IsCommandTransField == nil {
			m.GetWhat()
		}
	} else {
		f := false
		m.IsCommandTransField = &f
	}
	if m.IsCommandTransField != nil {
		return *m.IsCommandTransField, nil
	}
	return false, nil
}

// GetCommandParams returns message command parameters (*WHO*WHAT#Param1#Param2...#ParamN*...).
func (m *BaseOpenMessage) GetCommandParams() ([]int, error) {
	if m.CommandParams == nil {
		m.GetWhat()
	}
	return m.CommandParams, nil
}

// GetDimParams returns DIM parameters PAR1..PARN.
func (m *BaseOpenMessage) GetDimParams() ([]int, error) {
	if m.DimParams == nil {
		m.GetDim()
	}
	return m.DimParams, nil
}

// GetDimValues returns DIM values.
func (m *BaseOpenMessage) GetDimValues() ([]string, error) {
	if m.DimValues == nil {
		m.GetDim()
	}
	return m.DimValues, nil
}

func (m *BaseOpenMessage) parseParts() {
	parts, err := getPartsStrings(m.FrameValue)
	if err != nil {
		return
	}
	m.parsePartsFromSlice(parts)
}

func (m *BaseOpenMessage) parsePartsFromSlice(parts []string) {
	if m.IsCommand() {
		if len(parts) > 2 {
			m.WhatStr = parts[2] // second part is WHAT
		}
		if len(parts) > 3 {
			m.WhereStr = parts[3] // third part is WHERE (optional)
		}
	} else {
		if len(parts) > 2 && parts[2] != "" {
			m.WhereStr = parts[2] // second part is WHERE
		}
		if len(parts) >= 4 {
			m.DimStr = parts[3] // third part is DIM
			// copy last parts as DIM values
			if len(parts) > 4 {
				m.DimValues = make([]string, len(parts)-4)
				copy(m.DimValues, parts[4:])
			} else {
				m.DimValues = []string{}
			}
		}
	}
}

func (m *BaseOpenMessage) parseWhat() {
	if m.WhatStr == "" || m.WhatFromValue == nil {
		return
	}
	parts := strings.Split(m.WhatStr, "#")
	if len(parts) == 0 {
		return
	}
	partsIndex := 0
	firstVal, err := strconv.Atoi(parts[partsIndex])
	if err != nil {
		return
	}
	if firstVal == WhatCommandTranslation && len(parts) > 1 {
		t := true
		m.IsCommandTransField = &t
		partsIndex++
	} else {
		f := false
		m.IsCommandTransField = &f
	}
	whatVal, err := strconv.Atoi(parts[partsIndex])
	if err != nil {
		return
	}
	m.WhatField = m.WhatFromValue(whatVal)
	m.CommandParams = []int{}
	if m.WhatField == nil {
		return
	}
	if len(parts) > 1 {
		m.CommandParams = make([]int, len(parts)-partsIndex-1)
		for i := 0; i < len(m.CommandParams); i++ {
			v, err := strconv.Atoi(parts[i+partsIndex+1])
			if err != nil {
				return
			}
			m.CommandParams[i] = v
		}
	}
}

func (m *BaseOpenMessage) parseDim() {
	if m.DimStr == "" || m.DimFromValue == nil {
		return
	}
	ds := m.DimStr
	if strings.HasPrefix(ds, "#") { // Dim writing
		t := true
		m.IsDimWritingField = &t
		ds = ds[1:]
	} else {
		f := false
		m.IsDimWritingField = &f
	}
	dimParts := strings.Split(ds, "#")
	dimVal, err := strconv.Atoi(dimParts[0])
	if err != nil {
		return
	}
	m.DimField = m.DimFromValue(dimVal)
	if m.DimField == nil {
		return
	}
	// copy last parts as dim params
	if len(dimParts) > 1 {
		m.DimParams = make([]int, len(dimParts)-1)
		for i := 1; i < len(dimParts); i++ {
			v, err := strconv.Atoi(dimParts[i])
			if err != nil {
				return
			}
			m.DimParams[i-1] = v
		}
	} else {
		m.DimParams = []int{}
	}
}

// AddValues adds values separated by '*' at the end of a frame string.
func AddValues(msgStr string, vals ...string) string {
	str := msgStr[:len(msgStr)-2]
	for _, v := range vals {
		str = str + "*" + v
	}
	str += FrameEnd
	return str
}

func getPartsStrings(frame string) ([]string, error) {
	// remove trailing "##" and get frame parts separated by '*'
	trimmed := frame[:len(frame)-2]
	parts := strings.Split(trimmed, "*")
	if len(parts) == 0 {
		return nil, NewMalformedFrameError("Invalid frame")
	}
	if len(parts) < 3 {
		return nil, NewMalformedFrameError("Cmd/Dim frames must have at least 2 non-empty sections separated by '*'")
	}
	return parts, nil
}

// Parse parses a frame string and returns a new OpenMessage.
func Parse(frame string) (OpenMessage, error) {
	if frame == "" {
		return nil, NewMalformedFrameError("Frame is nil")
	}
	if frame == FrameACK || frame == FrameNACK || frame == FrameBusyNACK {
		return NewAckOpenMessage(frame), nil
	}
	if !strings.HasSuffix(frame, FrameEnd) {
		return nil, NewMalformedFrameError("Frame does not end with terminator " + FrameEnd)
	}
	isCmd := true
	if strings.HasPrefix(frame, FrameStartDim) {
		isCmd = false
	} else if !strings.HasPrefix(frame, FrameStart) {
		return nil, NewMalformedFrameError("Frame does not start with '*' or '*#'")
	}
	if len(frame) > MaxFrameLength {
		return nil, NewMalformedFrameError(fmt.Sprintf("Frame length is > %d", MaxFrameLength))
	}
	// check for bad characters
	for _, c := range frame {
		if !unicode.IsDigit(c) && c != '#' && c != '*' {
			return nil, NewMalformedFrameError("Frame can only contain '#', '*' or digits [0-9]")
		}
	}
	parts, err := getPartsStrings(frame)
	if err != nil {
		return nil, err
	}
	// parts[0] is empty, first is WHO
	whoStr := parts[1]
	if !isCmd {
		whoStr = parts[1][1:] // remove '#' from WHO part
	}
	msg, err := parseWho(whoStr, frame)
	if err != nil {
		return nil, err
	}
	msg.IsCmd = &isCmd
	msg.parsePartsFromSlice(parts)
	return msg, nil
}

func parseWho(whoPart string, frame string) (*BaseOpenMessage, error) {
	whoInt, err := strconv.Atoi(whoPart)
	if err != nil {
		return nil, NewMalformedFrameError("WHO not recognized: " + whoPart)
	}
	who, ok := WhoFromValue(whoInt)
	if !ok {
		return nil, NewMalformedFrameError("WHO not recognized: " + whoPart)
	}

	switch who {
	case WhoGatewayManagement:
		return newGatewayMgmtMessage(frame), nil
	case WhoLighting:
		return newLightingMessage(frame), nil
	case WhoAutomation:
		return newAutomationMessage(frame), nil
	case WhoEnergyManagement:
		return newEnergyManagementMessage(frame), nil
	case WhoThermoregulation:
		return newThermoregulationMessage(frame), nil
	case WhoCENScenarioScheduler:
		return newCENScenarioMessage(frame), nil
	case WhoCENPlusScenarioScheduler:
		return newCENPlusScenarioMessage(frame), nil
	case WhoEnergyManagementDiagnostic:
		return newEnergyMgmtDiagnosticMessage(frame), nil
	case WhoThermoregulationDiagnostic:
		return newThermoregulationDiagnosticMessage(frame), nil
	default:
		return nil, NewUnsupportedFrameError("WHO not recognized/supported: " + who.String())
	}
}
