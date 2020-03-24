MAKEFLAGS += -r --warn-undefined-variables
SHELL := /bin/bash
.SHELLFLAGS := -o pipefail -euc
.DEFAULT_GOAL := help

include Makefile.variables
include Makefile.local

.PHONY: help clean veryclean tag-build build build-image dep* format check test cover generate docs todo adhoc next-dev start-release finish-release pub-image xcompile

## display this help message
help:
	@echo 'Management commands for streaming-service:'
	@echo
	@echo 'Usage:'
	@echo
	@echo '  ## Build Commands'
	@echo '    build           Compile the project.'
	@echo '  ## Demo Command'
	@echo '    demo            automated demonstration'
	@echo '  ## Develop / Test Commands'
	@echo '    dep             Update go modules.'
	@echo '    dep-verify      Verify go modules.'
	@echo '    dep-why         Question the need of imported go modules.'
	@echo '    dep-cache       cache go modules.'
	@echo '    format          Run code formatter.'
	@echo '    check           Run static code analysis (lint).'
	@echo '    test            Run tests on project.'
	@echo '    cover           Run tests and capture code coverage metrics on project.'
	@echo '    todo            Generate a TODO list for project.'
	@echo
	@echo '  ## Local Commands'
	@echo '    drma            Removes all stopped containers.'
	@echo '    drmia           Removes all unlabelled images.'
	@echo '    drmvu           Removes all unused container volumes.'
	@echo

## Compile the project.
build: build/dev

build/dev:
	@rm -rf bin/
	@mkdir -p bin
	${DOCKER} bash ./scripts/build.sh

## Clean the directory tree of produced artifacts.
clean:
	@rm -rf bin build release cover *.out *.xml

## Same as clean but also removes cached dependencies.
veryclean: clean
	@rm -rf tmp

## builds the dev container
prepare: tmp/dev_image_id
tmp/dev_image_id: Dockerfile.dev
	@mkdir -p tmp
	@docker rmi -f ${DEV_IMAGE} > /dev/null 2>&1 || true
	echo "Spawning dev container..."
	@docker build		--quiet -t ${DEV_IMAGE} --build-arg DEVELOPER=$(shell whoami) -f Dockerfile.dev .
	@docker inspect -f "{{ .ID }}" ${DEV_IMAGE} > tmp/dev_image_id

# ----------------------------------------------
# develop and test

# ----------------------------------------------
# dependencies

## Update modules
dep: prepare
	@go mod vendor

## verify modules
dep-verify: prepare
	@go mod verify

## question modules
dep-why: prepare
	@go mod why

## cache modules
dep-cache: prepare
	@go mod download

## Run code formatter.
format: dep
	${DOCKER} bash ./scripts/format.sh

## Run static code analysis (lint).
check: format
	${DOCKER} bash ./scripts/check.sh

db-start:	db-remove
		@mkdir -p tmp
		@echo "starting db..."
		@docker run -p 9042:9042 --name cassandra-db-cnt -d cassandra:3 > tmp/cassandra_db_cnt_id
		@sleep 30

db-stop:
		@$(eval CASSANDRA_DB_CNT_ID=$(shell cat tmp/cassandra_db_cnt_id 2>/dev/null))
		@if [ -f tmp/cassandra_db_cnt_id ]; then \
						docker stop $(CASSANDRA_DB_CNT_ID) > /dev/null 2>&1 || : ; \
						rm -f tmp/cassandra_db_cnt_id; \
		fi

db-remove:
		@$(eval CASSANDRA_DB_CNT_ID=$(shell cat tmp/cassandra_db_cnt_id 2>/dev/null))
		@if [ -f tmp/cassandra_db_cnt_id ]; then \
						docker rm -f $(CASSANDRA_DB_CNT_ID) > /dev/null 2>&1 || : ; \
						rm -f tmp/cassandra_db_cnt_id; \
		fi

seed-db:
	@$(eval CASSANDRA_DB_CNT_ID=$(shell cat tmp/cassandra_db_cnt_id))
	@echo   "setting up db..."
	@docker cp schema.cql $(CASSANDRA_DB_CNT_ID):/tmp/schema.cql
	@docker exec -i $(CASSANDRA_DB_CNT_ID) sh -c 'cqlsh -f /tmp/schema.cql'

.pretest: check db-start seed-db

demo: db-start seed-db
	@$(eval ZATOO_CNT_ID=$(shell cat tmp/zatoo_cnt_id 2>/dev/null))
	@if [ -f tmp/zatoo_cnt_id ]; then \
		docker rm -f $(ZATOO_CNT_ID) > /dev/null 2>&1 || : ; \
		rm -f tmp/zattoo_cnt_id; \
	fi
	@echo "starting channel source"
	@docker run -p 9110:9110 --name zattoo-channel-source -d -e PORT=9110 -e EXECUTION_DURATION=1000sec fokhunov/zattoo:channelsource-1.0 > tmp/zatoo_cnt_id
	bash ./scripts/demo.sh

## Run tests on project.
test: .pretest
	${DOCKERTEST} bash ./scripts/test.sh
	@${MAKE} --no-print-directory db-remove

## Run tests and capture code coverage metrics on project.
cover: .pretest
	@rm -rf cover/
	@mkdir -p cover
	${DOCKERTEST} bash ./scripts/cover.sh
