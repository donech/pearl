#deployment for envory
ENVOY:=envoy-pearl

foundEnvoy := $(shell docker ps -a -f "name=$(ENVOY)" -q | grep -q . && echo Found || echo Not Found)
.PHONY: envoy
ifeq ($(foundEnvoy), Found)
envoy: stop-envoy start-envoy
	echo "stop and start envoy success"
else
envoy: start-envoy
	echo "start envoy success"
endif

.PHONY: stop-envoy
ifeq ($(foundEnvoy), Found)
stop-envoy:
	docker stop $(ENVOY)
else
stop-envoy:
	echo "no envory running"
endif

.PHONY: start-envoy
start-envoy:
	docker run -itd --rm --name $(ENVOY) \
    	  -p 51051:51051 \
    	  -p 10000:10000 \
          -v "$(shell pwd)/internal/proto/gen/descriptor:/data/descriptor" \
          -v "$(shell pwd)/deployment/envoy/envoy-config-dev.yml:/etc/envoy/envoy.yaml" \
          envoyproxy/envoy

.PHONY: in-envory
in-envoy:
	docker exec -it $(ENVOY) /bin/bash

.PHONY: valid-envoy
valid-envoy:
	docker run --rm \
              -v "$(shell pwd)/internal/proto/gen/descriptor:/data/descriptor" \
              -v "$(shell pwd)/envoy/envoy-config-dev.yml:/etc/envoy/envoy.yaml" \
              envoyproxy/envoy \
              --mode validate -c /etc/envoy/envoy.yaml
