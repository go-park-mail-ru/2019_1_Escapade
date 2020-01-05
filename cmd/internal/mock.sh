#!/bin/sh
echo ""
echo "  --------------------------"
echo "  ------Mock generate-------"
echo "  --------------------------"
echo ""
trap 'echo " stop" ' INT TERM
#chmod +x gomock.sh && ./gomock.sh
go get github.com/stretchr/testify/mock
go get github.com/vektra/mockery/.../
#apt-get install mockgen
export PATH=$PATH:$GOPATH/bin/mockery
go generate ./../../internal/...

# echo '1. pkg/database'
# mockgen -source=./../../internal/pkg/database/interfaces.go -destination=interfaces_mock.go -package=database
# mv interfaces_mock.go ./../../internal/pkg/database/

# echo '2. game/database'
# mockgen -source=./../../internal/services/game/database/interfaces.go -destination=interfaces_mock.go -package=database
# mv interfaces_mock.go ./../../internal/services/game/database/