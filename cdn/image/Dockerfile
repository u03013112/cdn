From alpine

ADD Shanghai /etc/localtime
RUN echo "https://mirrors.aliyun.com/alpine/latest-stable/main" > /etc/apk/repositories
RUN echo "https://mirrors.aliyun.com/alpine/latest-stable/community" >> /etc/apk/repositories
#RUN echo "https://mirrors.ustc.edu.cn/alpine/latest-stable/main" > /etc/apk/repositories
#RUN echo "https://mirrors.ustc.edu.cn/alpine/latest-stable/community" >> /etc/apk/repositories

RUN apk update

RUN apk add tmux
RUN apk add tcpdump
RUN apk add supervisor
RUN apk add py2-pip

RUN mkdir -p /etc/supervisor.d


VOLUME ["/var/log/supervisor"]

CMD ["supervisord", "-n"]
