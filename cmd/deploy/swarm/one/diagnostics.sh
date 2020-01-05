#!/bin/sh
echo ""
echo "  -------------------------------------"
echo "  ---$2 check ports open---"
echo "  -------------------------------------"
echo ""
trap 'echo " stop" ' INT TERM
#chmod +x droplet_diagnostics.sh && ./droplet_diagnostics.sh 1

#$1 - the IP of this node
#$2 - the IP of tested node

# -v Вывод информации о процессе работы (verbose)
# -z Не посылать данные (сканирование портов)
# -n Отключить DNS и поиск номеров портов по /etc/services
# -u Использовать для подключения UDP протокол

# https://www.digitalocean.com/community/tutorials/how-to-configure-the-linux-firewall-for-docker-swarm-on-ubuntu-16-04

ssh root@$1 "
netcat -vnz $2 22
netcat -vnz $2 2376
netcat -vnz $2 2377
netcat -vnz $2 7946
netcat -vnzu $2 7946
netcat -vnzu $2 4789
" #netcat -vnz $2 8786"