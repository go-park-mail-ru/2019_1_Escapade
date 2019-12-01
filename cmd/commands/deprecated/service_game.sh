#!/bin/sh
echo ""
echo "  --------------------------"
echo "  ---Launch game service----"
echo "  --------------------------"
echo ""
trap 'echo " stop" ' INT TERM
#chmod +x service_game.sh && ./service_game.sh

done=0

echo "  1. Build app"
go build -o bin/game ../game/main.go &&\

echo "  2. Run app" &&\
# 1 parameter - path to main configuration json file
# 2 parameter - path to public photo configuration
# 3 parameter - path to private photo configuration
# 4 parameter - path to field constants(need for game)
# 5 parameter - path to room constants(need for game)
# 6 parameter - port of this server
# 7 parameter - port of Consul server
./bin/game ../game/game.json ../internal/photo/photo.json \
    ../secret.json ../internal/constants/field.json \
    ../internal/constants/room.json 3002 8500 &&\

#./bin/game cmd/services/game/game.json internal/photo/photo.json secret.json internal/constants/field.json internal/constants/room.json 3002 8500 &&\
done=1
    
echo "  3. Remove app" 
rm game

echo ""
if [ "$done" -eq 1 ]
then 
echo "  ----------Done!-----------"
else
echo "  ----------Error!-----------"
exit 1
fi