### Global Variables
.SECONDEXPANSION:

DESTDIR = dist
MODULE = github.com/stonedu1011/devenvctl

EXECS = devenvctl

GEN_LIST = \

RES_LIST = \

### Main
.PHONY: generate clean test lint build copy-resources $(GEN_LIST) $(EXECS) $(RES_LIST)

## Required Variables by Local Targets
GO ?= go
CLI ?= lanai-cli

# target patterns
pGenerate = generate@%
pResource = resource@%

## Build AdHoc Targets
# generate:
# 	Invoke "go generate" on defined targets
# 	This target typically run on CI/CD working machine
generate: $(GEN_LIST)

# test:
# 	Invoke "go test" on defined modules.
# 	This target typically run on CI/CD working machine
# 	Optional parameter: ARGS="..."
ifneq ($(filter true True TRUE,$(SKIP_TEST)),)
test: generate
	@echo "Test Skipped..."
else
test: generate
	set -o pipefail; \
	gotestsum -f pkgname --jsonfile="$(DESTDIR)/tests.json" --raw-command -- \
  		$(GO) test -json -count=1 -failfast -timeout=0 -coverprofile $(DESTDIR)/coverage.out \
  		-coverpkg $(MODULE)/pkg/...,\
  		$(MODULE)/pkg/... $(ARGS)
endif

# lint:
# 	Invoke "go vet" and other linters
lint:
	$(GO) vet ./... 2>&1 | tee $(DESTDIR)/go-vet-report.out
	golangci-lint -c golangci.yml \
    	  --timeout 10m \
    	  --out-format colored-line-number,checkstyle:$(DESTDIR)/golangci-lint-report.xml \
    	  --issues-exit-code 0 run ./...

# build:
# 	Generate executable binary and copy resources to $(DESTDIR)
# 	this target should be run on targeted OS.
#	e.g. build is executed inside Docker container when building Docker image
# 	Optional Vars:
#		- VERSION version value without leading "v". Used for build info ldflags
build: $(EXECS) copy-resources

# install:
# 	Generate executable binary and install to $GOBIN
# 	this target should be run on targeted OS.
#	e.g. build is executed inside Docker container when building Docker image
# 	Optional Vars:
#		- VERSION version value without leading "v". Used for build info ldflags
install:
	$(GO) install github.com/stonedu1011/devenvctl

# copy-resources:
#	Copy resources to $(DESTDIR) based on $(RES_LIST)
# 	This target should be run on targeted OS.
#	e.g. build is executed inside Docker container when building Docker image
copy-resources: $(RES_LIST)

# clean:
# 	Undo previous "build".  clean $(DESTDIR) and build cache
# 	This target should be run on targeted OS.
clean:
	$(GO) clean
	rm -rf $(DESTDIR)/*

## Local Targets
# Generate
$(GEN_LIST):
	$(GO) generate $(@:$(pGenerate)=%)

# Build
devenvctl:
	#$(CLI) build -v "$(VERSION)" -- -o $(DESTDIR)/$@ main.go
	$(GO) build -o $(DESTDIR)/$@ main.go



