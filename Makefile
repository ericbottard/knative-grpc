.PHONY: deploy build-image

build-image: server/main.go repeater/repeater.pb.go
	docker build . -t $(IMAGE)
	docker push $(IMAGE)

deploy: build-image service.yaml
	sed "s@github.com/ericbottard/knative-grpc@$(IMAGE)@" service.yaml | kubectl apply -f -

repeater/repeater.pb.go: repeater.proto
	protoc repeater.proto --go_out=plugins=grpc:repeater
