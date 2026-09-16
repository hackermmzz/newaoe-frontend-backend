#!/bin/bash

set -e


echo "======================"
echo "启动 Docker daemon"
echo "======================"


# 启动dind
dockerd-entrypoint.sh &


echo "等待docker启动..."


# 等docker稳定
for i in {1..30}
do
    if docker ps -a >/dev/null 2>&1
    then
        echo "Docker启动成功"
        break
    fi

    sleep 1
done


if ! docker ps -a >/dev/null 2>&1
then
    echo "Docker启动失败"
    exit 1
fi


echo "======================"
echo "pull qt-env"
echo "======================"

if docker image inspect hackermmzz/qt-env:latest >/dev/null 2>&1
then

    echo "qt-env 已存在，跳过加载"

else

    if [ -f /root/qt-env.tar ]
    then

        echo "第一次加载 qt-env 镜像"

        docker load \
            -i /root/qt-env.tar


        echo "qt-env 加载完成"

        # 删除tar释放空间
        rm -f /root/qt-env.tar

    else

        echo "未找到 qt-env.tar，尝试在线pull"

        docker pull hackermmzz/qt-env || true

    fi

fi


##################################
# git pull函数
##################################

git_retry()
{
    local dir=$1

    cd "$dir"

    for i in {1..3}
    do
        echo "git pull 第 $i 次"

        if git pull
        then
            echo "pull成功"
            return
        fi


        echo "失败，等待3秒"

        sleep 3

    done


    echo "git pull $dir 三次失败，程序退出"
    exit 1
}



##################################
# frontend-backend
##################################

echo "更新frontend-backend"


git_retry /root/newaoe-frontend-backend



##################################
# new-aoe
##################################

echo "更新new-aoe"


git_retry /root/newaoe-frontend-backend/new-aoe



##################################
# pip install
##################################

echo "安装python依赖"


cd /root/newaoe-frontend-backend


if ! pip3 install -r requirements.txt
then
    echo "安装python依赖失败，程序退出"
    exit 1
fi


##################################
# cp到/tmp分区
##################################

echo "cp到/tmp分区"
dir="/tmp/judge_$(date +%Y%m%d_%H%M%S)"
mkdir -p "$dir"
cp -r /root/newaoe-frontend-backend/* $dir/

##################################
# 切换到/tmp/judge目录下
##################################

echo "切换到/tmp/judge目录下"
cd "$dir"


##################################
# 启动程序
##################################

echo "启动main.py"


python3 main.py


exit $?