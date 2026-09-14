FROM hackermmzz/judge-base


RUN apk add --no-cache \
    bash


WORKDIR /root


COPY qt-env.tar /root/qt-env.tar


RUN cd /root/newaoe-frontend-backend && \
    git pull || true


RUN cd /root/new-aoe && \
    git pull || true


RUN cd /root/newaoe-frontend-backend && \
    pip3 install -r requirement.txt || true


COPY entrypoint.sh /entrypoint.sh

RUN chmod +x /entrypoint.sh


ENTRYPOINT ["/entrypoint.sh"]