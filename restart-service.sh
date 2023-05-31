#!/bin/bash
export statLogPath=./statlog
export statAdminEmail=johnathanhartsfield@gmail.com
trap -- '' SIGTERM
git pull
go build -o statNotify
pkill -f statNotify
nohup ./statNotify > /dev/null & disown
sleep 2
