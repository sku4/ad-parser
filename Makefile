CURDIR=$(shell pwd)
BINDIR=${CURDIR}/.bin
GOVER=$(shell go version | perl -nle '/(go\d\S+)/; print $$1;')
MOCKGEN=${BINDIR}/mockgen_${GOVER}
SMARTIMPORTS=${BINDIR}/smartimports_${GOVER}
LINTVER=v1.49.0
LINTBIN=${BINDIR}/lint_${GOVER}_${LINTVER}

.PHONY:
run:
	docker-compose up --remove-orphans --build app

test:
	go test ./... -coverprofile cover.out

test-coverage:
	go tool cover -func cover.out | grep total | awk '{print $3}'

bindir:
	mkdir -p ${BINDIR}

lint:
	golangci-lint run

install-lint: bindir
	test -f ${LINTBIN} || \
		(GOBIN=${BINDIR} go install github.com/golangci/golangci-lint/cmd/golangci-lint@${LINTVER} && \
		mv ${BINDIR}/golangci-lint ${LINTBIN})
