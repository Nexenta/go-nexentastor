# go-nexentastor

[![Build Status](https://travis-ci.org/Nexenta/go-nexentastor.svg?branch=master)](https://travis-ci.org/Nexenta/go-nexentastor)
[![Go Report Card](https://goreportcard.com/badge/github.com/Nexenta/go-nexentastor)](https://goreportcard.com/report/github.com/Nexenta/go-nexentastor)
[![Conventional Commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-yellow.svg)](https://conventionalcommits.org)
[![Releases](https://img.shields.io/github/tag/nexenta/go-nexentastor.svg)](https://github.com/Nexenta/go-nexentastor/releases)

Go API library for [NexentaStor](https://nexenta.com/products/nexentastor).

## Versions

See [CHANGELOG.md](./CHANGELOG.md)

## Documentation

To see docs for a specific version, change Git branch to version tag.

Current version documentation is [here](./docs).

## Development

Commits should follow [Conventional Commits Spec](https://conventionalcommits.org).

### Tests

See [Makefile](Makefile) for more examples.

```bash
# Test options:
# - TEST_NS_SINGLE=https://10.3.199.254:8443 # single NS
# - TEST_NS_HA_1=https://10.3.199.252:8443   # first node of HA NS cluster
# - TEST_NS_HA_2=https://10.3.199.253:8443   # second node of HA NS cluster
# - NOCOLORS=true                            # disable colors

# run all tests on local machine:
TEST_NS_SINGLE=https://10.3.199.251:8443 \
TEST_NS_HA_1=https://10.3.199.252:8443 \
TEST_NS_HA_2=https://10.3.199.253:8443 \
make test

# run all tests in containers:
TEST_NS_SINGLE=https://10.3.199.251:8443 \
TEST_NS_HA_1=https://10.3.199.252:8443 \
TEST_NS_HA_2=https://10.3.199.253:8443 \
make test-container
```

End-to-end NexentaStor test parameters:
```bash
# Tests for NexentaStor API provider (same options for `./resolver/resolver_test.go`)
go test ./tests/e2e/ns/provider/provider_test.go -v -count 1 \
    --address="https://10.3.199.254:8443" \
    --username="admin" \
    --password="pass" \
    --pool="myPool" \
    --dataset="myDataset" \
    --filesystem="myFs" \
    --cluster=true \
    --log=true
```

### Deps

To update deps run:
```bash
# go get -u github.com/golang/dep/cmd/dep
~/go/bin/dep ensure
```

### Release

Each API change MUST be released as git tag (X.X.X):

```bash
export NEXT_TAG=X.X.X
make pre-release-container
git diff
git add .
git commit -m "release ${NEXT_TAG}"
git push
git tag "${NEXT_TAG}"
git push --tags
```
