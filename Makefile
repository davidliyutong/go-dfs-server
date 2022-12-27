
ROOT_PACKAGE = $(shell pwd)
DEMO_DATA_DIR = $(ROOT_PACKAGE)/demo_data
_BINARY_PREFIX = go-dfs-
AUTHOR = davidliyutong

include scripts/make-rules/common.mk
include scripts/make-rules/golang.mk

define USAGE_OPTIONS
	N_SERVERS: number of servers to start
endef
export USAGE_OPTIONS

.PHONY: clean
clean: demo.clean
	@echo "===========> Cleaning all build output"
	@-rm -vrf $(OUTPUT_DIR)

.PHONY: build
build:
	@$(MAKE) go.build

IMAGES_DIR ?= $(wildcard ${ROOT_DIR}/build/docker/*)
# Determine images names by stripping out the dir names
IMAGES ?= $(filter-out tools,$(foreach image,${IMAGES_DIR},$(notdir ${image})))

.PHONY: image.build.%
image.build.%:
	$(eval IMAGE := $*)
	$(eval IMAGE_PLAT := $(subst _,/,$(PLATFORM)))
	@echo "===========> Building docker image $(IMAGE) $(VERSION) for $(IMAGE_PLAT)"
	@docker build --platform $(IMAGE_PLAT) -t "$(AUTHOR)/$(_BINARY_PREFIX)$(IMAGE):$(VERSION)-$(GOOS)-$(GOARCH)" --file ./build/docker/$(IMAGE)/Dockerfile .

.PHONY: image.build
image.build: $(foreach image,${IMAGES},image.build.${image})

.PHONY: image.push.%
image.push.%:
	$(eval IMAGE := $*)
	@echo "===========> Pushing docker image $(IMAGE) $(VERSION)"
	@docker push "$(AUTHOR)/$(_BINARY_PREFIX)$(IMAGE):$(VERSION)-$(GOOS)-$(GOARCH)"

.PHONY: image.push
image.push: $(foreach image,${IMAGES},image.push.${image})

.PHONY: image
image:
	@$(MAKE) image.build

.PHONY: image.clean
image.clean:
	@echo "===========> Cleaning all docker images"
	@-docker rmi -f $(shell docker images -q $(AUTHOR)/$(_BINARY_PREFIX)*)

demo.create:
	@-mkdir -p $(DEMO_DATA_DIR)

.PHONY: demo
demo: demo.create
	@$(MAKE) demo.start

.PHONY: demo.start
 demo.start: demo.create
	$(eval TAG := $(VERSION)-$(GOOS)-$(GOARCH))
	$(eval N_SERVERS ?= 4)
	@bash ./scripts/launch_all_servers.sh $(DEMO_DATA_DIR) $(TAG) $(N_SERVERS) dfs

.PHONY: demo.prepare
demo.prepare: demo.create image.build

#.PHONY: demo.stop
demo.stop:
	@-bash ./scripts/stop_all_servers.sh

.PHONY: demo.clean
demo.clean:
	@-rm -vrf $(DEMO_DATA_DIR)

.PHONY: all
all: clean build