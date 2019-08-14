#!/bin/sh
echo ""
echo "  --------------------------"
echo "  ---Launch api service----"
echo "  --------------------------"
echo ""
trap 'echo " stop" ' INT TERM
#chmod +x service_api.sh && ./service_api.sh

done=0

echo "  1. Build app"
go build -o bin/api ../api/main.go &&\

echo "  2. Run app" &&\
# 1 parameter - path to main configuration json file
# 2 parameter - path to public photo configuration
# 3 parameter - path to private photo configuration. If you dont know it, set random string
# 4 parameter - port of this server
# 5 parameter - port of Consul server
./bin/api ../api/api.json ../internal/photo/photo.json ../secret.json 3010 8500 &&\
done=1
    
echo "  3. Remove app"
rm bin/api

echo ""
if [ "$done" -eq 1 ]
then 
echo "  ----------Done!-----------"
else
echo "  ----------Error!-----------"
exit 1
fi