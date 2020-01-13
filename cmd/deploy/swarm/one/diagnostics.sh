#!/bin/sh
echo ""
echo "  -------------------------------------"
echo "  ---$2 check ports open---"
echo "  -------------------------------------"
echo ""
trap 'echo " stop" ' INT TERM
#chmod +x droplet_diagnostics.sh && ./droplet_diagnostics.sh 1

main=$1   # the IP of this node
tested=$2 # the IP of tested node

if [ -z "${tested}" ]; then
    chmod +x ./../print_error.sh && ./../print_error.sh \
        "$0 <ip of any node> <ip of tested node>"
    exit 1
fi

# -v Вывод информации о процессе работы (verbose)
# -z Не посылать данные (сканирование портов)
# -n Отключить DNS и поиск номеров портов по /etc/services
# -u Использовать для подключения UDP протокол

# https://www.digitalocean.com/community/tutorials/how-to-configure-the-linux-firewall-for-docker-swarm-on-ubuntu-16-04

ssh root@$main "
netcat -vnz $tested 22
netcat -vnz $tested 2376
netcat -vnz $tested 2377
netcat -vnz $tested 7946
netcat -vnzu $tested 7946
netcat -vnzu $tested 4789
" #netcat -vnz $tested 8786"