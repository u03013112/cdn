#!/bin/sh

yum install -y docker-ce
curl -o get-pip.py https://bootstrap.pypa.io/get-pip.py
python get-pip.py
pip install supervisor
mkdir -p /etc/supervisor/conf.d
cp -f supervisord.conf /etc/supervisor/
cp -f cdnapi.conf /etc/supervisor/conf.d/