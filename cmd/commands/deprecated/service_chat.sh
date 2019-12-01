#!/bin/sh
echo ""
echo "  --------------------------"
echo "  ---Launch chat service----"
echo "  --------------------------"
echo ""
trap 'echo " stop" ' INT TERM
#chmod +x service_chat.sh && ./service_chat.sh

done=0

echo "  1. Build app"
go build -o bin/chat ../chat/main.go &&\

echo "  2. Run app" &&\
# 1 parameter - path to main configuration json file
# 2 parameter - port of this server
# 3 parameter - port of Consul server
./bin/chat ../chat/chat.json 3041 8500
done=1
    
echo "  3. Remove app" 
rm bin/chat

echo ""
if [ "$done" -eq 1 ]
then 
echo "  ----------Done!-----------"
else
echo "  ----------Error!-----------"
exit 1
fi