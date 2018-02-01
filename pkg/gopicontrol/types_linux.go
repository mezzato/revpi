// Created by cgo -godefs - DO NOT EDIT
// cgo -godefs=true -- -I/home/pi/go/src/github.com/mezzato/revpi/pkg/gopicontrol/../../kunbus/interface/piControl ctypes_linux.go

package gopicontrol

type SDeviceInfo struct {
	I8uAddress		uint8
	Pad_cgo_0		[3]byte
	I32uSerialnumber	uint32
	I16uModuleType		uint16
	I16uHW_Revision		uint16
	I16uSW_Major		uint16
	I16uSW_Minor		uint16
	I32uSVN_Revision	uint32
	I16uInputLength		uint16
	I16uOutputLength	uint16
	I16uConfigLength	uint16
	I16uBaseOffset		uint16
	I16uInputOffset		uint16
	I16uOutputOffset	uint16
	I16uConfigOffset	uint16
	I16uFirstEntry		uint16
	I16uEntries		uint16
	I8uModuleState		uint8
	I8uActive		uint8
	I8uReserve		[30]uint8
	Pad_cgo_1		[2]byte
}

type SEntryInfo struct {
	I8uAddress	uint8
	I8uType		uint8
	I16uIndex	uint16
	I16uBitLength	uint16
	I8uBitPos	uint8
	Pad_cgo_0	[1]byte
	I16uOffset	uint16
	Pad_cgo_1	[2]byte
	I32uDefault	uint32
	StrVarName	[32]uint8
}

type SPIValue struct {
	I16uAddress	uint16
	I8uBit		uint8
	I8uValue	uint8
}

type SPIVariable struct {
	StrVarName	[32]uint8
	I16uAddress	uint16
	I8uBit		uint8
	Pad_cgo_0	[1]byte
	I16uLength	uint16
}

type SDIOResetCounter struct {
	I8uAddress	uint8
	Pad_cgo_0	[1]byte
	I16uBitfield	uint16
}

const (
	PICONTROL_DEVICE		= "/dev/piControl0"
	KB_RESET			= 0x4b0c
	KB_GET_DEVICE_INFO		= 0x4b0e
	KB_GET_DEVICE_INFO_LIST		= 0x4b0d
	KB_GET_VALUE			= 0x4b0f
	KB_SET_VALUE			= 0x4b10
	KB_FIND_VARIABLE		= 0x4b11
	KB_DIO_RESET_COUNTER		= 0x4b14
	KB_UPDATE_DEVICE_FIRMWARE	= 0x4b13
	KB_GET_LAST_MESSAGE		= 0x4b15
	KB_INTERN_IO_MSG		= 0x4b65
	KB_WAIT_FOR_EVENT		= 0x4b32
	PICONTROL_NOT_CONNECTED		= 0x8000
	PICONTROL_NOT_CONNECTED_MASK	= 0x7fff
	PICONTROL_SW_MODBUS_TCP_SLAVE	= 0x6001
	PICONTROL_SW_MODBUS_RTU_SLAVE	= 0x6002
	PICONTROL_SW_MODBUS_TCP_MASTER	= 0x6003
	PICONTROL_SW_MODBUS_RTU_MASTER	= 0x6004
)
