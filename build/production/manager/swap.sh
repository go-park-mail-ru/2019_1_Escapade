#!/bin/sh
echo ""
echo "  --------------------------"
echo "  ---- Create swap ----"
echo "  --------------------------"
echo ""
trap 'echo " stop" ' INT TERM
#chmod +x install.sh && ./install.sh

place=$(cat /proc/swaps | grep /swapfile)
if [ "$?" = "0" ]; then
    echo "✔ " $place
else
    # https://sheensay.ru/swap
    df -h
    fallocate -l 2G /swapfile
    ls -lh /swapfile
    chmod 600 /swapfile
    ls -lh /swapfile
    mkswap /swapfile
    swapon /swapfile
    swapon -s
fi