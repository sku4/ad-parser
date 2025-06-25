CURDIR=$(shell pwd)
BINDIR=${CURDIR}/.bin
GOVER=$(shell go version | perl -nle '/(go\d\S+)/; print $$1;')
MOCKGEN=${BINDIR}/mockgen_${GOVER}
SMARTIMPORTS=${BINDIR}/smartimports_${GOVER}
LINTVER=v1.61.0
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

helm-install:
	helm upgrade --install "ad-parser" .helm --namespace=ad-prod

helm-install-local:
	helm upgrade --install "ad-parser" .helm \
		--namespace=ad-prod \
		-f ./.helm/values-local.yaml \
		--wait \
		--timeout 300s \
		--atomic \
		--debug

helm-template:
	helm template --name-template="ad-parser" \
		--namespace=ad-prod \
		-f .helm/values-local.yaml .helm \
		> .helm/helm.txt

helm-package:
	helm package .helm
	mv ad-parser*.tgz docs/charts
	helm repo index docs/charts --url https://raw.githubusercontent.com/sku4/ad-parser/refs/heads/master/docs/charts/
