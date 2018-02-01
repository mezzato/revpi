set(PICONTROL_SOURCE
   piControlIf.c
   piControlIf.cpp
   piControlIf.hpp
)

add_library(piControl_static STATIC ${PICONTROL_SOURCE})
add_library(piControl SHARED ${PICONTROL_SOURCE})

install(TARGETS piControl DESTINATION lib)
install(TARGETS piControl_static DESTINATION lib)

install(FILES
   piControl.h
   piControlIf.h
   piControlIf.hpp
   DESTINATION include/interface/piControl
)