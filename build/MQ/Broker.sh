docker run -d \
--name rmqbroker \
--network rocketmq \
-p 10912:10912 -p 10911:10911 -p 10909:10909 \
-p 8080:8080 -p 8081:8081 \
-e "NAMESRV_ADDR=rmqnamesrv:9876" \
-v ./broker.conf:/home/rocketmq/rocketmq-5.4.0/conf/broker.conf \
apache/rocketmq:5.4.0 sh mqbroker --enable-proxy \
-c /home/rocketmq/rocketmq-5.4.0/conf/broker.conf
