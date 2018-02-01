export CFLAGS="-I`pwd`/../../kunbus/interface/piControl"
go tool cgo -godefs=true -- $CFLAGS ctypes_linux.go > types_linux.go
