package message

// OpenDeviceType represents device types according to OpenWebNet specs.
type OpenDeviceType int

const (
	DeviceUnknown                      OpenDeviceType = 0
	DeviceScenarioControl              OpenDeviceType = 2
	DeviceZigBeeOnOffSwitch            OpenDeviceType = 256
	DeviceZigBeeDimmerControl          OpenDeviceType = 257
	DeviceZigBeeDimmerSwitch           OpenDeviceType = 258
	DeviceZigBeeSwitchMotionDetector   OpenDeviceType = 259
	DeviceZigBeeDaylightSensor         OpenDeviceType = 260
	DeviceSCSOnOffSwitch               OpenDeviceType = 261
	DeviceSCSDimmerControl             OpenDeviceType = 262
	DeviceSCSDimmerSwitch              OpenDeviceType = 263
	DeviceZigBeeWaterproof1GangSwitch  OpenDeviceType = 264
	DeviceZigBeeAutomaticDimmerSwitch  OpenDeviceType = 265
	DeviceZigBeeToggleControl          OpenDeviceType = 266
	DeviceSCSToggleControl             OpenDeviceType = 267
	DeviceZigBeeMotionDetector         OpenDeviceType = 268
	DeviceZigBeeSwitchMotionDetectorII OpenDeviceType = 269
	DeviceZigBeeMotionDetectorII       OpenDeviceType = 270
	DeviceMultifunctionScenarioControl OpenDeviceType = 273
	DeviceZigBeeOnOffControl           OpenDeviceType = 274
	DeviceZigBeeAuxiliaryMotionControl OpenDeviceType = 271
	DeviceSCSAuxiliaryToggleControl    OpenDeviceType = 272
	DeviceZigBeeAuxOnOff1GangSwitch    OpenDeviceType = 275
	DeviceZigBeeShutterControl         OpenDeviceType = 512
	DeviceZigBeeShutterSwitch          OpenDeviceType = 513
	DeviceSCSShutterControl            OpenDeviceType = 514
	DeviceSCSShutterSwitch             OpenDeviceType = 515
	DeviceSCSThermoSensor              OpenDeviceType = 410
	DeviceSCSThermoZone                OpenDeviceType = 420
	DeviceSCSThermoCentralUnit         OpenDeviceType = 430
	DeviceSCS1System14Gateway          OpenDeviceType = 1024
	DeviceSCS2System14Gateway          OpenDeviceType = 1025
	DeviceNetworkRepeater              OpenDeviceType = 1029
	DeviceOpenWebNetInterface          OpenDeviceType = 1030
	DeviceVideoSwitcher                OpenDeviceType = 1536
	DeviceSCSEnergyMeter               OpenDeviceType = 1830
	DeviceSCSDryContactIR              OpenDeviceType = 2510
)

var deviceTypeValues = map[int]OpenDeviceType{
	0: DeviceUnknown, 2: DeviceScenarioControl,
	256: DeviceZigBeeOnOffSwitch, 257: DeviceZigBeeDimmerControl,
	258: DeviceZigBeeDimmerSwitch, 259: DeviceZigBeeSwitchMotionDetector,
	260: DeviceZigBeeDaylightSensor, 261: DeviceSCSOnOffSwitch,
	262: DeviceSCSDimmerControl, 263: DeviceSCSDimmerSwitch,
	264: DeviceZigBeeWaterproof1GangSwitch, 265: DeviceZigBeeAutomaticDimmerSwitch,
	266: DeviceZigBeeToggleControl, 267: DeviceSCSToggleControl,
	268: DeviceZigBeeMotionDetector, 269: DeviceZigBeeSwitchMotionDetectorII,
	270: DeviceZigBeeMotionDetectorII, 273: DeviceMultifunctionScenarioControl,
	274: DeviceZigBeeOnOffControl, 271: DeviceZigBeeAuxiliaryMotionControl,
	272: DeviceSCSAuxiliaryToggleControl, 275: DeviceZigBeeAuxOnOff1GangSwitch,
	512: DeviceZigBeeShutterControl, 513: DeviceZigBeeShutterSwitch,
	514: DeviceSCSShutterControl, 515: DeviceSCSShutterSwitch,
	410: DeviceSCSThermoSensor, 420: DeviceSCSThermoZone,
	430: DeviceSCSThermoCentralUnit, 1024: DeviceSCS1System14Gateway,
	1025: DeviceSCS2System14Gateway, 1029: DeviceNetworkRepeater,
	1030: DeviceOpenWebNetInterface, 1536: DeviceVideoSwitcher,
	1830: DeviceSCSEnergyMeter, 2510: DeviceSCSDryContactIR,
}

// DeviceTypeFromValue returns the OpenDeviceType for the given integer.
func DeviceTypeFromValue(value int) (OpenDeviceType, bool) {
	dt, ok := deviceTypeValues[value]
	return dt, ok
}
