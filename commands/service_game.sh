#!/bin/sh
echo ""
echo "  --------------------------"
echo "  ---Launch game service----"
echo "  --------------------------"
echo ""
#chmod +x service_game.sh && ./service_game.sh

echo "  1. Build app"
go build -o game ../game/main.go

echo "  1. Run app"
# 1 parameter - path to main configuration json file
# 2 parameter - path to public photo configuration
# 3 parameter - path to private photo configuration
# 4 parameter - path to field constants(need for game)
# 4 parameter - path to room constants(need for game)
./game ../game/game.json ../internal/photo/photo.json \
    ../secret.json ../internal/constants/field.json \
    ../internal/constants/room.json
    
echo "  3. Remove app"
rm game

echo "  ----------Done!-----------"