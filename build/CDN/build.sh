docker run -d  \
-p 9000:9000 -p 9001:9001 \
-v ./minio/:/data 
-e MINIO_ROOT_USER=newaoe \
-e MINIO_ROOT_PASSWORD=newaoe-newaoe \
minio/minio:latest server /data --console-address ":9001"           