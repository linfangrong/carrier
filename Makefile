export GOPATH=$(shell pwd)

default: clean fmt install

clean:
	@echo "clean..."
	rm -rf pkg bin

fmt:
	@echo "format..."
	gofmt -w src

install:
	@echo "install..."
	go install carrier/carrier_svr
