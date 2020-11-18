export GO111MODULE=on
export GOSUMDB=off

.PHONY: linter
linter:
	$(info # Installing golangci-lint ...)
	$ curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./bin v1.31.0

.PHONY: test
test:
	$(info # Running app tests ...)
	$ go test -cover ./...

.PHONY: lint
lint:
	$(info # Running linter ...)
	$ ./bin/golangci-lint run --new-from-rev=origin/master --config=.pipeline.yaml --timeout=180s ./...

.PHONY: generate
generate:
	$(info # Generating stuff ...)
	$ go generate ./...

# CODE QUALITY
# https://github.com/codeclimate/codeclimate#packages
# $ docker pull --quiet "$CODE_QUALITY_IMAGE"
# registry.gitlab.com/gitlab-org/ci-cd/codequality:0.85.10-gitlab.1