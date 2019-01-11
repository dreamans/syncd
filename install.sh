#!/bin/bash

syncd_repo="github.com/dreamans/syncd"
build_repo="github.com/dreamans/syncd/syncd"

install_path=$( cd `dirname $0`; pwd )/syncd-deploy

if [ -d ${install_path} ];then
    install_path=${install_path}-$( date +%Y%m%d%H%M%S )
fi

build_syncd() {
    go get ${build_repo}
    cd $GOPATH/src/${build_repo}
    go run build.go
}

install_syncd() {
    mkdir ${install_path}
    cd ${install_path}
    mkdir bin log etc
    cp $GOPATH/src/${build_repo}/syncd ./bin/
    cp $GOPATH/src/${syncd_repo}/syncd.example.ini ./etc/syncd.ini
    cp -r $GOPATH/src/${syncd_repo}/public ./public
}


build_syncd

install_syncd

echo "Installing syncd binary: ${install_path}/bin"
echo "Installing web public: ${install_path}/public"
echo "Installing syncd.ini: ${install_path}/etc"
echo "Install complete"
