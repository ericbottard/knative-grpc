.PHONY: deploy build-image push-image run-server-local clean

build-image: .docker-image

.docker-image: cmd/server/main.go repeater/repeater.pb.go
	test -n "$(IMAGE)"  # $$IMAGE
	docker build . -t $(IMAGE) && touch .docker-image

push-image: build-image
	test -n "$(IMAGE)"  # $$IMAGE
	docker push $(IMAGE)

run-server-local: build-image
	docker run -p8080:8080 $(IMAGE)

deploy: push-image service.yaml
	test -n "$(IMAGE)"  # $$IMAGE
	sed "s@github.com/ericbottard/knative-grpc@$(IMAGE)@" service.yaml | kubectl apply -f -

repeater/repeater.pb.go: repeater.proto
	protoc repeater.proto --go_out=plugins=grpc:repeater

client: cmd/client/main.go repeater/repeater.pb.go
	GO111MODULE=on go build ./cmd/client

server: cmd/server/main.go repeater/repeater.pb.go
	GO111MODULE=on go build ./cmd/server

clean:
	rm -f server
	rm -f client
	rm .docker-image