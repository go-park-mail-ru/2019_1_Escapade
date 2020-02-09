#!/bin/sh
echo ""
echo "  --------------------------"
echo "  ------Mock generate-------"
echo "  --------------------------"
echo ""
#chmod +x gomock.sh && ./gomock.sh

export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin:/usr/local/go/bin

echo "Install mockery"
place=$(whereis mockery | grep bin)
if [ "$?" = "0" ]; then
    echo "✔ " $place "\n"
else
    go get github.com/vektra/mockery/.../
    echo "\n✔ installed \n"
fi

echo "Apply mockery\n"
export PATH=$PATH:$GOPATH/bin/mockery
go generate ./../../internal/...
