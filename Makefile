# Digital Platform operator makefile
#
# The following env variables must be set to run make targets
#    SERVICE_NAME : Name of the service in cmd package e.g product-operator
#    DOCKER_REGISTRY_URL: Docker Registry URL where to push image
# The following env variable is optional and only required when pushing to Docker Hub
#    DOCKER_PUSH_LIBRARY : Registry library of the image repository

.PHONY: clean fmt-check test-local test code-analysis lint vet coverage-local documentation swagger build build-for-linux build-for-windows build-docker-image build-image publish push-image push-docker-image

ifdef DOCKER_PUSH_LIBRARY
DOCKER_REGISTRY_URL := $(DOCKER_REGISTRY_URL)/$(DOCKER_PUSH_LIBRARY)
endif

VERSION           := 0.0.7
SERVICE_NAME      := $(SERVICE_NAME)

ifdef BUILD_NUMBER
VERSION           := $(VERSION)-$(BUILD_NUMBER)
endif

IMAGE_NAME        := $(DOCKER_REGISTRY_URL)/$(SERVICE_NAME):$(VERSION)
IMAGE_NAME_TEST   := $(DOCKER_REGISTRY_URL)/$(SERVICE_NAME)-test:$(VERSION)
IMAGE_NAME_DOCS   := $(DOCKER_REGISTRY_URL)/$(SERVICE_NAME)-docs:$(VERSION)

# PKG_LIST is the list of packages for code coverage with exclusion of generated code or integration test
BUILD_TIME   := $(shell date +%Y-%m-%d-%H:%M)
GIT_DIRTY    = `git diff-index --quiet HEAD -- || echo "✗-"`
GIT_REVISION = `git rev-parse --short HEAD`
GIT_BRANCH   = `git rev-parse --symbolic-full-name --abbrev-ref HEAD`
LDFLAGS      = "-s \
	-X eurocontrol.io/digital-platform-product-deployment/pkg/actuator.buildTime=${BUILD_TIME} \
	-X eurocontrol.io/digital-platform-product-deployment/pkg/actuator.gitRevision=${GIT_DIRTY}${GIT_REVISION} \
	-X eurocontrol.io/digital-platform-product-deployment/pkg/actuator.gitBranch=${GIT_BRANCH} \
	-X eurocontrol.io/digital-platform-product-deployment/pkg/actuator.version=${VERSION} \
	-X eurocontrol.io/digital-platform-product-deployment/pkg/actuator.name=${SERVICE_NAME}"

COVER_PROFILE          := dist/test-results/coverage.out
COVERAGE_REPORT_HTML   := dist/test-results/coverage.html
COVERAGE_REPORT_XML    := dist/test-results/coverage.xml
COVERAGE_REPORT_PLAIN  := dist/test-results/coverage.txt
COVERAGE_REPORT_JSON   := dist/test-results/coverage.json
TEST_OUTPUT            := dist/test-results/test.out
TEST_REPORT_XML        := dist/test-results/test-report.xml

# Use advanced caching features for multistage docker builds
export DOCKER_BUILDKIT := 1


# Reinitialize directories
clean: cleanup-docker-images
	@rm -rf bin dist
	@mkdir -p bin dist/{artifacts,code-analysis,documentation,test-results}


fmt-check:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"


# Unit tests for local usage
test-local: clean
	go test -v -parallel 1 -count 1 ./pkg/... ./cmd/... -coverprofile "$(COVER_PROFILE)" -coverpkg $(shell go list ./... | grep -v -e pb -e generated -e fake -e test -e zz_ -e hack -e docs -e certificate -e registry -e deployment/operations | tr '\n' ',')


# Unit tests
test: .create-docker-image-for-testing
	go test -v -parallel 1 -p 1 -count 1 ./pkg/... ./cmd/... -coverprofile "$(COVER_PROFILE)" -coverpkg $(shell go list ./... | grep -v -e pb -e generated -e fake -e test -e zz_ -e hack -e docs -e certificate -e registry -e deployment/operations | tr '\n' ',') 2>&1 > $(TEST_OUTPUT)
	@docker run --rm \
      --user="$(id -u):$(id -g)" \
	  -v "$(PWD)/dist:/src/dist" \
	  -e TEST_OUTPUT=$(TEST_OUTPUT) \
	  -e TEST_REPORT_XML=$(TEST_REPORT_XML) \
	  $(IMAGE_NAME_TEST) /bin/bash -c "/junittest.sh"


# Code-analysis
code-analysis: lint vet fmt-check


# Check code style
lint:
	@docker run --rm -v "$(PWD)/dist:/src/dist" $(IMAGE_NAME_TEST) golint ./pkg/... ./cmd/...


# Check code correctness
vet:
	@docker run --rm -v "$(PWD)/dist:/src/dist" $(IMAGE_NAME_TEST) go vet -v -structtag=false ./cmd/... ./pkg/...


# Generate code coverage reports for local usage
coverage-local:
	go tool cover -html $(COVER_PROFILE) -o $(COVERAGE_REPORT_HTML)
	gocov convert $(COVER_PROFILE) > $(COVERAGE_REPORT_JSON)
	cat $(COVERAGE_REPORT_JSON) | gocov-xml > $(COVERAGE_REPORT_XML)
	cat $(COVERAGE_REPORT_JSON) | gocov report > $(COVERAGE_REPORT_PLAIN)
	cat $(COVERAGE_REPORT_PLAIN)


# Project Documentation
documentation: .create-docker-image-for-documentation swagger


# Generate swagger API documentation
swagger:
	@docker run --rm \
	  --security-opt label=disable \
	  -v "$(PWD)/dist:/tmp" \
	  $(IMAGE_NAME_DOCS) swagger generate spec -w ./cmd/$(SERVICE_NAME)/ -m -o /tmp/documentation/swagger.yaml


# Local build, exe will be in bin directory
build:
	rm -rf bin/$(SERVICE_NAME)
	go build -ldflags $(LDFLAGS) -v -o bin/$(SERVICE_NAME) ./cmd/$(SERVICE_NAME)


# Build linux cross compilation
build-for-linux:
	rm -rf dist/artifacts/linux/$(SERVICE_NAME)
	GOOS=linux GOARCH=amd64 go build -ldflags $(LDFLAGS) -v -o dist/artifacts/linux/$(SERVICE_NAME)/app ./cmd/$(SERVICE_NAME)


# Build windows cross compilation
build-for-windows:
	rm -rf dist/artifacts/windows/$(SERVICE_NAME)
	GOOS=windows GOARCH=386 go build -ldflags $(LDFLAGS) -v -o dist/artifacts/windows/$(SERVICE_NAME).exe ./cmd/$(SERVICE_NAME)


# Build docker image based on linux
build-docker-image:
	@docker build --rm --build-arg SERVICE_NAME=$(SERVICE_NAME) . -t $(IMAGE_NAME)


# Build local image for linux and copy in docker image
build-image:
	@docker build --rm --build-arg DOCKER_REGISTRY_URL=$(DOCKER_REGISTRY_URL) --build-arg SERVICE_NAME=$(SERVICE_NAME) --build-arg LDFLAGS=$(LDFLAGS) . -t $(IMAGE_NAME) -f Dockerfile.build --target release


# Publish docker image to registry
publish:
	@echo 'Publishing $(IMAGE_NAME)'
	@docker push $(IMAGE_NAME)
	@docker rmi $(IMAGE_NAME)


# Push image to docker registry after build-image or build-docker-image
push-docker-image:
	@echo 'Publishing $(IMAGE_NAME)'
	@docker push $(IMAGE_NAME)
	@docker tag $(IMAGE_NAME) $(DOCKER_REGISTRY_URL)/$(SERVICE_NAME):latest
	@docker push $(DOCKER_REGISTRY_URL)/$(SERVICE_NAME):latest


# Build and push image to docker registry
push-image: build-image
	@echo 'Publishing $(IMAGE_NAME)'
	@docker push $(IMAGE_NAME)
	@docker tag $(IMAGE_NAME) $(DOCKER_REGISTRY_URL)/$(SERVICE_NAME):latest
	@docker push $(DOCKER_REGISTRY_URL)/$(SERVICE_NAME):latest


cleanup-docker-images: .cleanup-docker-image-for-testing .cleanup-docker-image-for-documentation
	-@docker rmi -f $(docker images -f "dangling=true" -q)


.cleanup-docker-image-for-testing:
	-@docker rmi -f $(IMAGE_NAME_TEST)

.create-docker-image-for-testing:
	@docker build --rm . -t $(IMAGE_NAME_TEST) -f Dockerfile.test --target test

.cleanup-docker-image-for-documentation:
	-@docker rmi -f $(IMAGE_NAME_DOCS)

.create-docker-image-for-documentation:
	@docker build --rm . -t $(IMAGE_NAME_DOCS) -f Dockerfile.docs --target documentation


.check-env:
ifndef SERVICE_NAME
$(error SERVICE_NAME is undefined)
endif
ifndef DOCKER_REGISTRY_URL
$(error DOCKER_REGISTRY_URL is undefined)
endif