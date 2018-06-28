FROM registry.cn-hangzhou.aliyuncs.com/fleacloud/cdnimage

RUN apk add docker
ADD build/cdn /usr/local/bin/cdn
ADD conf /etc/cdn

ADD supervisor.d /etc/supervisor.d
VOLUME ["/opt/cdn"]
VOLUME ["/var/run/docker"]

ENV DOCKER_HOST unix:///var/run/docker.sock
ENV MACVLAN_PARENT eth0
EXPOSE 8000
CMD ["supervisord", "-n"]