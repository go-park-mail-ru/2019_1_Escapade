#!/bin/sh
echo ""
echo "  --------------------------"
echo "  ------Proto generate------"
echo "  --------------------------"
echo ""
#chmod +x proto.sh && ./proto.sh

export GOPATH=$HOME/go
export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
done=0

#echo "  0. Download protobuf "
#go get github.com/golang/protobuf
#go get github.com/golang/protobuf/proto
#go get -u github.com/golang/protobuf/protoc-gen-go

export CHATPROTO=$PWD/../internal/services/chat/
export PROTO=$GOPATH/src

echo "  1. Copy .proto files to protobuf directory"

cp $CHATPROTO/chat.proto $PROTO && \

echo "  2.1 Apply proto to chat" && \
cd $PROTO && \
protoc --go_out=plugins=grpc:. chat.proto && \
cp $PROTO/chat.pb.go $CHATPROTO && \

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