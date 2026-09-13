#!/bin/bash

# ==========================
# 第一步：当前目录 git pull
# ==========================

MAX_RETRY=10
COUNT=1

while [ $COUNT -le $MAX_RETRY ]
do
    echo "当前目录 git pull，第 $COUNT/$MAX_RETRY 次..."

    if git pull; then
        echo "当前目录 git pull 成功"
        break
    fi

    echo "当前目录 git pull 失败，5秒后重试..."
    COUNT=$((COUNT+1))
    sleep 5
done

if [ $COUNT -gt $MAX_RETRY ]; then
    echo "当前目录 git pull 失败10次，退出容器"
    exit 1
fi


# ==========================
# 第二步：进入 new-aoe
# ==========================

cd new-aoe/ || exit 1


# ==========================
# 第三步：new-aoe git pull
# ==========================

COUNT=1

while [ $COUNT -le $MAX_RETRY ]
do
    echo "new-aoe git pull，第 $COUNT/$MAX_RETRY 次..."

    if git pull; then
        echo "new-aoe git pull 成功"
        break
    fi

    echo "new-aoe git pull 失败，5秒后重试..."
    COUNT=$((COUNT+1))
    sleep 5
done

if [ $COUNT -gt $MAX_RETRY ]; then
    echo "new-aoe git pull 失败10次，退出容器"
    exit 1
fi


# ==========================
# 第四步：启动 Python
# ==========================

cd ..

echo "启动 Python 服务..."
exec python3 main.py