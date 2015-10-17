#!/bin/bash

TARGET=(carrier/carrier_svr)

if [ $# -eq 1 ]
then
    TARGET=($1)
fi

CURDIR=`pwd`
export GOPATH="$CURDIR"

gofmt -w src
if [ $? != 0 ]; then
    echo "gofmt error"
    exit 1
fi

for t in ${TARGET[@]}; do
    go install $t
    if [ $? != 0 ]; then
        echo -e "\e[31minstall \e[4m$t\e[0m\e[31m error\e[0m"
        exit 1
    fi
    echo -e "\e[32minstall \e[4m$t\e[0m\e[32m success\e[0m"
done

echo -e "\e[36mdone\e[0m"
