#!/bin/sh
echo ""
echo "  --------------------------"
echo "  ---Launch auth service----"
echo "  --------------------------"
echo ""
trap 'echo " stop" ' INT TERM
#chmod +x service_auth.sh && ./service_auth.sh

done=0

echo "  1. Build app"
go build -o bin/auth ../auth/main.go &&\

echo "  2. Run app" &&\
# 1 parameter - path to main configuration json file
# 2 parameter - port of this server
# 3 parameter - port of Consul server
# 4 parameter - port of auth postgresql server
./bin/auth ../auth/auth.json 3045 8500 2345 &&
done=1
    
echo "  3. Remove app" 
rm bin/auth

echo ""
if [ "$done" -eq 1 ]
then 
echo "  ----------Done!-----------"
else
echo "  ----------Error!-----------"
exit 1
fi