set(PI_TEST_SOURCE
   piTest.c)


add_executable(piTest ${PI_TEST_SOURCE})
target_link_libraries (piTest piControl_static)

install (TARGETS piTest DESTINATION bin)