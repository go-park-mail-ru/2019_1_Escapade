#!/bin/sh
echo ""
echo "  --------------------------"
echo "  ------Proto generate------"
echo "  --------------------------"
echo ""
#chmod +x proto.sh && ./proto.sh

export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin:/usr/local/go/bin

echo "Install protoc\n"
place=$(whereis protoc | grep bin)
if [ "$?" = "0" ]; then
    echo "✔ " $place "\n"
else
    go get github.com/golang/protobuf
    go get -u github.com/golang/protobuf/protoc-gen-go
    yes | apt install golang-goprotobuf-dev
    echo "\n✔ installed \n"
fi

done=0

go mod vendor

export PATH=$PATH:/usr/bin/protoc
export CHATPROTO=$PWD/../../internal/services/chat/proto
export PROTO=$PWD/../../vendor

echo "  1. Copy .proto files to protobuf directory"

cp $CHATPROTO/chat.proto $PROTO && \

echo "  2.1 Apply proto to chat" && \
cd $PROTO && \
protoc  --go_out=plugins=grpc:. chat.proto && \
cp $PROTO/chat.pb.go $CHATPROTO && \
done=1

echo "  3. Remove our .proto and .go files from GOPATH" && \
rm $PROTO/chat.pb.go && \
rm $PROTO/chat.proto && 

echo ""
if [ "$done" -eq 1 ]
then 
echo "  ----------Done!-----------"
else
echo "  ----------Error!-----------"
exit 1
fi