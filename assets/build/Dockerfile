FROM hackermmzz/judge-base


WORKDIR /root

COPY qt-env.tar /root/qt-env.tar


COPY entrypoint.sh /entrypoint.sh

RUN chmod +x /entrypoint.sh


ENTRYPOINT ["/entrypoint.sh"]