docker run -d \
  -p 6379:6379 \
  -v ./conf:/etc/redis/redis.conf \
  redis:7-alpine \
  redis-server /etc/redis/redis.conf