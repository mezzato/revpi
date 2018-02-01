package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mezzato/revpi/pkg/gopicontrol"
)

func main() {

	// Subcommands
	readCmd := flag.NewFlagSet("read", flag.ExitOnError)
	writeCmd := flag.NewFlagSet("write", flag.ExitOnError)
	lsCmd := flag.NewFlagSet("ls", flag.ExitOnError)
	resetCmd := flag.NewFlagSet("reset", flag.ExitOnError)
	variableCmd := flag.NewFlagSet("variable", flag.ExitOnError)

	readCmdVarName := readCmd.String("n", "", "variable name. (required)")
	readCmdVarFormat := readCmd.String("f", "d", "variable format. (optional)")

	// List subcommand flag pointers
	writeCmdVarName := writeCmd.String("n", "", "variable name. (required)")
	writeCmdVarValue := writeCmd.Uint64("v", 0, "variable value. (required)")

	variableCmdVarName := variableCmd.String("n", "", "variable name. (required)")

	// Verify that a subcommand has been provided
	// os.Arg[0] is the main command
	// os.Arg[1] will be the subcommand
	if len(os.Args) < 2 || os.Args[1] == "-h" {
		fmt.Printf(`a subcommand is required, valid options are [read|write|ls|reset]:
read:     read variable value
write:    write variable value
variable: show variable info
ls:       list devices
reset:    reset the driver

Type 
%s -h
for help with a verb

For example to read the RevPi Core LED:
%s read -n RevPiLED
`, os.Args[0])
		os.Exit(1)
	}

	// Switch on the subcommand
	// Parse the flags for appropriate FlagSet
	// FlagSet.Parse() requires a set of arguments to parse as input
	// os.Args[2:] will be all arguments starting after the subcommand at os.Args[1]
	switch os.Args[1] {
	case "write":
		writeCmd.Parse(os.Args[2:])
	case "read":
		readCmd.Parse(os.Args[2:])
	case "variable":
		variableCmd.Parse(os.Args[2:])
	case "ls":
		lsCmd.Parse(os.Args[2:])
	case "reset":
		resetCmd.Parse(os.Args[2:])
	default:
		fmt.Printf("invalid command\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	rpctl := gopicontrol.NewRevPiControl()
	defer rpctl.Close()

	// Check which subcommand was Parsed using the FlagSet.Parsed() function. Handle each case accordingly.
	// FlagSet.Parse() will evaluate to false if no flags were parsed (i.e. the user did not provide any flags)
	if writeCmd.Parsed() {
		// Required Flags
		if *writeCmdVarName == "" {
			writeCmd.PrintDefaults()
			os.Exit(1)
		}

		fmt.Printf("writing variable: %s, value: %d\n", *writeCmdVarName, *writeCmdVarValue)
		if err := writeVariableValue(rpctl, *writeCmdVarName, (uint32)(*writeCmdVarValue)); err != nil {
			fmt.Println(err)
			return
		}

	}

	if readCmd.Parsed() {
		// Required Flags
		if *readCmdVarName == "" {
			readCmd.PrintDefaults()
			os.Exit(1)
		}

		fmt.Printf("reading variable: %s\n", *readCmdVarName)
		if err := readVariableValue(rpctl, *readCmdVarName, (*readCmdVarFormat)[0], false); err != nil {
			fmt.Println(err)
			return
		}
	}

	if variableCmd.Parsed() {
		// Required Flags
		if *variableCmdVarName == "" {
			variableCmd.PrintDefaults()
			os.Exit(1)
		}

		fmt.Printf("reading variable info: %s\n", *variableCmdVarName)
		if err := showVariableInfo(rpctl, *variableCmdVarName); err != nil {
			fmt.Println(err)
			return
		}
	}

	if lsCmd.Parsed() {

		devices, err := rpctl.GetDeviceInfoList()
		if err != nil {
			fmt.Println(err)
			return
		}

		if err := showDeviceList(devices); err != nil {
			fmt.Println(err)
			return
		}
	}

	if resetCmd.Parsed() {

		if err := rpctl.Reset(); err != nil {
			fmt.Println(err)
			return
		}

	}

}

func showDeviceList(asDevList []gopicontrol.SDeviceInfo) (err error) {
	devcount := len(asDevList)

	fmt.Printf("Found %d devices:\n", devcount)
	for dev := 0; dev < devcount; dev++ {
		mn := gopicontrol.GetModuleName(asDevList[dev].I16uModuleType)

		// Show device number, address and module type
		fmt.Printf("Address: %d module type: %d (0x%x) %s V%d.%d\n", asDevList[dev].I8uAddress,
			asDevList[dev].I16uModuleType, asDevList[dev].I16uModuleType,
			mn,
			asDevList[dev].I16uSW_Major, asDevList[dev].I16uSW_Minor)

		if asDevList[dev].I8uActive > 0 {
			fmt.Printf("Module is present\n")
		} else {
			if gopicontrol.IsModuleConnected(asDevList[dev].I16uModuleType) {
				fmt.Printf("Module is NOT present, data is NOT available!!!\n")
			} else {
				fmt.Printf("Module is present, but NOT CONFIGURED!!!\n")
			}
		}

		// Show offset and length of input section in process image
		fmt.Printf("     input offset: %d length: %d\n", asDevList[dev].I16uInputOffset,
			asDevList[dev].I16uInputLength)

		// Show offset and length of output section in process image
		fmt.Printf("    output offset: %d length: %d\n", asDevList[dev].I16uOutputOffset,
			asDevList[dev].I16uOutputLength)
		fmt.Printf("\n")
	}

	return nil
}
