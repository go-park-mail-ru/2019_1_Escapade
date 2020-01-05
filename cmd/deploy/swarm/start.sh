#!/bin/bash
echo ""
echo "  -------------------"
echo "  -------Start-------"
echo "  -------------------"
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

#./all/add.sh $args
#./all/firewall.sh $addrs
#./all/diagnostics.sh $addrs
#./all/metrics.sh $addrs
./all/roles.sh $addrs
