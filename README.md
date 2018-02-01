# Go library for Revolution Pi

**This first code base is still experimental. The Go library has not been tested for a firmware update.**

Please post any issues you may find, I will try to investigate them even though I can not promise to be very prompt due to lack of time, feel free to create a pull request with a fix.

Refer to the Go doc in the code stored in the folders:

- [gopicontrol package](pkg/gopicontrol): this is a Go port of the piControl C/C++ driver methods wrapped in a Go object. Native syscalls have been used to access the process image in the kernel. These are a one-to-one translate of the piControl.c interface.
- [gopitest application](cmd/gopitest): this is a Go sample application which mimics the functionality of the piTest C application available as a standard command line tool using piControl.

## Go library for Revolution Pi and cross-compilation tools

This project was originally born to port the RevPi piControl C library to the Go language. The initial attempt included porting C structures to Go code and use cgo to keep the RevPi original code base.

For this reason this project besides the Go piControl port contains CMake scripts to cross-compile C/C++ applications for RaspberryPi (RevPi Core).
Readers interested in either topic should refer to the REAMDE files in the specific folders. Below the two areas and more details about this project if you are keen to know its origin.

### How to use the Go package and cross-compile for the RevPi

Install Go on any Linux-based system.

To use the Go package inside your project you just need to refer to the gopicontrol package by using a standard Go import like:

```go
import "github.com/mezzato/revpi/pkg/gopicontrol"
```

It is normally faster to cross-compile the code on a decent machine and upload it to the RevPi, eg launch:

```go
GOOS=linux GOARCH=arm GOARM=6 go build
```

There is a running example of this, the [gopitest application](cmd/gopitest).
Build the app with the `build.sh` script and try it out, for help:

```go
./gopitest -h
```

To switch on the internal LED:

```go
./gopitest write -n RevPiLED -v 2
```

To read the internal LED value:

```go
./gopitest read -n RevPiLED
```

### How to keep the Go code in sync with the piControl C headers

There is a shell script [generate_godefs.sh](pkg/gopicontrol/generate_godefs.sh) to generate Go structs and constants from the C headers via cgo `-godefs` option.
You just need to keep aligned the file [ctypes_linux.go](pkg/gopicontrol/ctypes_linux.go) with the piControl headers to export what you need.

IMPORTANT: **You should launch the script on the RevPi not on a generic Linux machine.**

### Why not a cgo direct port

The first attempt was to use cgo to convert the code, there is still a working example: [cgopicontrol.go](pkg/gopicontrol/cgopicontrol.go). The file has been now excluded from the go build.

This works but it is tricky and has downsides as mentioned by the Go team as well: [Dave Cheney: cgo-is-not-go](https://dave.cheney.net/2016/01/18/cgo-is-not-go).

Nevertheless, it provides a working package, before compiling with cgo you need to compile the C code to a library and that's part of the following topic.

## C/C++ cross-compilation utilities

The standard RevPi examples are based on make files, while this is the standard way of building C code CMake provides simplicity and flexibility, particularly if you need to build your own solution and you are interested in cross-compilation.

Refer to [REAMDE in the kunbus folder](kunbus/README.md), in particular to the [buildme](kunbus/buildme) script to see native and cross-compilation options with CMake.

## Why a Go port

Based on previous experience with C++ and Python I prefer to use Go when possible and sensible, I thought it reasonable to use it for a more modern look and feel to a project like RevPi born from open source technology but geared to the industrial use.

Some years ago I decided to build my own home automation software and went for Go, a decision I still don't regret since the code I wrote still runs beautifully on a RaspberryPi after years with virtually no code changes. When the RevPi Core was released I was eager to replace my Arduino/Raspberry Pi prototype with a more robust solution which could be fully controlled through the Web and would be based on a world-tested platform, that is why the choice fell in favour of the RevPi.

The idea is then to use my Go experience to contribute some code to the open source community, I personally find in many cases that Go is more fun than C++ and Python, event if I do not suggest it be used in all situations. As always it depends on what you are up to. Here are my motivations, I hope they help make your choice.

- Python is dynamically interpreted but runtime errors are best avoided in microcontrollers, the more the compiler can do for you the better. Duck-typing is good for quick scripting but extra care is needed when building reliable services that run round-the-clock. Of course you need to have a Python runtime on the target machine.
- C++ and C memory management and compilation need plenty of attention. If you need full memory control and speed consider using these languages though.
- Go like Python is compact, open source and batteries included.
- Go static compilation produces a single file which can be copied and pasted, cross-compilation is possible too.
- Modern applications involve more and more concurrency patterns. Go has a modern API and nice concurrency support built into the language, threads and locks are normally not needed, which are a common source of pitfalls and hard-to-find live issues.
- It is fun to use Go for me.
