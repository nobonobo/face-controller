package service

import (
	"encoding/binary"
	"machine/usb"

	"machine/usb/descriptor"
	"machine/usb/hid/joystick"
)

type SendReporter interface {
	SendReport(reportID byte, b []byte)
}

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

type JS struct {
	js       SendReporter
	buf      [13]byte
	axis     [4]int16
	triggers [2]uint8
	buttons  [10]bool
	hat      uint8
}

func (j *JS) Button(index int) bool {
	return j.buttons[index]
}

func (j *JS) SetButton(index int, push bool) {
	j.buttons[index] = push
}

func (j *JS) Hat(index int) joystick.HatDirection {
	return joystick.HatDirection(j.hat)
}

func (j *JS) SetHat(index int, dir joystick.HatDirection) {
	j.hat = uint8(dir)
}

func (j *JS) Axis(index int) int {
	return int(j.axis[index])
}

func (j *JS) SetAxis(index int, v int) {
	j.axis[index] = int16(v)
}

func (j *JS) SendState() {
	binary.LittleEndian.PutUint16(j.buf[0:2], uint16(j.axis[0]))
	binary.LittleEndian.PutUint16(j.buf[2:4], uint16(j.axis[1]))
	binary.LittleEndian.PutUint16(j.buf[4:6], uint16(j.axis[2]))
	binary.LittleEndian.PutUint16(j.buf[6:8], uint16(j.axis[3]))
	j.buf[8] = j.triggers[0]
	j.buf[9] = j.triggers[1]
	j.buf[10] = 0
	j.buf[11] = 0
	for i, v := range j.buttons[:8] {
		if v {
			j.buf[10] |= 1 << i
		}
	}
	for i, v := range j.buttons[8:] {
		if v {
			j.buf[11] |= 1 << i
		}
	}
	j.buf[12] = j.hat
	j.js.SendReport(1, j.buf[:])
}

func init() {
	usb.VendorID = 0x2786
	usb.ProductID = 0x000a
	usb.Product = "Gamepad Emulator"
	usb.Manufacturer = "Switch Science"

	js = &JS{
		js: joystick.UseSettings(joystick.Definitions{
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
		}, nil, nil, desc),
	}
}

var desc = descriptor.Append([][]byte{
	descriptor.HIDUsagePageGenericDesktop,
	descriptor.HIDUsageDesktopGamepad,

	descriptor.HIDCollectionApplication,

	descriptor.HIDReportID(1),
	descriptor.HIDUsageDesktopPointer,

	descriptor.HIDCollectionPhysical,

	descriptor.HIDUsageDesktopX,
	descriptor.HIDUsageDesktopY,
	descriptor.HIDUsageDesktopRx,
	descriptor.HIDUsageDesktopRy,
	descriptor.HIDLogicalMinimum(-8191),
	descriptor.HIDLogicalMaximum(8191),
	descriptor.HIDReportSize(16),
	descriptor.HIDReportCount(4),
	descriptor.HIDInputConstVarAbs,
	descriptor.HIDUsagePageGenericDesktop,
	descriptor.HIDUsageDesktopZ,
	descriptor.HIDUsageDesktopRz,
	descriptor.HIDLogicalMinimum(0),
	descriptor.HIDLogicalMaximum(255),
	descriptor.HIDReportSize(8),
	descriptor.HIDReportCount(2),
	descriptor.HIDInputConstVarAbs,
	descriptor.HIDUsagePageButton,
	descriptor.HIDUsageMinimum(1),
	descriptor.HIDUsageMaximum(10),
	descriptor.HIDLogicalMinimum(0),
	descriptor.HIDLogicalMaximum(1),
	descriptor.HIDReportSize(1),
	descriptor.HIDReportCount(10),
	descriptor.HIDInputConstVarAbs,
	[]byte{0x75, 0x06}, // Padding
	descriptor.HIDReportCount(1),
	descriptor.HIDInputConstVarAbs,
	descriptor.HIDUsagePageGenericDesktop,
	descriptor.HIDUsageDesktopHatSwitch,
	descriptor.HIDLogicalMinimum(0),
	descriptor.HIDLogicalMaximum(7),
	descriptor.HIDPhysicalMinimum(0),
	descriptor.HIDPhysicalMaximum(315),
	descriptor.HIDUnit(0x14), // UNIT (Eng Rotation: Centimeter)
	descriptor.HIDReportSize(4),
	descriptor.HIDReportCount(1),
	descriptor.HIDInputDataVarAbs,
	[]byte{0x75, 0x04}, // Padding
	descriptor.HIDReportCount(1),
	descriptor.HIDInputConstVarAbs,

	descriptor.HIDCollectionEnd,

	descriptor.HIDCollectionEnd,
})
