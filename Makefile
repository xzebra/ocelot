LINTER_URL=golang.org/x/lint/golint
LINTER=golint
LINTER_FLAGS=-min_confidence 0.5

export GO111MODULE=on

build: get lint
	@go build

dep:
	@go install $(LINTER_URL)

get:
	@go get

test:
	@go test -vc

lint: get dep
	@$(LINTER) $(LINTER_FLAGS)
