TARGETS := wlint

.PHONY: \
	all \
	clean \
	coverage \
	format \
	install \
	lint \
	test \
	${TARGETS}

all: ${TARGETS}

clean:
	go clean

format:
	find -name '*.go' | xargs gofmt -w -s

wlint:
	go build -o ${@}

test:
	go test -race ./...

lint:
	go vet ./...
	staticcheck ./...

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

install: wlint
	if [ -n "${DESTDIR}" ]; then \
		mkdir -p "${DESTDIR}/bin/"; \
		cp wlint "${DESTDIR}/bin/"; \
	else \
		mkdir -p "$(shell go env GOPATH)/bin"; \
		cp wlint "$(shell go env GOPATH)/bin"; \
	fi
