PROJECT:=pearl
CONFIG:=app.yaml
APP:=grpc
include deployment/envoy/Makefile

.PHONY: run
run: lint
	@echo "wellcome for donech land, this will run the erp system for you."
	@go run main.go $(APP) -c $(CONFIG)

.PHONY: gen
gen: stop-envoy
	go generate ./...
	rm -Rf internal/proto/gen/*
	docker run --rm -v "$(shell pwd)/internal/proto:/work" uber/prototool prototool all

.PHONY: lint
lint:
	docker run --rm -v "$(shell pwd)/internal/proto:/work" uber/prototool prototool lint
	golangci-lint run
