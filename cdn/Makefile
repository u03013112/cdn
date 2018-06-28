IMAGE_PREFIX=registry.cn-hangzhou.aliyuncs.com/fleacloud
IMAGE_NAME=$(IMAGE_PREFIX)/cdnapi

all:
	make -C image
	make -C node
	
cdnapi:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/cdn main/main.go
	docker build -t $(IMAGE_NAME) .
push-cdnapi:
	docker push $(IMAGE_NAME)
push:
	make push -C image
	make push -C node
	docker push $(IMAGE_NAME)
clean:
	make clean -C image 
	make clean -C node
	docker rmi  $(IMAGE_NAME)
	rm -rf build/*

