export GO111MODULE=on
export GOPROXY=https://goproxy.cn,direct

default: fmt carrier_svr

fmt:
	go fmt ./...

carrier_svr:
	go build -o bin/carrier_svr ./cmd/carrier_svr

demo:
	go build -o bin/demo ./cmd/demo
