#!/bin/sh
echo ""
echo "  -------------------------------"
echo "  --Add nodes to docker-machine--"
echo "  -------------------------------"
echo ""
trap 'echo " stop" ' INT TERM
#chmod +x prepare.sh && ./prepare.sh

echo "  0. Prepare other .sh"
chmod +x ./one/add.sh
# $1 - $5 - адреса серверов
addr1=$1
addr2=$2
addr3=$3
addr4=$4
addr5=$5

# Укажите в переменной path путь до ssh ключа.
# Либо оставьте "id_rsa", но скопируйте id_rsa и id_rsa.pub
# в данный каталог 2019_1_Escapade/cmd/commands/
path="id_rsa"

# названия машин
name1="api1"
name2="api2"
name3="api3"
name4="api4"
name5="api5"

echo "  1. Create machines"

 echo "  1.1. Create machine - "$addr1
 ./one/add.sh $addr1 $path $name1 && \

echo "  1.2. Create machine - "$addr2
 ./one/add.sh $addr2 $path $name2 && \

echo "  1.3. Create machine - "$addr3
./one/add.sh $addr3 $path $name3 && \

echo "  1.4. Create machine - "$addr4
 ./one/add.sh $addr4 $path $name4 && \

echo "  1.5. Create machine - "$addr5
 ./one/add.sh $addr5 $path $name5
