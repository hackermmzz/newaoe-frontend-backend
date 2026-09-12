docker network create rocketmq 
docker run -d --name rmqnamesrv -p 9876:9876 --network rocketmq apache/rocketmq:5.4.0 sh mqnamesrv