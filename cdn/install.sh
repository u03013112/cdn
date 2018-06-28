#!/bin/sh
yum install -y yum-utils
yum-config-manager --add-repo https://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo
yum install -y docker-ce
systemctl enable docker
systemctl start docker

docker pull registry.cn-hangzhou.aliyuncs.com/fleacloud/cdnimage
docker pull registry.cn-hangzhou.aliyuncs.com/fleacloud/cdnnode
docker pull registry.cn-hangzhou.aliyuncs.com/fleacloud/cdnapi
echo "docker exec cdnapi cdn cdnctl $@" > /usr/local/bin/cdnctl
chmod a+x /usr/local/bin/cdnctl

curl -L https://github.com/docker/compose/releases/download/1.17.0/docker-compose-`uname -s`-`uname -m` -o /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose
# docker run -d -e "MACVLAN_PARENT=eth0" -v /var/run:/var/run -v /opt/cdn:/opt/cdn -p 8000:8000 --restart=always --name=cdnapi registry.cn-hangzhou.aliyuncs.com/fleacloud/cdnapi
# docker network create -d macvlan --subnet=172.18.0.0/24 --gateway=172.18.0.1 -o parent=enp10s0 cdnbr0
