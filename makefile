# Check to see if we can use ash, in Alpine images, or default to BASH.
SHELL_PATH = /bin/ash
SHELL = $(if $(wildcard $(SHELL_PATH)),/bin/ash,/bin/bash)

# Deploy First Mentality


# ==============================================================================
# Define dependencies

GOLANG          := golang:1.26
ALPINE          := alpine:3.23
KIND            := kindest/node:v1.35.0
POSTGRES        := postgres:18.3
GRAFANA         := grafana/grafana:12.4.0
PROMETHEUS      := prom/prometheus:v3.10.0
TEMPO           := grafana/tempo:2.10.0
LOKI            := grafana/loki:3.6.0
PROMTAIL        := grafana/promtail:3.6.0

KIND_CLUSTER    := ahmed-starter-cluster
NAMESPACE       := sales-system
SALES_APP       := sales
AUTH_APP        := auth
BASE_IMAGE_NAME := ahmedhrid
VERSION         := 0.0.1
SALES_IMAGE     := $(BASE_IMAGE_NAME)/$(SALES_APP):$(VERSION)
METRICS_IMAGE   := $(BASE_IMAGE_NAME)/metrics:$(VERSION)
AUTH_IMAGE      := $(BASE_IMAGE_NAME)/$(AUTH_APP):$(VERSION)

# VERSION       := "0.0.1-$(shell git rev-parse --short HEAD)"

# ==============================================================================
# Detect operating system and set the appropriate open command

UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Darwin)
	OPEN_CMD := open
else
	OPEN_CMD := xdg-open
endif

# ==============================================================================
# Install dependencies

dev-gotooling:
	go install github.com/divan/expvarmon@latest
	go install github.com/rakyll/hey@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install golang.org/x/vuln/cmd/govulncheck@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

dev-brew:
	brew update
	brew list kind || brew install kind
	brew list kubectl || brew install kubectl
	brew list kustomize || brew install kustomize
	brew list pgcli || brew install pgcli
	brew list watch || brew install watch
	brew list protobuf || brew install protobuf
	brew list grpcurl || brew install grpcurl

dev-docker:
	docker pull docker.io/$(GOLANG) & \
	docker pull docker.io/$(ALPINE) & \
	docker pull docker.io/$(KIND) & \
	docker pull docker.io/$(POSTGRES) & \
	docker pull docker.io/$(GRAFANA) & \
	docker pull docker.io/$(PROMETHEUS) & \
	docker pull docker.io/$(TEMPO) & \
	docker pull docker.io/$(LOKI) & \
	docker pull docker.io/$(PROMTAIL) & \
	wait;

# ==============================================================================
# Building containers

build: sales auth

sales:
	docker build \
		-f zarf/docker/dockerfile.sales \
		-t $(SALES_IMAGE) \
		--build-arg BUILD_TAG=$(VERSION) \
		--build-arg BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
		.

auth:
	docker build \
		-f zarf/docker/dockerfile.auth \
		-t $(AUTH_IMAGE) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
		.
# ==============================================================================
# Running from within k8s/kind
# Docker Desktop 28.3.2 changed how it stores image layers, causing KIND's kind
# load docker-image command to fail with "content digest not found" errors. The
# workaround uses docker save | docker exec to bypass this incompatibility for
# the critical images allowing this to work without a network.

dev-up:
	kind create cluster \
		--image $(KIND) \
		--name $(KIND_CLUSTER) \
		--config zarf/k8s/dev/kind-config.yaml

	kubectl wait --timeout=120s --namespace=local-path-storage --for=condition=Available deployment/local-path-provisioner

# 	docker save $(POSTGRES) | docker exec -i $(KIND_CLUSTER)-control-plane ctr --namespace=k8s.io images import - & \
# 	docker save $(GRAFANA) | docker exec -i $(KIND_CLUSTER)-control-plane ctr --namespace=k8s.io images import - & \
# 	docker save $(PROMETHEUS) | docker exec -i $(KIND_CLUSTER)-control-plane ctr --namespace=k8s.io images import - & \
# 	docker save $(TEMPO) | docker exec -i $(KIND_CLUSTER)-control-plane ctr --namespace=k8s.io images import - & \
# 	docker save $(LOKI) | docker exec -i $(KIND_CLUSTER)-control-plane ctr --namespace=k8s.io images import - & \
# 	docker save $(PROMTAIL) | docker exec -i $(KIND_CLUSTER)-control-plane ctr --namespace=k8s.io images import - & \
# 	wait;

dev-down:
	kind delete cluster --name $(KIND_CLUSTER)

dev-up-registry:
	KIND_CLUSTER_NAME=$(KIND_CLUSTER) ./zarf/k8s/dev/kind-with-registry.sh

dev-down-registry:
	kind delete cluster --name $(KIND_CLUSTER)
	docker stop kind-registry 2>/dev/null || true
	docker rm kind-registry 2>/dev/null || true

dev-status-all:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --watch --all-namespaces

dev-status:
	watch -n 2 kubectl get pods -o wide --all-namespaces

# ------------------------------------------------------------------------------

dev-load:
	kind load docker-image $(SALES_IMAGE) --name $(KIND_CLUSTER)
	kind load docker-image $(AUTH_IMAGE) --name $(KIND_CLUSTER)

dev-apply:
	kustomize build zarf/k8s/dev/auth | kubectl apply -f -
	kubectl wait pods --namespace=$(NAMESPACE) --selector app=$(AUTH_APP) --timeout=120s --for=condition=Ready

	kustomize build zarf/k8s/dev/sales | kubectl apply -f -
	kubectl wait pods --namespace=$(NAMESPACE) --selector app=$(SALES_APP) --timeout=120s --for=condition=Ready

dev-restart:
	kubectl rollout restart deployment $(AUTH_APP) --namespace=$(NAMESPACE)
	kubectl rollout restart deployment $(SALES_APP) --namespace=$(NAMESPACE)

dev-run: build dev-up dev-load dev-apply

dev-update: build dev-load dev-restart

dev-update-apply: build dev-load dev-apply

dev-logs:
	kubectl logs --namespace=$(NAMESPACE) -l app=$(SALES_APP) --all-containers=true -f --tail=100 --max-log-requests=6 | go run apis/tooling/logfmt/main.go -service=$(SALES_APP)

dev-logs-auth:
	kubectl logs --namespace=$(NAMESPACE) -l app=$(AUTH_APP) --all-containers=true -f --tail=100 | go run apis/tooling/logfmt/main.go



# ------------------------------------------------------------------------------

dev-logs-init:
	kubectl logs --namespace=$(NAMESPACE) -l app=$(SALES_APP) -f --tail=100 -c init-migrate-seed

dev-describe-node:
	kubectl describe node

dev-describe-deployment:
	kubectl describe deployment --namespace=$(NAMESPACE) $(SALES_APP)

dev-describe-sales:
	kubectl describe pod --namespace=$(NAMESPACE) -l app=$(SALES_APP)

dev-describe-auth:
	kubectl describe pod --namespace=$(NAMESPACE) -l app=$(AUTH_APP)


# ==============================================================================
# Metrics and Tracing

metrics-view-sc:
	expvarmon -ports="localhost:3010" -vars="build,requests,goroutines,errors,panics,mem:memstats.HeapAlloc,mem:memstats.HeapSys,mem:memstats.Sys"

metrics-view:
	expvarmon -ports="localhost:4020" -endpoint="/metrics" -vars="build,requests,goroutines,errors,panics,mem:memstats.HeapAlloc,mem:memstats.HeapSys,mem:memstats.Sys"

grafana:
	$(OPEN_CMD) http://localhost:3100/

statsviz:
	$(OPEN_CMD) http://localhost:3010/debug/statsviz

# ==============================================================================

# 

run:
	go run apis/services/sales/main.go | go run apis/tooling/logfmt/main.go

tidy:
	go mod tidy
	go mod vendor

help:
	go run apis/services/sales/main.go


# 	=======================================
admin:
	go run apis/tooling/admin/main.go

# admin token
# export TOKEN=eyJhbGciOiJSUzI1NiIsImtpZCI6IjU0YmIyMTY1LTcxZTEtNDFhNi1hZjNlLTdkYTRhMGUxZTJjMSIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzZXJ2aWNlIHByb2plY3QiLCJzdWIiOiIzOGRjOWQ4NC0wMThiLTRhMTUtYjk1OC0wYjc4YWYxMWMzMDEiLCJleHAiOjE3NDU1OTY1MTEsImlhdCI6MTcxNDA2MDUxMSwiUm9sZXMiOlsiQURNSU4iXX0.nioy4PpggnfrwTxNTQKbviCs3duF53Q5jcoRQqdngQSv7lccKYgTmxzuyMano-Yd-KijtHBCZxWAOFEv5w6xGCfqmQRThKXzQXiHN5Cv0OGab5lmThPGRuCHv35lEQzImKU9E1skSwHvCwyX89pRJpnku9PKJMT_Z4oT6amwFA50HU8jSM8j1HQ0ao60jSMgELKFFb9m3u4ZKIj4w7qxwwV9JD2_wH8HWjt5P1L2V5YtnP9vgMOBZ617TTGRysDS8WXQGXqEAiVQVZSneJCUR-4HofXvOTYIKyQG3iUAs3WTf91EubJQbeW6cFmwudE4xx4t20EaVIUkMr00jFxHtg

# user token
# export TOKEN=eyJhbGciOiJSUzI1NiIsImtpZCI6IjU0YmIyMTY1LTcxZTEtNDFhNi1hZjNlLTdkYTRhMGUxZTJjMSIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzZXJ2aWNlIHByb2plY3QiLCJzdWIiOiIzOGRjOWQ4NC0wMThiLTRhMTUtYjk1OC0wYjc4YWYxMWMzMDEiLCJleHAiOjE3NDU1OTY1NTUsImlhdCI6MTcxNDA2MDU1NSwiUm9sZXMiOlsiVVNFUiJdfQ.bBomVD8igAkiXzKvsBQGj5Nb5ho2MYq-eOXeFjtfdauayE4lQamFtWHUjKq1KaIzMvxlybdU_37py620vOBrQ7tUIYTdY91ggf-AgakBTm3UTs494BjIv3rvPM0baXKSMSe8Ao2acjhTSPA9Crz0mEgWynd3gcuSXVAXrSO3eH8bPuBqR2Ohp1_JsE4Y_dj0Mfmzo8wRQHmBI-QVMT0udUlxVGsLhurGmFpLSrQ7Vzr-xuCJPMgjkTZmSvelCrxSeql6-scwZLTHjgqkzqIMS5EceyKAQfYuuRgqrwwIAtGhyM6SrPHH3WFEy2RHW2ebqQKa-fc3JAcJH44hgiiGsA

curl-auth:
	curl -il \
	-H "Authorization: Bearer ${TOKEN}" "http://localhost:3000/testauth"

token:
	curl -il \
	--user "admin@example.com:gophers" http://localhost:6000/auth/token/54bb2165-71e1-41a6-af3e-7da4a0e1e2c1

curl-auth2:
	curl -il \
	-H "Authorization: Bearer ${TOKEN}" "http://localhost:6000/auth/authenticate"
