#!/bin/bash

set -e


echo "启动docker"


dockerd-entrypoint.sh &



echo "等待docker"


for i in {1..30}
do
    if docker ps -a >/dev/null 2>&1
    then
        break
    fi

    sleep 1
done



#################################
# 第一次加载qt-env
#################################

if ! docker image inspect hackermmzz/qt-env:latest >/dev/null 2>&1
then

    echo "加载qt-env镜像"

    docker load \
        -i /root/qt-env.tar

else

    echo "qt-env已经存在"

fi



rm -f /root/qt-env.tar



cd /root/newaoe-frontend-backend


python3 main.py