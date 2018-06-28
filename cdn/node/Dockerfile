FROM registry.cn-hangzhou.aliyuncs.com/fleacloud/cdnimage

RUN apk update
RUN apk add nginx
#RUN apk add --upgrade nginx
#RUN apk add nginx-mod-stream
RUN apk add openssl-dev


ADD supervisor.d /etc/supervisor.d
ADD nginx.conf /etc/nginx/nginx.conf
ADD nginx /usr/sbin

RUN mkdir -p /run
VOLUME ["/etc/nginx/conf.d"]

CMD ["supervisord","-n"]
