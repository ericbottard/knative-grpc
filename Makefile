.PHONY: deploy build-image

build-image: server/main.go repeater/repeater.pb.go
	test -n "$(IMAGE)"  # $$IMAGE
	docker build . -t $(IMAGE)
	docker push $(IMAGE)

deploy: build-image service.yaml
	test -n "$(IMAGE)"  # $$IMAGE
	sed "s@github.com/ericbottard/knative-grpc@$(IMAGE)@" service.yaml | kubectl apply -f -

repeater/repeater.pb.go: repeater.proto
	protoc repeater.proto --go_out=plugins=grpc:repeater

run-client: client/main.go repeater/repeater.pb.go
	GO111MODULE=on go build -o run-client ./client
	./run-client
