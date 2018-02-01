# Run this file only once per cmake run. This is so that
# users projects can override those global settings before
# parsing subdirectories (which can also include this file
# directly so that they can be built as standalone projects).
if (DEFINED GLOBAL_SETTINGS_SET)
   return ()
endif ()
set (GLOBAL_SETTINGS_SET "TRUE")
