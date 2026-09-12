go build -o newaoe ./ && \
docker run -d \
  --name backend \
  -v $(pwd)/assets/config.yaml:/app/config.yaml \
  -v $(pwd)/assets/newaoe:/app/newaoe \
  -w /app \
  --network newaoe-net \
  ubuntu \
  ./newaoe