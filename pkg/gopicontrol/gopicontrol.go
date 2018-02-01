package gopicontrol

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"syscall"

	"unsafe"

	"golang.org/x/sys/unix"
)

// ByteToUint8Array converts a byte slice to a uint8 array.
func ByteToUint8Array(s []byte) (r [32]uint8) {
	for i, c := range s {
		r[i] = uint8(c)
	}
	return r
}

// NumToBytes converts a generic fixed-size value to its byte representation.
func NumToBytes(num interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, num)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
		return nil, err
	}
	return buf.Bytes(), nil
}

// see /home/enrico/go_src/src/golang.org/x/sys/unix/zsyscall_linux_arm64.go
// https://github.com/golang/crypto/blob/master/ssh/terminal/util.go

// Do the interface allocations only once for common
// Errno values.
var (
	errEAGAIN error = syscall.EAGAIN
	errEINVAL error = syscall.EINVAL
	errENOENT error = syscall.ENOENT
)

// errnoErr returns common boxed Errno values, to prevent
// allocations at runtime.
func errnoErr(e syscall.Errno) error {
	switch e {
	case 0:
		return nil
	case unix.EAGAIN:
		return errEAGAIN
	case unix.EINVAL:
		return errEINVAL
	case unix.ENOENT:
		return errENOENT
	}
	return e
}

// ioctl invokes a Unix syscall.
func ioctl(fd uintptr, req uint, arg uintptr) (r1 uintptr, r2 uintptr, err error) {
	r1, r2, e1 := unix.Syscall(unix.SYS_IOCTL, fd, uintptr(req), uintptr(arg))
	if e1 != 0 {
		err = errnoErr(e1)
		return r1, r2, err
	}
	return r1, r2, nil
}

// RevPiControl is an object representing an open file handle to the piControl driver file descriptor.
type RevPiControl struct {
	handle *os.File
}

// NewRevPiControl creates a new RevPiControl object.
func NewRevPiControl() *RevPiControl {
	return &RevPiControl{}
}

// Open opens the file handle.
// see also: golang.org/x/sys/unix/syscall_unix_test.go
func (c *RevPiControl) Open() (err error) {
	/* open handle if needed */
	if c.handle != nil {
		return nil
	}

	c.handle, err = os.OpenFile(PICONTROL_DEVICE, os.O_RDWR, 0)
	if err != nil {
		return err
	}
	return nil
}

// Close closes the file handle.
func (c *RevPiControl) Close() (err error) {
	/* open handle if needed */
	if c.handle != nil {
		if err = c.handle.Close(); err != nil {
			return err
		}
		c.handle = nil
	}
	return nil
}

// Reset initializes the Pi Control Interface.
func (c *RevPiControl) Reset() (err error) {
	if err = c.Open(); err != nil {
		return err
	}

	if _, _, err = ioctl(c.handle.Fd(), KB_RESET, uintptr(0)); err != nil {
		return err
	}
	return nil
}

// Read gets process data from a specific position, reads len(pData) bytes from file.
// Returns number of bytes read or error.
func (c *RevPiControl) Read(offset uint32, pData []byte) (n int, err error) {

	if err = c.Open(); err != nil {
		return -1, err
	}
	if _, err = c.handle.Seek(int64(offset), 0); err != nil {
		return -1, err
	}

	// read
	return c.handle.Read(pData)
}

// Write writes process data at a specific position, writes len(pData) bytes to file.
// Returns number of bytes read or error
func (c *RevPiControl) Write(offset uint32, pData []byte) (n int, err error) {
	if err = c.Open(); err != nil {
		return -1, err
	}
	if _, err = c.handle.Seek(int64(offset), 0); err != nil {
		return -1, err
	}

	// write
	return c.handle.Write(pData)
}

// GetDeviceInfo gets a description of a connected device.
func (c *RevPiControl) GetDeviceInfo(devInfo *SDeviceInfo) (result int, err error) {
	if err = c.Open(); err != nil {
		return 0, err
	}

	var r uintptr
	if r, _, err = ioctl(c.handle.Fd(), KB_GET_DEVICE_INFO, uintptr(unsafe.Pointer(devInfo))); err != nil {
		return 0, err
	}

	return int(r), nil
}

// GetDeviceInfoList gets a description of connected devices as an array of 20 elements.
// Returns the number of detected devices.
func (c *RevPiControl) GetDeviceInfoList() (devInfo []SDeviceInfo, err error) {
	if err = c.Open(); err != nil {
		return nil, err
	}
	asDevList := make([]SDeviceInfo, 255)
	var r uintptr
	if r, _, err = ioctl(c.handle.Fd(), KB_GET_DEVICE_INFO_LIST, uintptr(unsafe.Pointer(&asDevList[0]))); err != nil {
		return nil, err
	}

	// cut off the slice
	devInfo = asDevList[:int(r)]

	return devInfo, nil
}

// GetBitValue gets the value of one bit in the process image.
func (c *RevPiControl) GetBitValue(pSpiValue *SPIValue) (err error) {
	if err = c.Open(); err != nil {
		return err
	}

	pSpiValue.I16uAddress += uint16(pSpiValue.I8uBit) / 8
	pSpiValue.I8uBit %= 8

	if _, _, err = ioctl(c.handle.Fd(), KB_GET_VALUE, uintptr(unsafe.Pointer(pSpiValue))); err != nil {
		return err
	}
	return nil
}

// SetBitValue sets the value of one bit in the process image.
func (c *RevPiControl) SetBitValue(pSpiValue *SPIValue) (err error) {
	if err = c.Open(); err != nil {
		return err
	}

	pSpiValue.I16uAddress += uint16(pSpiValue.I8uBit) / 8
	pSpiValue.I8uBit %= 8

	if _, _, err = ioctl(c.handle.Fd(), KB_SET_VALUE, uintptr(unsafe.Pointer(pSpiValue))); err != nil {
		return err
	}
	return nil
}

// GetVariableInfo gets information about a variable by name.
func (c *RevPiControl) GetVariableInfo(name string) (pSpiVariable *SPIVariable, err error) {
	if err = c.Open(); err != nil {
		return nil, err
	}

	var v SPIVariable
	v.StrVarName = ByteToUint8Array(([]byte)(name))
	var r uintptr
	if r, _, err = ioctl(c.handle.Fd(), KB_FIND_VARIABLE, uintptr(unsafe.Pointer(&v))); err != nil {
		return nil, err
	}

	if int(r) < 0 {
		return nil, fmt.Errorf("could not find variable %s", name)
	}

	return &v, nil
}

// FindVariable checks if a variable with a specific name exists.
func (c *RevPiControl) FindVariable(name string) (found bool) {

	if _, err := c.GetVariableInfo(name); err != nil {
		return false
	}

	return true
}

// ResetCounter resets a counter.
func (c *RevPiControl) ResetCounter(address uint8, bitfield uint16) (result int, err error) {

	if err = c.Open(); err != nil {
		return -1, err
	}

	var tel SDIOResetCounter
	var r uintptr

	tel.I8uAddress = address
	tel.I16uBitfield = bitfield

	if r, _, err = ioctl(c.handle.Fd(), KB_DIO_RESET_COUNTER, uintptr(unsafe.Pointer(&tel))); err != nil {
		return int(r), err
	}

	if int(r) < 0 {
		return int(r), fmt.Errorf("could not reset counter")
	}

	return int(r), nil
}

// WaitForEvent waits for Reset of Pi Control Interface
func (c *RevPiControl) WaitForEvent() (err error) {
	var event int

	if err = c.Open(); err != nil {
		return err
	}

	if ioctl(c.handle.Fd(), KB_WAIT_FOR_EVENT, uintptr(unsafe.Pointer(&event))); err != nil {
		return err
	}
	return nil
}

// UpdateFirmware update a device firmware, check on the Kunubs website for details about updating firmware.
func (c *RevPiControl) UpdateFirmware(addrP uint32) (result int, err error) {

	if err = c.Open(); err != nil {
		return -1, err
	}

	var r uintptr

	if addrP == 0 {
		r, _, err = ioctl(c.handle.Fd(), KB_UPDATE_DEVICE_FIRMWARE, 0)
	} else {
		r, _, err = ioctl(c.handle.Fd(), KB_UPDATE_DEVICE_FIRMWARE, uintptr(unsafe.Pointer(&addrP)))
	}

	if err != nil {
		return int(r), err
	}

	if int(r) < 0 {
		return int(r), fmt.Errorf("firmware update failed")
	}

	cMsg := make([]byte, 255)
	if r, _, err = ioctl(c.handle.Fd(), KB_GET_LAST_MESSAGE, uintptr(unsafe.Pointer(&cMsg[0]))); err != nil && r == 0 && cMsg[0] != 0 {
		fmt.Println(string(cMsg))
	}

	return int(r), nil
}

// GetModuleName returns a friendly name for a RevPi module type.
func GetModuleName(moduletype uint16) string {
	moduletype = moduletype & PICONTROL_NOT_CONNECTED_MASK
	switch moduletype {
	case 95:
		return "RevPi Core"
	case 96:
		return "RevPi DIO"
	case 97:
		return "RevPi DI"
	case 98:
		return "RevPi DO"
	case 103:
		return "RevPi AIO"
	case PICONTROL_SW_MODBUS_TCP_SLAVE:
		return "ModbusTCP Slave Adapter"
	case PICONTROL_SW_MODBUS_RTU_SLAVE:
		return "ModbusRTU Slave Adapter"
	case PICONTROL_SW_MODBUS_TCP_MASTER:
		return "ModbusTCP Master Adapter"
	case PICONTROL_SW_MODBUS_RTU_MASTER:
		return "ModbusRTU Master Adapter"
	case 100:
		return "Gateway DMX"
	case 71:
		return "Gateway CANopen"
	case 73:
		return "Gateway DeviceNet"
	case 74:
		return "Gateway EtherCAT"
	case 75:
		return "Gateway EtherNet/IP"
	case 93:
		return "Gateway ModbusTCP"
	case 76:
		return "Gateway Powerlink"
	case 77:
		return "Gateway Profibus"
	case 79:
		return "Gateway Profinet IRT"
	case 81:
		return "Gateway SercosIII"
	default:
		return "unknown moduletype"
	}
}

// IsModuleConnected checks whether a RevPi module is conneted.
func IsModuleConnected(moduletype uint16) bool {
	return moduletype&PICONTROL_NOT_CONNECTED > 0
}
