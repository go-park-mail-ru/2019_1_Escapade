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
go build -o chat ../chat/main.go &&\

echo "  2. Run app" &&\
# 1 parameter - path to main configuration json file
# 2 parameter - path to public photo configuration
# 3 parameter - path to private photo configuration
# 4 parameter - path to field constants(need for game)
# 4 parameter - path to room constants(need for game)
./chat ../chat/chat.json &&\
done=1
    
echo "  3. Remove app" 
rm chat

echo ""
if [ "$done" -eq 1 ]
then 
echo "  ----------Done!-----------"
else
echo "  ----------Error!-----------"
exit 1
fi