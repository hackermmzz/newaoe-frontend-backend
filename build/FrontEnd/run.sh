docker run -d \
-p 5119:80 \
-v ./nginx.conf:/etc/nginx/nginx.conf \
-v ./dist:/usr/share/nginx/html/ \
dhi.io/nginx:1-debian13-dev