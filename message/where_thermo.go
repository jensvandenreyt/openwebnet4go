package message

import (
	"fmt"
	"strconv"
	"strings"
)

// WhereThermo is WHERE for Thermoregulation frames.
type WhereThermo struct {
	BaseWhere
	zone       int
	probe      int
	actuator   int
	standalone bool
}

var WhereThermoAllMasterProbes = &WhereThermo{
	BaseWhere:  BaseWhere{WhereStr: "0"},
	zone:       0,
	probe:      -1,
	actuator:   -1,
	standalone: true,
}

func NewWhereThermo(w string) (*WhereThermo, error) {
	bw, err := NewBaseWhere(w)
	if err != nil {
		return nil, err
	}
	wt := &WhereThermo{BaseWhere: *bw, probe: -1, actuator: -1}

	pos := strings.Index(w, "#")
	if pos >= 0 { // '#' is present
		if pos == 0 { // case '#x'
			wt.standalone = false
			z, err := strconv.Atoi(w[1:])
			if err != nil {
				return nil, fmt.Errorf("WHERE address '%s' is invalid: %v", w, err)
			}
			wt.zone = z
		} else { // case 'x#x'
			wt.standalone = true
			z, err := strconv.Atoi(w[:pos])
			if err != nil {
				return nil, fmt.Errorf("WHERE address '%s' is invalid: %v", w, err)
			}
			wt.zone = z
			a, err := strconv.Atoi(w[pos+1:])
			if err != nil {
				return nil, fmt.Errorf("WHERE address '%s' is invalid: %v", w, err)
			}
			wt.actuator = a
		}
	} else { // no '#' present
		wt.standalone = true
		z, err := strconv.Atoi(w)
		if err != nil {
			return nil, fmt.Errorf("WHERE address '%s' is invalid: %v", w, err)
		}
		if z > 99 { // case 'pZZ'
			p, err := strconv.Atoi(w[:1])
			if err != nil {
				return nil, fmt.Errorf("WHERE address '%s' is invalid: %v", w, err)
			}
			wt.probe = p
			zz, err := strconv.Atoi(w[1:])
			if err != nil {
				return nil, fmt.Errorf("WHERE address '%s' is invalid: %v", w, err)
			}
			wt.zone = zz
		} else if strings.HasPrefix(w, "0") && len(w) > 1 { // case '0ZZ'
			wt.probe = 0
			zz, err := strconv.Atoi(w[1:])
			if err != nil {
				return nil, fmt.Errorf("WHERE address '%s' is invalid: %v", w, err)
			}
			wt.zone = zz
		} else {
			wt.zone = z
		}
	}

	if wt.zone < 0 || wt.zone > 99 {
		return nil, fmt.Errorf("WHERE address '%s' is invalid: zone not in range [0-99]", w)
	}
	if wt.probe < -1 || wt.probe > 9 {
		return nil, fmt.Errorf("WHERE address '%s' is invalid: probe not in range [0-9]", w)
	}
	if wt.actuator < -1 || wt.actuator > 9 {
		return nil, fmt.Errorf("WHERE address '%s' is invalid: actuator not in range [0-9]", w)
	}

	return wt, nil
}

// GetZone returns the Zone for this WHERE.
func (wt *WhereThermo) GetZone() int { return wt.zone }

// GetProbe returns the probe for this WHERE, 0 for all probes, or -1 if no probe is present.
func (wt *WhereThermo) GetProbe() int { return wt.probe }

// GetActuator returns the actuator (1-9) for this WHERE, 0 for all actuators, or -1 if no actuator is present.
func (wt *WhereThermo) GetActuator() int { return wt.actuator }

// IsStandalone returns true if WHERE is a standalone configuration.
func (wt *WhereThermo) IsStandalone() bool { return wt.standalone }

// IsCentralUnit returns true if WHERE is Central Unit (where=#0).
func (wt *WhereThermo) IsCentralUnit() bool { return wt.zone == 0 && !wt.standalone }

// IsProbe returns true if WHERE is a probe address (where=pZZ).
func (wt *WhereThermo) IsProbe() bool { return wt.probe >= 0 }
