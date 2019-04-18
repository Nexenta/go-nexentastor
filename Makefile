# NexentaStor Go library makefile
#
# Test options to be set before run tests:
# - NOCOLOR=1                                # disable colors
# - TEST_NS_SINGLE=https://10.3.199.254:8443 # single NS provider/resolver tests
# - TEST_NS_HA_1=https://10.3.199.252:8443   # HA cluster NS provider/resolver tests
# - TEST_NS_HA_2=https://10.3.199.253:8443
#

DOCKER_FILE_TESTS = Dockerfile.tests
DOCKER_IMAGE_TESTS = go-nexentastor-tests

.PHONY: all
all: test

.PHONY: test
test: test-unit test-e2e-ns-single test-e2e-ns-cluster
.PHONY: test-container
test-container: \
	test-unit-container \
	test-e2e-ns-single-container \
	test-e2e-ns-cluster-container

.PHONY: test-unit
test-unit:
	go test ./tests/unit/rest -v -count 1
	go test ./tests/unit/ns -v -count 1
.PHONY: test-unit-container
test-unit-container:
	docker build -f ${DOCKER_FILE_TESTS} -t ${DOCKER_IMAGE_TESTS} .
	docker run -i --rm -e NOCOLORS=${NOCOLORS} ${DOCKER_IMAGE_TESTS} test-unit

.PHONY: test-e2e-ns-single
test-e2e-ns-single: check-env-TEST_NS_SINGLE
	go test ./tests/e2e/ns/provider/provider_test.go -v -count 1 --address="${TEST_NS_SINGLE}"
	go test ./tests/e2e/ns/resolver/resolver_test.go -v -count 1 --address="${TEST_NS_SINGLE}"
.PHONY: test-e2e-ns-single-container
test-e2e-ns-single-container: check-env-TEST_NS_SINGLE
	docker build -f ${DOCKER_FILE_TESTS} -t ${DOCKER_IMAGE_TESTS} .
	docker run -i --rm -v ${HOME}/.ssh:/root/.ssh:ro \
		-e NOCOLORS=${NOCOLORS} -e TEST_NS_SINGLE=${TEST_NS_SINGLE} \
		${DOCKER_IMAGE_TESTS} test-e2e-ns-single

.PHONY: test-e2e-ns-cluster
test-e2e-ns-cluster: check-env-TEST_NS_HA_1 check-env-TEST_NS_HA_2
	go test ./tests/e2e/ns/provider/provider_test.go -v -count 1 --address="${TEST_NS_HA_1}" --cluster=true
	go test ./tests/e2e/ns/resolver/resolver_test.go -v -count 1 --address="${TEST_NS_HA_1},${TEST_NS_HA_2}"
.PHONY: test-e2e-ns-cluster-container
test-e2e-ns-cluster-container: check-env-TEST_NS_HA_1 check-env-TEST_NS_HA_2
	docker build -f ${DOCKER_FILE_TESTS} -t ${DOCKER_IMAGE_TESTS} .
	docker run -i --rm -v ${HOME}/.ssh:/root/.ssh:ro \
		-e NOCOLORS=${NOCOLORS} -e TEST_NS_HA_1=${TEST_NS_HA_1} -e TEST_NS_HA_2=${TEST_NS_HA_2} \
		${DOCKER_IMAGE_TESTS} test-e2e-ns-cluster

.PHONY: check-env-TEST_NS_SINGLE
check-env-TEST_NS_SINGLE:
ifeq ($(strip ${TEST_NS_SINGLE}),)
	$(error "Error: environment variable TEST_NS_SINGLE is not set (e.i. https://10.3.199.254:8443)")
endif

.PHONY: check-env-TEST_NS_HA_1
check-env-TEST_NS_HA_1:
ifeq ($(strip ${TEST_NS_HA_1}),)
	$(error "Error: environment variable TEST_NS_HA_1 is not set (e.i. https://10.3.199.254:8443)")
endif

.PHONY: check-env-TEST_NS_HA_2
check-env-TEST_NS_HA_2:
ifeq ($(strip ${TEST_NS_HA_2}),)
	$(error "Error: environment variable TEST_NS_HA_2 is not set (e.i. https://10.3.199.254:8443)")
endif

.PHONY: clean
clean:
	go clean -r -x
