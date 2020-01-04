#!/bin/bash
echo ""
echo "  ----------------------------"
echo "  -------------test-----------"
echo "  ----------------------------"
echo ""

loadAddrs() {
    local addrs=""
    input=`cat $1 | grep "^[^#]"| grep "addr"`
    set -- $input
    while [ $1 ]
        do
        addrs=$addrs' '$3
        shift 3
    done
    echo $addrs
}

loadNames() {
    local names=""
    input=`cat $1 | grep "^[^#]"| grep "name"`
    set -- $input
    while [ $1 ]
        do
        names=$names' '$3
        shift 3
    done
    echo $names
}

loadPath() {
    input=`cat $1 | grep "^[^#]"| grep "path"`
    set -- $input
    echo $3
}

argsForAdd() {
    # addrs=$(echo $1 | tr ' ' '\n')
    # names=$(echo $2 | tr ' ' '\n')
    path=$3

    IFS=' ' read -ra addrs <<< "$1"
    IFS=' ' read -ra names <<< "$2"
    local args=""
    for i in ${!addrs[@]}
    do
        args=$args' '${addrs[i]}' '$path' '${names[i]}
    done
    echo $args
}

addrs=$(loadAddrs "info.txt")
names=$(loadNames "info.txt")
path=$(loadPath "info.txt")
args=$(argsForAdd "$addrs" "$names" $path)
echo 'args:'$args
#./all/add_exp.sh $args
#./all/firewall_exp.sh $addrs
#./all/diagnostics_exp.sh $addrs
#./all/metrics_exp.sh $addrs
#./one/manager.sh $addr1  #dont forget to set manager token to worker.sh 
#./all/worker.sh $addr1 $addr2 $addr3 $addr4 $addr5
./all/roles.sh $addrs
