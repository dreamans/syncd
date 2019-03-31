#!/bin/bash

syncd_remote_repo=https://github.com/dreamans/syncd.git
syncd_repo_path=$HOME/.syncd-repo
syncd_install_path=$( cd `dirname $0`; pwd )/syncd-deploy

checkCommand() {
    type $1 > /dev/null 2>&1
    if [ $? -ne 0 ]; then
        echo "error: $1 must be installed"
        echo "install exit"
        exit 1
    fi
}

checkCommand "go"
checkCommand "git"
checkCommand "make"

if [ -d ${syncd_install_path} ];then
    syncd_install_path=${syncd_install_path}-$( date +%Y%m%d%H%M%S )
fi

rm -fr ${syncd_repo_path}
git clone ${syncd_remote_repo} ${syncd_repo_path}

cd ${syncd_repo_path}
make

rm -fr ${syncd_install_path}
cp -r ${syncd_repo_path}/output ${syncd_install_path}

rm -fr ${syncd_repo_path}

cat << EOF

Installing Syncd Path:  ${syncd_install_path}
Install complete.

EOF