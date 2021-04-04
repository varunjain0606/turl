# Assignment is done to avoid re-evaluation, that is used in subsequent target
BUILD_DATE ?=$(shell date -u '+%Y-%m-%d_%I:%M:%S%p')

# Allow customizing build components.
SUBDIRS ?= api_service

all build clean docs image publish publish-no-deps test : $(SUBDIRS)
$(SUBDIRS):
	@echo "Make in $@"
	@$(MAKE) -C $@ BUILD_DATE=$(BUILD_DATE)

.PHONY: all build clean docs image publish publish-no-deps test $(SUBDIRS)

.PHONY: docker-rm
docker-rm:
	@# pass ARGS=-f to kill the processes forcefully.
	docker ps --all | awk '{print $$1}' | xargs docker rm ${ARGS}

.PHONY: docker-image-rm
docker-image-rm:
	@# pass ARGS=-f to remove the images forcefully.
	(docker image ls | awk '{print $$3}' | xargs docker image rm ${ARGS}) > /dev/null || true

.PHONY: list-build-modules
list-build-modules:
	@# Run following command to clear the module cache.
	@# sudo rm -rf ~/go/pkg/mod

	@# clean up all the components.
	(make clean SUBDIRS='api_service')

	@# make build will fill out all the third-party modules
	(make build SUBDIRS='api_service')