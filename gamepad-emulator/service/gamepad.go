package service

import (
	"machine/usb"

	"machine/usb/descriptor"
	"machine/usb/hid/joystick"
)

type JoySticker interface {
	Button(index int) bool
	SetButton(index int, push bool)
	Hat(index int) joystick.HatDirection
	SetHat(index int, dir joystick.HatDirection)
	Axis(index int) int
	SetAxis(index int, v int)
	SendState()
}

var js JoySticker

func init() {
	//usb.VendorID = 0x045E
	usb.ProductID = 0x000b
	usb.Product = "Gamepad Emulator"
	usb.Manufacturer = "Switch Science"
	js = joystick.UseSettings(joystick.Definitions{
		ReportID:     1,
		ButtonCnt:    10,
		HatSwitchCnt: 1,
		AxisDefs: []joystick.Constraint{
			{MinIn: -32767, MaxIn: 32767, MinOut: -32767, MaxOut: 32767},
			{MinIn: -32767, MaxIn: 32767, MinOut: -32767, MaxOut: 32767},
			{MinIn: -32767, MaxIn: 32767, MinOut: -32767, MaxOut: 32767},
			{MinIn: -32767, MaxIn: 32767, MinOut: -32767, MaxOut: 32767},
			{MinIn: 0, MaxIn: 1023, MinOut: 0, MaxOut: 1023},
			{MinIn: 0, MaxIn: 1023, MinOut: 0, MaxOut: 1023},
		},
	}, nil, nil, desc)
}

var desc = descriptor.Append([][]byte{
	descriptor.HIDUsagePageGenericDesktop,
	descriptor.HIDUsageDesktopGamepad,
	descriptor.HIDCollectionApplication,

	descriptor.HIDReportID(1),
	descriptor.HIDUsagePageButton,
	descriptor.HIDUsageMinimum(1),
	descriptor.HIDUsageMaximum(16),
	descriptor.HIDLogicalMinimum(0),
	descriptor.HIDLogicalMaximum(1),
	descriptor.HIDReportSize(1),
	descriptor.HIDReportCount(16),
	descriptor.HIDUnitExponent(0),
	descriptor.HIDUnit(0),
	descriptor.HIDInputDataVarAbs,

	descriptor.HIDUsagePageGenericDesktop,
	descriptor.HIDUsageDesktopHatSwitch,
	descriptor.HIDLogicalMinimum(0),
	descriptor.HIDLogicalMaximum(7),
	descriptor.HIDPhysicalMinimum(0),
	descriptor.HIDPhysicalMaximum(315),
	descriptor.HIDUnit(0x14),
	descriptor.HIDReportCount(1),
	descriptor.HIDReportSize(4),
	descriptor.HIDInputDataVarAbs,

	descriptor.HIDUsageDesktopHatSwitch,
	descriptor.HIDReportCount(1),
	descriptor.HIDReportSize(4),
	descriptor.HIDInputConstVarAbs,

	descriptor.HIDUsageDesktopPointer,
	descriptor.HIDLogicalMinimum(-32767),
	descriptor.HIDLogicalMaximum(32767),
	descriptor.HIDPhysicalMinimum(0),
	descriptor.HIDPhysicalMaximum(0),
	descriptor.HIDUnit(0),
	descriptor.HIDReportCount(2),
	descriptor.HIDReportSize(16),
	descriptor.HIDCollectionPhysical,
	descriptor.HIDUsageDesktopX,
	descriptor.HIDUsageDesktopY,
	descriptor.HIDInputConstVarAbs,
	descriptor.HIDCollectionEnd,

	descriptor.HIDUsageDesktopPointer,
	descriptor.HIDLogicalMinimum(-32767),
	descriptor.HIDLogicalMaximum(32767),
	descriptor.HIDPhysicalMinimum(0),
	descriptor.HIDPhysicalMaximum(0),
	descriptor.HIDUnit(0),
	descriptor.HIDReportCount(2),
	descriptor.HIDReportSize(16),
	descriptor.HIDCollectionPhysical,
	descriptor.HIDUsageDesktopRx,
	descriptor.HIDUsageDesktopRy,
	descriptor.HIDInputConstVarAbs,
	descriptor.HIDCollectionEnd,

	descriptor.HIDUsageDesktopPointer,
	descriptor.HIDLogicalMinimum(0),
	descriptor.HIDLogicalMaximum(1023),
	descriptor.HIDPhysicalMinimum(0),
	descriptor.HIDPhysicalMaximum(0),
	descriptor.HIDUnit(0),
	descriptor.HIDReportCount(2),
	descriptor.HIDReportSize(16),
	descriptor.HIDCollectionPhysical,
	descriptor.HIDUsageDesktopZ,
	descriptor.HIDUsageDesktopRz,
	descriptor.HIDInputDataVarAbs,
	descriptor.HIDCollectionEnd,

	descriptor.HIDCollectionEnd,
})
