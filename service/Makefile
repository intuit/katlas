PACKAGE  = service
DATE    ?= $(shell date +%FT%T%z)
VERSION ?= $(shell git describe --tags --always --dirty --match=v* 2> /dev/null || \
			cat $(CURDIR)/.version 2> /dev/null || echo v0)
GOPATH   = $(CURDIR)/.gopath~
BIN      = $(GOPATH)/bin
BASE     = $(GOPATH)/src/github.com/intuit/katlas/$(PACKAGE)
FULLPATH = github.com/intuit/katlas/$(PACKAGE)
PKGS     = $(or $(PKG),$(shell cd $(BASE) && env GOPATH=$(GOPATH) $(GO) list ./... | grep -v "^$(PACKAGE)/vendor/" | grep -v "$(FULLPATH)/ext_service/idps"))
TESTPKGS = $(shell env GOPATH=$(GOPATH) $(GO) list -f '{{ if or .TestGoFiles .XTestGoFiles }}{{ .ImportPath }}{{ end }}' $(PKGS))

GO      = go
GODOC   = godoc
GOFMT   = gofmt
DEP     = dep
TIMEOUT = 15
V = 0
Q = $(if $(filter 1,$V),,@)
M = $(shell printf "\033[34;1m▶\033[0m")

.PHONY: all
all: fmt lint vendor | $(BASE) ; $(info $(M) building executable…) @ ## Build program binary
	@echo $(PKGS)
	@echo $(TESTPKGS)
	$Q cd $(BASE) && CGO_ENABLED=0 GOOS=linux $(GO) build \
		-tags release \
		-ldflags '-X $(PACKAGE)/cmd.Version=$(VERSION) -X $(PACKAGE)/cmd.BuildDate=$(DATE)' \
		-a -o bin/katlas server.go

$(BASE): ; $(info $(M) setting GOPATH…)
	mkdir -vp $(dir $@)
	mkdir -vp $(BIN)
	@ln -sf $(CURDIR) $@

# Tools

GOLINT = $(BIN)/golint
$(BIN)/golint: | $(BASE) ; $(info $(M) building golint…)
	$Q go get -u golang.org/x/lint/golint
GOSWAGGER = $(BIN)/swagger
$(BIN)/swagger: | $(BASE) ; $(info $(M) building swagger…)
	$Q go get github.com/go-swagger/go-swagger/cmd/swagger
DEP = $(BIN)/dep
$(BIN)/dep: | $(BASE) ; $(info $(M) building dep…)
	$Q go get github.com/golang/dep/cmd/dep

GOCOVMERGE = $(BIN)/gocovmerge
$(BIN)/gocovmerge: | $(BASE) ; $(info $(M) building gocovmerge…)
	$Q go get github.com/wadey/gocovmerge

GOCOV = $(BIN)/gocov
$(BIN)/gocov: | $(BASE) ; $(info $(M) building gocov…)
	$Q go get github.com/axw/gocov/...

GOCOVXML = $(BIN)/gocov-xml
$(BIN)/gocov-xml: | $(BASE) ; $(info $(M) building gocov-xml…)
	$Q go get github.com/AlekSi/gocov-xml

GO2XUNIT = $(BIN)/go2xunit
$(BIN)/go2xunit: | $(BASE) ; $(info $(M) building go2xunit…)
	$Q go get github.com/tebeka/go2xunit

# Tests

TEST_TARGETS := test-default test-bench test-short test-verbose test-race
.PHONY: $(TEST_TARGETS) test-xml check test tests
test-bench:   ARGS=-run=__absolutelynothing__ -bench=. ## Run benchmarks
test-short:   ARGS=-short        ## Run only short tests
test-verbose: ARGS=-v            ## Run tests in verbose mode with coverage reporting
test-race:    ARGS=-race         ## Run tests with race detector
$(TEST_TARGETS): NAME=$(MAKECMDGOALS:test-%=%)
$(TEST_TARGETS): test
check test tests: fmt lint vendor | $(BASE) ; $(info $(M) running $(NAME:%=% )tests…) @ ## Run tests
	$Q cd $(BASE) && $(GO) test -timeout $(TIMEOUT)s $(ARGS) $(TESTPKGS)

test-xml: fmt lint vendor | $(BASE) $(GO2XUNIT) ; $(info $(M) running $(NAME:%=% )tests…) @ ## Run tests with xUnit output
	$Q cd $(BASE) && 2>&1 $(GO) test -timeout 20s -v $(TESTPKGS) | tee test/tests.output
	$(GO2XUNIT) -fail -input test/tests.output -output test/tests.xml

COVERAGE_MODE = atomic
COVERAGE_PROFILE = $(COVERAGE_DIR)/profile.out
COVERAGE_XML = $(COVERAGE_DIR)/coverage.xml
COVERAGE_HTML = $(COVERAGE_DIR)/index.html
.PHONY: test-coverage test-coverage-tools
test-coverage-tools: | $(GOCOVMERGE) $(GOCOV) $(GOCOVXML)
test-coverage: COVERAGE_DIR := $(CURDIR)/test/coverage.$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
test-coverage: fmt lint vendor test-coverage-tools | $(BASE) ; $(info $(M) running coverage tests…) @ ## Run coverage tests
	$Q mkdir -p $(COVERAGE_DIR)/coverage
	$Q cd $(BASE) && for pkg in $(TESTPKGS); do \
		$(GO) test -timeout 60s \
			-coverpkg=$$($(GO) list -f '{{ join .Deps "\n" }}' $$pkg | \
					grep '^$(PACKAGE)/' | grep -v '^$(PACKAGE)/vendor/' | \
					tr '\n' ',')$$pkg \
			-covermode=$(COVERAGE_MODE) \
			-coverprofile="$(COVERAGE_DIR)/coverage/`echo $$pkg | tr "/" "-"`.cover" $$pkg ;\
	 done
	$Q $(GOCOVMERGE) $(COVERAGE_DIR)/coverage/*.cover > $(COVERAGE_PROFILE)
	$Q $(GO) tool cover -html=$(COVERAGE_PROFILE) -o $(COVERAGE_HTML)
	$Q $(GOCOV) convert $(COVERAGE_PROFILE) | $(GOCOVXML) > $(COVERAGE_XML)

test-integration: fmt lint vendor test-coverage-tools | $(BASE) ; $(info $(M) running integration tests…) @ ## Run integration tests
	$Q cd $(BASE) && $(GO) test -v -tags=integration

.PHONY: lint
lint: vendor | $(BASE) $(GOLINT) ; $(info $(M) lint running golint…) @ ## Run golint
	$Q cd $(BASE) && ret=0 && for pkg in $(PKGS); do \
		test -z "$$($(GOLINT) $$pkg )" || ret=1 ; \
	 done ; exit $$ret

.PHONY: swagger
swagger: $(BASE) $(GOSWAGGER) ; $(info $(M) creating swagger…) @ ## Run golint
	$Q cd $(BASE) && ret=0 && for pkg in $(PKGS); do \
		test -z "$$($(GOSWAGGER) generate spec -o ./swagger.json  | tee /dev/stderr)" || ret=1 ; \
	 done ; exit $$ret

.PHONY: swaggerserver
swaggerserver: $(BASE) $(GOSWAGGER) ; $(info $(M) running swagger server…) @ ## Run golint
	$Q cd $(BASE) && ret=0 && for pkg in $(PKGS); do \
		test -z "$$($(GOSWAGGER) generate server -f ./swagger.json -A msok --principal msok  | tee /dev/stderr)" || ret=1 ; \
	 done ; exit $$ret


.PHONY: fmt
fmt: ; $(info $(M) running gofmt…) @ ## Run gofmt on all source files
	@ret=0 && for d in $$($(GO) list -f '{{.Dir}}' ./... | grep -v /vendor/); do \
		$(GOFMT) -l -w $$d/*.go || ret=$$? ; \
	 done ; exit $$ret

server.crt: ; $(info $(M) Creating self signed cert...) @ ## Run openssl
	@ret=0 && openssl req -new -newkey rsa:4096 -days 365 -nodes -x509 \
        -subj "/C=US/ST=CA/L=Mountain View/O=Dis/CN=selfsigned.com" \
        -keyout server.key  -out server.crt  || ret=$$?; \
        exit $$ret

.PHONY: package 
package: server.crt ; $(info $(M) Creating creating docker image...) @ ## Run openssl
	@ret=0 && \
        docker build .. -f ../Dockerfile -t localhost:5000/msok:latest || ret=$$?; \
        exit $$ret

.PHONY: stop
stop:  ; $(info $(M) Stoping docker container...) @ ## Stop running docker if any
	@ret=0 && \
	docker stop msok ; \
        exit 0


.PHONY: run
run:  stop ; $(info $(M) Running docker container...) @ ## Run openssl
	@ret=0 && \
	docker run -d --rm -p 8443:8443 -v ~/.idps/:/.idps/ --name msok localhost:5000/msok:latest || ret=$$?; \
        exit $$ret

.PHONY: logs
logs:  logs ; $(info $(M) Get docker logs for msok image...) @ ## Run openssl
	@ret=0 && \
	docker logs msok || ret=$$?; \
        exit $$ret

# Dependency management

Gopkg.lock: Gopkg.toml | $(BASE) ; $(info $(M) updating dependencies…)
	$Q cd $(BASE) && $(DEP) ensure -update
	@touch $@
#glide.lock: glide.yaml | $(BASE) ; $(info $(M) updating dependencies…)
# 	$Q cd $(BASE) && $(DEP) update
# 	@touch $@
vendor: Gopkg.lock | $(BASE) $(DEP) ; $(info $(M) retrieving dependencies…)
	$Q cd $(BASE) && $(DEP) ensure
	@ln -nsf . vendor/src
	@touch $@

#add-dependency: Gopkg.lock | $(BASE) ; $(info $(M) retrieving dependencies…)
#	$Q cd $(BASE) && $(DEP) ensure --add $(ADD_DEP)
#	@ln -nsf . vendor/src
#	@touch $@

# Misc

.PHONY: clean
clean: ; $(info $(M) cleaning…)	@ ## Cleanup everything
	@rm -rf $(GOPATH)
	@rm -rf bin
	@rm -rf test/tests.* test/coverage.*

.PHONY: help
help:
	@grep -E '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

.PHONY: version
version:
	@echo $(VERSION)

