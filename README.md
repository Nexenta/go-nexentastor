# go-nexentastor

[![Build Status](https://travis-ci.org/Nexenta/go-nexentastor.svg?branch=master)](https://travis-ci.org/Nexenta/go-nexentastor)
[![Go Report Card](https://goreportcard.com/badge/github.com/Nexenta/go-nexentastor)](https://goreportcard.com/report/github.com/Nexenta/go-nexentastor)
[![Conventional Commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-yellow.svg)](https://conventionalcommits.org)

Go API library for [NexentaStor](https://nexenta.com/products/nexentastor).

## Versions

See [CHANGELOG.md](./CHANGELOG.md)

## Usage

TODO

## Development

Commits should follow [Conventional Commits Spec](https://conventionalcommits.org).

### Tests

See [Makefile](Makefile) for more examples.

```bash
# update deps
~/go/bin/dep ensure

# Test options to be set before run tests:
# - TEST_NS_SINGLE=https://10.3.199.254:8443 # single NS provider/resolver tests
# - TEST_NS_HA_1=https://10.3.199.252:8443   # HA cluster NS provider/resolver tests
# - TEST_NS_HA_2=https://10.3.199.253:8443
TEST_NS_SINGLE=https://10.3.199.251:8443 \
TEST_NS_HA_1=https://10.3.199.252:8443 \
TEST_NS_HA_2=https://10.3.199.253:8443 \
make test

# run tests in container:
# - RSA keys from host's ~/.ssh directory will be used by container.
#   Make sure all remote hosts used in tests have host's RSA key added as trusted
#   (ssh-copy-id -i ~/.ssh/id_rsa.pub user@host)
# - "export NOCOLORS=true" to run w/o colors
TEST_NS_SINGLE=https://10.3.199.251:8443 \
TEST_NS_HA_1=https://10.3.199.252:8443 \
TEST_NS_HA_2=https://10.3.199.253:8443 \
make test-container
```

End-to-end NexentaStor/K8s test parameters:
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

### Release

Each API change MUST be released as git tag.

```bash
make test-container
git-chglog --next-tag X.X.X -o CHANGELOG.md # go get -u github.com/git-chglog/git-chglog/cmd/git-chglog
git add CHANGELOG.md
git commit -m "release X.X.X"
git tag X.X.X
git push --tags
```
