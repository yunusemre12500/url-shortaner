DOCKER ?= docker
GO ?= go
RELEASE ?= unknown

.PHONY: api-server-image
api-server-image:
	@$(DOCKER) buildx build \
		--attest=type=sbom,mode=max \
		--attest=type=provenance \
		--build-arg=RELEASE=$(RELEASE) \
		--progress=plain \
		--tag ghcr.io/yunusemre12500/url-shortaner/api-server:$(RELEASE) \
		.

.PHONY: tidy
tidy:
	@$(GO) mod tidy

.PHONY: vendor
vendor:
	@$(GO) mod vendo
