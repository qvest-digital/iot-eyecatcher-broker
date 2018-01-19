#!/usr/bin/env bash

go get github.com/golang/mock/gomock
go get github.com/golang/mock/mockgen

# should be in your path
MOCKGEN=mockgen
SED=sed
GOFMT=gofmt
MKDIR=mkdir

generate_mock() {
    SRC=$1
    PKG=$(dirname $SRC)
    DST=$PKG/mock_$(basename $SRC)

    $MKDIR -p $(dirname $DST)
    $MOCKGEN -source ./$SRC -destination ./$DST -package $(basename $PKG)
    $GOFMT -w ./$DST
}

generate_vendor_mock() {
    PKG=$1
    INTERFACES=$2
    DST=mocks/$PKG/mock_$(basename $PKG).go

    $MKDIR -p $(dirname $DST)
    $MOCKGEN -destination ./$DST -package $(basename $PKG) $PKG $INTERFACES
    $GOFMT -w ./$DST
}

# generate project/internal mocks
generate_mock broker/interfaces.go
generate_mock ws/interfaces.go
generate_mock json/interfaces.go

# generate vendor mocks
#generate_vendor_mock github.com/docker/machine/libmachine API

# generate go intrinsic mocks
generate_vendor_mock net Addr
