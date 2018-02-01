# setup environment for cross compile to arm-linux

if (DEFINED CMAKE_TOOLCHAIN_FILE)
else()
   message(WARNING
	"  *********************************************************\n"
   	"  *   CMAKE_TOOLCHAIN_FILE not defined                    *\n"
	"  *   This is correct for compiling on the Raspberry Pi   *\n"
	"  *                                                       *\n"
	"  *   If you are cross-compiling on some other machine    *\n"
	"  *   then DELETE the build directory and re-run with:    *\n"
	"  *   -DCMAKE_TOOLCHAIN_FILE=toolchain_file.cmake         *\n"
	"  *                                                       *\n"
   	"  *   Toolchain files are in makefiles/cmake/toolchains.  *\n"
	"  *********************************************************"
       )
endif()

set(SHARED "SHARED")


# All linux systems have sbrk()
add_definitions(-D_HAVE_SBRK)

# pull in declarations of lseek64 and friends
add_definitions(-D_LARGEFILE64_SOURCE)
	
# test for glibc malloc debugging extensions
try_compile(HAVE_MTRACE
            ${CMAKE_BINARY_DIR}
            ${PROJECT_SOURCE_DIR}/makefiles/cmake/srcs/test-mtrace.c
            OUTPUT_VARIABLE foo)

# test for existence of execinfo.h header
include(CheckIncludeFile)
check_include_file(execinfo.h HAVE_EXECINFO_H)

add_definitions(-DHAVE_CMAKE_CONFIG)
configure_file (
    "makefiles/cmake/cmake_config.h.in"
    "${PROJECT_BINARY_DIR}/cmake_config.h"
    )
 
