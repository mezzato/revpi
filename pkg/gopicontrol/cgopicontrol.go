// +build ignore

package gopicontrol

//go:generate sh -c "go tool cgo -godefs=true ctypes_linux.go > types_linux.go"

// #cgo LDFLAGS: -L${SRCDIR}/../../kunbus/build/native/release/interface/piControl -lpiControl_static
// #cgo CFLAGS: -I${SRCDIR}/../../kunbus/interface/piControl
// #include <piControlIf.h>
// #include <stdlib.h>
// #include <string.h>
import "C"
import (
	"errors"
	"fmt"
	"unsafe"
)

var WriteError []string = []string{
	"Cannot connect to control process",
	"Offset seek error",
	"Cannot write to control process",
	"Unknown error",
}

// MyError is an error implementation that includes a time and message.
type PiControlWriteError struct {
	Code int
}

func (e PiControlWriteError) Error() string {
	switch e.Code {
	case -1:
		return WriteError[0]
	case -2:
		return WriteError[1]
	case -3:
		return WriteError[2]
	default:
		return WriteError[3]
	}
}

type PiControlReadError struct {
	Code int
}

func (e PiControlReadError) Error() string {
	return fmt.Sprintf("read error, code:%d", e.Code)
}

func Reset() (err error) {
	if c := int(C.piControlReset()); c < 0 {
		return fmt.Errorf("error resetting driver: %s", StrError(c))
	}
	return nil
}

func StrError(code int) string {
	return C.GoString(C.strerror(C.int(code)))
}

func Read(offset uint32, length uint32, pData []byte) (err error) {

	r := int(C.piControlRead(C.uint32_t(offset), C.uint32_t(length), (*C.uint8_t)(unsafe.Pointer(&pData[0]))))

	if r < 0 {
		return PiControlReadError{Code: r}
	}
	return nil
}

func Write(offset uint32, length uint32, pData []byte) (err error) {
	r := int(C.piControlWrite(C.uint32_t(offset), C.uint32_t(length), (*C.uint8_t)(unsafe.Pointer(&pData[0]))))

	if r < 0 {
		return PiControlWriteError{Code: r}
	}
	return nil
}

func GetDeviceInfo(devInfo *SDeviceInfo) (result int, err error) {
	return int(C.piControlGetDeviceInfo((*C.struct_SDeviceInfoStr)(unsafe.Pointer(devInfo)))), nil
}

func GetDeviceInfoList(devInfo []SDeviceInfo) (devcount int, err error) {
	r := int(C.piControlGetDeviceInfoList((*C.SDeviceInfo)(unsafe.Pointer(&devInfo[0]))))
	if r < 0 {
		return 0, fmt.Errorf("error getting device info list, error: %s", StrError(r))
	}
	return r, nil
}

func GetBitValue(pSpiValue *SPIValue) (err error) {
	r := int(C.piControlGetBitValue((*C.struct_SPIValueStr)(unsafe.Pointer(pSpiValue))))

	if r < 0 {
		return PiControlReadError{Code: r}
	}
	return nil
}

func SetBitValue(pSpiValue *SPIValue) (err error) {
	r := int(C.piControlSetBitValue((*C.struct_SPIValueStr)(unsafe.Pointer(pSpiValue))))

	if r < 0 {
		return PiControlWriteError{Code: r}
	}
	return nil
}

func GetVariableInfo(name string, pSpiVariable *SPIVariable) (err error) {

	pSpiVariable.StrVarName = ByteToInt8Array(([]byte)(name))

	v := (*C.struct_SPIVariableStr)(unsafe.Pointer(pSpiVariable))

	r := int(C.piControlGetVariableInfo(v))

	if r < 0 {
		return errors.New(fmt.Sprintf("cannot find variable '%s'", name))
	}

	return nil
}

func FindVariable(name string) (err error) {
	// s := "RevPiLED"
	cs := C.CString(name)
	defer C.free(unsafe.Pointer(cs))
	r := int(C.piControlFindVariable(cs))
	if r < 0 {
		return errors.New(fmt.Sprintf("cannot find variable '%s'", name))
	}

	return nil
}

func ResetCounter(address int, bitfield int) (result int, err error) {
	return int(C.piControlResetCounter(C.int(address), C.int(bitfield))), nil
}

func WaitForEvent() (result int, err error) {
	return int(C.piControlWaitForEvent()), nil
}

func UpdateFirmware() (result int, err error) {
	return int(C.piControlUpdateFirmware()), nil
}
