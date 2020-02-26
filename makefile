CC=go
RM=rm
MV=mv


SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

VERSION:=$(shell grep -m1 "version" *.go | sed 's/[", ]//g' | cut -d= -f2)
suffix=$(shell grep -m1 "version" *.go | sed 's/[", ]//g' | cut -d= -f2 | sed 's/[0-9.]//g')
snapshot=$(shell date +%FT%T)

ifeq ($(suffix),rc)
	appversion=$(VERSION)$(snapshot)
else 
	appversion=$(VERSION)
endif 

.DEFAULT_GOAL:=build


build: 
#	@echo "Compilation for linux"
#	CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o m4backup $(SOURCEDIR)/main.go
#	zip m4backup-$(appversion)-linux.zip m4backup 
	@echo "Compilation for windows"
	CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc  CXX=x86_64-w64-mingw32-g++ GOOS=windows GOARCH=amd64  go build ${LDFLAGS} -o m4backup.exe $(SOURCEDIR)/main.go
	zip m4backup-$(appversion)-windows.zip -j m4backup.exe windows/*
	@echo "Compilation for macos"
	GOOS=darwin go build ${LDFLAGS} -o m4backup $(SOURCEDIR)/main.go
	zip m4backup-$(appversion)-macos.zip m4backup 
