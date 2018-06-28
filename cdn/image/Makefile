IMAGE_PREFIX=registry.cn-hangzhou.aliyuncs.com/fleacloud
IMAGE_NAME=$(IMAGE_PREFIX)/cdnimage

all:
	docker build -t $(IMAGE_NAME) .
clean:
	docker rmi $(IMAGE_NAME)
push:
	docker push $(IMAGE_NAME)