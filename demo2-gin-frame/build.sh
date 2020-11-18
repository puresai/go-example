#!/bin/bash
SERVER="demo2-gin-frame"

function status() 
{
	if [ "`pgrep $SERVER -u $UID`" != "" ];then
		echo $SERVER is running
	else
		echo $SERVER is not running
	fi
}

function build() 
{
	echo "build..."

    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./$SERVER main.go
    if [ $? -ne "0" ];then  
        echo "built error!!!"
        return 
    fi  

    echo "built success!"
}


case "$1" in
	'status')
	status
	;; 
    'build')
	build
	;;  
	*)  
	echo "unknown, please: $0 {status or build}"
	exit 1
	;;  
esac
