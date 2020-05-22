
all: help


help:
	@echo ''
	@echo 'Here are the targets:'
	@echo ''
	@echo 'To test                :    "make check"'
	@echo 'To develop             :    "make develop"'
	@echo 'To install             :    "make install"'
	@echo 'To format              :    "make format"'
	@echo 'To lint                :    "make lint"'
	@echo ''



install:
	go get github.com/ianamason/yices2_go_bindings/cmd/yices_info

develop:
	go install github.com/ianamason/yices2_go_bindings/cmd/yices_info

# make sure we do not run tests in parallel!
check: develop
	go test -v -p 1 ./test_api

format:
	gofmt -s -w yices_api/*.go yices/*.go cmd/*/*.go test_api/*.go

lint:
	golint ./test_api/  ./yices/ ./yices_api/ ./cmd/...


ineffassign:
	ineffassign ./test_api/  ./yices/ ./yices_api/ ./cmd/
