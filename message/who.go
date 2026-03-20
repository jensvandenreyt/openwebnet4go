package message

// Who represents OpenWebNet WHO types.
type Who int

const (
	WhoScenario                   Who = 0
	WhoLighting                   Who = 1
	WhoAutomation                 Who = 2
	WhoLoadControl                Who = 3 // Deprecated
	WhoThermoregulation           Who = 4
	WhoThermoregulationDiagnostic Who = 1004
	WhoBurglarAlarm               Who = 5
	WhoDoorEntrySystem            Who = 6
	WhoVideoDoorEntrySystem       Who = 7
	WhoAux                        Who = 9
	WhoGatewayManagement          Who = 13
	WhoLightShutterActuatorsLock  Who = 14
	WhoCENScenarioScheduler       Who = 15
	WhoCENPlusScenarioScheduler   Who = 25
	WhoSoundSystem1               Who = 16 // Deprecated
	WhoSoundSystem2               Who = 22
	WhoScenarioProgramming        Who = 17
	WhoEnergyManagement           Who = 18
	WhoEnergyManagementDiagnostic Who = 1018
	WhoLightingManagement         Who = 24
	WhoDiagnostic                 Who = 1000
	WhoAutomationDiagnostic       Who = 1001
	WhoDeviceDiagnostic           Who = 1013
	WhoUnknown                    Who = 9999
)

// whoValues maps int values to Who constants.
var whoValues = map[int]Who{
	0:    WhoScenario,
	1:    WhoLighting,
	2:    WhoAutomation,
	3:    WhoLoadControl,
	4:    WhoThermoregulation,
	1004: WhoThermoregulationDiagnostic,
	5:    WhoBurglarAlarm,
	6:    WhoDoorEntrySystem,
	7:    WhoVideoDoorEntrySystem,
	9:    WhoAux,
	13:   WhoGatewayManagement,
	14:   WhoLightShutterActuatorsLock,
	15:   WhoCENScenarioScheduler,
	25:   WhoCENPlusScenarioScheduler,
	16:   WhoSoundSystem1,
	22:   WhoSoundSystem2,
	17:   WhoScenarioProgramming,
	18:   WhoEnergyManagement,
	1018: WhoEnergyManagementDiagnostic,
	24:   WhoLightingManagement,
	1000: WhoDiagnostic,
	1001: WhoAutomationDiagnostic,
	1013: WhoDeviceDiagnostic,
	9999: WhoUnknown,
}

// Value returns the integer value of the Who.
func (w Who) Value() int {
	return int(w)
}

// IsValidWhoValue checks if the given integer is a valid Who value.
func IsValidWhoValue(value int) bool {
	_, ok := whoValues[value]
	return ok
}

// WhoFromValue returns the Who corresponding to the given integer value.
// Returns WhoUnknown and false if the value is not recognized.
func WhoFromValue(value int) (Who, bool) {
	w, ok := whoValues[value]
	return w, ok
}

// String returns the string name of the Who.
func (w Who) String() string {
	switch w {
	case WhoScenario:
		return "SCENARIO"
	case WhoLighting:
		return "LIGHTING"
	case WhoAutomation:
		return "AUTOMATION"
	case WhoLoadControl:
		return "LOAD_CONTROL"
	case WhoThermoregulation:
		return "THERMOREGULATION"
	case WhoThermoregulationDiagnostic:
		return "THERMOREGULATION_DIAGNOSTIC"
	case WhoBurglarAlarm:
		return "BURGLAR_ALARM"
	case WhoDoorEntrySystem:
		return "DOOR_ENTRY_SYSTEM"
	case WhoVideoDoorEntrySystem:
		return "VIDEO_DOOR_ENTRY_SYSTEM"
	case WhoAux:
		return "AUX"
	case WhoGatewayManagement:
		return "GATEWAY_MANAGEMENT"
	case WhoLightShutterActuatorsLock:
		return "LIGHT_SHUTTER_ACTUATORS_LOCK"
	case WhoCENScenarioScheduler:
		return "CEN_SCENARIO_SCHEDULER"
	case WhoCENPlusScenarioScheduler:
		return "CEN_PLUS_SCENARIO_SCHEDULER"
	case WhoSoundSystem1:
		return "SOUND_SYSTEM_1"
	case WhoSoundSystem2:
		return "SOUND_SYSTEM_2"
	case WhoScenarioProgramming:
		return "SCENARIO_PROGRAMMING"
	case WhoEnergyManagement:
		return "ENERGY_MANAGEMENT"
	case WhoEnergyManagementDiagnostic:
		return "ENERGY_MANAGEMENT_DIAGNOSTIC"
	case WhoLightingManagement:
		return "LIGHTING_MANAGEMENT"
	case WhoDiagnostic:
		return "DIAGNOSTIC"
	case WhoAutomationDiagnostic:
		return "AUTOMATION_DIAGNOSTIC"
	case WhoDeviceDiagnostic:
		return "DEVICE_DIAGNOSTIC"
	case WhoUnknown:
		return "UNKNOWN"
	default:
		return "UNKNOWN"
	}
}
