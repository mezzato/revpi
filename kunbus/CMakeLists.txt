cmake_minimum_required(VERSION 2.8)

project(revpi)

SET(PROJECT_VER_MAJOR 1)
SET(PROJECT_VER_MINOR 0)
SET(PROJECT_VER_PATCH 0)
SET(PROJECT_VER "${PROJECT_VER_MAJOR}.${PROJECT_VER_MINOR}.${PROJECT_VER_PATCH}")
SET(PROJECT_APIVER "${PROJECT_VER}")

include(makefiles/cmake/global_settings.cmake)
include(makefiles/cmake/arm-linux.cmake)

enable_language(ASM)

# Global include paths
include_directories(${PROJECT_SOURCE_DIR})
include_directories(${PROJECT_BINARY_DIR})
include_directories("${PROJECT_SOURCE_DIR}/interface/piControl")

# Global compiler flags
if(CMAKE_COMPILER_IS_GNUCC)
   set(CMAKE_C_FLAGS "${CMAKE_C_FLAGS} -Wno-multichar -Wall -Wno-unused-but-set-variable -fPIC")
endif()

add_definitions(-D_REENTRANT)
add_definitions(-DUSE_VCHIQ_ARM -DVCHI_BULK_ALIGN=1 -DVCHI_BULK_GRANULARITY=1)
add_definitions(-DOMX_SKIP64BIT)
add_definitions(-DEGL_SERVER_DISPMANX)
add_definitions(-D_LARGEFILE_SOURCE -D_LARGEFILE64_SOURCE -D_FILE_OFFSET_BITS=64)
add_definitions(-D_GNU_SOURCE)

# piControl library
add_subdirectory(interface/piControl) 

# piTest
add_subdirectory(piTest)


include_directories("${PROJECT_BINARY_DIR}")
include(FindPkgConfig QUIET)
if(PKG_CONFIG_FOUND)
	# Produce a pkg-config file
	foreach(PCFILE revpi.pc)
		configure_file("pkgconfig/${PCFILE}.in" "${PCFILE}" @ONLY)
		install(FILES       "${CMAKE_CURRENT_BINARY_DIR}/${PCFILE}"
			DESTINATION "${CMAKE_INSTALL_PREFIX}/lib/pkgconfig")
	endforeach()
endif()

