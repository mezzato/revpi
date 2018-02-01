package main

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/mezzato/revpi/pkg/gopicontrol"
)

func writeVariableValue(ctrl *gopicontrol.RevPiControl, variableName string, v uint32) (err error) {

	var (
		sPIValue    gopicontrol.SPIValue
		sPiVariable *gopicontrol.SPIVariable
	)

	sPiVariable, err = ctrl.GetVariableInfo(variableName)
	if err != nil {
		return
	}

	var data interface{}
	if sPiVariable.I16uLength == 1 {
		sPIValue.I16uAddress = sPiVariable.I16uAddress
		sPIValue.I8uBit = sPiVariable.I8uBit
		sPIValue.I8uValue = uint8(v)
		data = sPIValue.I8uValue
		if err = ctrl.SetBitValue(&sPIValue); err != nil {
			return
		}
	} else {
		switch sPiVariable.I16uLength {
		case 8:
			data = uint8(v)
		case 16:
			data = uint16(v)
		case 32:
			data = v
		}
		b, e := gopicontrol.NumToBytes(data)
		if e != nil {
			return e
		}

		if _, err = ctrl.Write(uint32(sPiVariable.I16uAddress), b); err != nil {
			return
		}
	}

	fmt.Printf("written value %d dec (=%02x hex) to offset %d.\n", data, data, sPiVariable.I16uAddress)
	return nil
}
func showVariableInfo(ctrl *gopicontrol.RevPiControl, variableName string) (err error) {
	sPiVariable, err := ctrl.GetVariableInfo(variableName)
	if err != nil {
		return
	}
	fmt.Printf("variable name: %s\n", sPiVariable.StrVarName)
	fmt.Printf("       offset: %d\n", sPiVariable.I16uAddress)
	fmt.Printf("       length: %d\n", sPiVariable.I16uLength)
	fmt.Printf("          bit: %d\n", sPiVariable.I8uBit)

	return nil
}

func readVariableValue(ctrl *gopicontrol.RevPiControl, variableName string, format byte, quiet bool) (err error) {
	var (
		sPIValue    gopicontrol.SPIValue
		sPiVariable *gopicontrol.SPIVariable
		i32uValue   uint32
	)

	sPiVariable, err = ctrl.GetVariableInfo(variableName)
	if err != nil {
		return
	}
	if sPiVariable.I16uLength == 1 {
		sPIValue.I16uAddress = sPiVariable.I16uAddress
		sPIValue.I8uBit = sPiVariable.I8uBit

		err = ctrl.GetBitValue(&sPIValue)
		if err != nil {
			return
		}

		if !quiet {
			fmt.Printf("Bit value: %d\n", sPIValue.I8uValue)
		} else {
			fmt.Printf("%d\n", sPIValue.I8uValue)
		}

	} else {
		sizeRemainder := sPiVariable.I16uLength % 8
		if sizeRemainder != 0 {
			return fmt.Errorf("could not read variable %s. Internal Error", variableName)
		}
		size := sPiVariable.I16uLength / 8

		switch sPiVariable.I16uLength {
		case 8, 16, 32:
			data := make([]byte, size)
			if _, err = ctrl.Read(uint32(sPiVariable.I16uAddress), data); err != nil {
				return
			}
			// fmt.Printf("read from address %d, data: %x\n", sPiVariable.I16uAddress, data)
			buf := bytes.NewReader(data)
			switch sPiVariable.I16uLength {
			case 8:
				var ui8 uint8
				err = binary.Read(buf, binary.LittleEndian, &ui8)
				i32uValue = uint32(ui8)
			case 16:
				var ui16 uint16
				err = binary.Read(buf, binary.LittleEndian, &ui16)
				i32uValue = uint32(ui16)
			case 32:
				err = binary.Read(buf, binary.LittleEndian, &i32uValue)
			}
			if err != nil {
				return
			}

			if format == 'h' {
				if !quiet {
					//f := "%d byte-value of %s: %0" + strconv.Itoa(size*2) + "x hex (=%d dec)\n"
					fmt.Printf("%d byte-value of %s: %x hex bytes (=%d dec)\n", size, variableName, data, i32uValue)
				} else {
					fmt.Printf("%x\n", i32uValue)
				}
			} else if format == 'b' {
				if !quiet {
					fmt.Printf("%d byte value of %s: ", size, variableName)
				}

				bn, _ := gopicontrol.NumToBytes(i32uValue)
				fmt.Printf("binary value:% x\n", bn)
			} else {
				if !quiet {
					//f := "%d byte-value of %s: %d dec (=%0" + strconv.Itoa(size*2) + "x hex)\n"
					fmt.Printf("%d byte-value of %s: %d dec (=%x hex bytes)\n", size, variableName, i32uValue, data)
				} else {
					fmt.Printf("%d\n", i32uValue)
				}
			}
		default:
			return fmt.Errorf("invalid byte size %d for variable %s", size, variableName)
		}

	}
	return nil
}
