# C/C++ cross-compilation utilities for RevPi

This folder contains the source code for the RevPi libraries used on Raspberry Pi.

Use:

```go
./buildme
```

to build, this has been adapted by the RaspberryPi buildme script used for cross-compilation, use [original RaspberryPi build tools](https://github.com/raspberrypi/tools).

It requires CMake to be installed and an arm cross compiler. It is set up to use this one:
<https://github.com/raspberrypi/tools/tree/master/arm-bcm2708/gcc-linaro-arm-linux-gnueabihf-raspbian>

On a Debian-based machine:

- Debian Jessie: <https://wiki.debian.org/CrossToolchains#In_jessie_.28Debian_8.29>
- Debian 9 stretch: `sudo apt-get install gcc-6-arm-linux-gnueabihf`

## In case of issues use the official RaspberryPi tools

If the standard packages do not work use directly ("segmentation fault" error for instance) use [original RaspberryPi build tools](https://github.com/raspberrypi/tools)

See also:

- Linux

<https://raspberrypi.stackexchange.com/questions/42015/cross-compile-error-using-arm-linux-gnueabihf-gcc>
<https://github.com/raspberrypi/tools/tree/master/arm-bcm2708/gcc-linaro-arm-linux-gnueabihf-raspbian-x64>

Add "tools/arm-bcm2708/gcc-linaro-arm-linux-gnueabihf-raspbian-x64/bin" to you PATH.

- Windows, if you are really up to it:

<http://gnutoolchains.com/raspberry/tutorial/>
