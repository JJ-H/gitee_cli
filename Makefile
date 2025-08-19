# Makefile for cross-platform build
BINARY_NAME = gitee
NPM_VERSION = 0.1.2
GO = go
OSES = darwin linux windows
ARCHS = amd64 arm64


# Repository information
GITEE_OWNER ?= "JJ-H"
GITEE_REPO ?= "gitee_cli"

# Flags
LDFLAGS = -ldflags "-s -w"
BUILD_FLAGS = -o bin/gitee $(LDFLAGS)

define show_usage_info
	@echo "\033[32m\n🤖🤖 Build Success 🤖🤖\033[0m"
endef

build:
	$(GO) build $(BUILD_FLAGS) -v main.go
	@echo "Build complete."
	$(call show_usage_info)

# Clean up generated binaries
clean:
	rm -f bin/gitee
	@echo "Clean up complete."


.PHONY: build-all-platforms
build-all-platforms:
	$(foreach os,$(OSES),$(foreach arch,$(ARCHS), \
		GOOS=$(os) GOARCH=$(arch) go build $(BUILD_FLAGS) -o $(BINARY_NAME)-$(os)-$(arch)$(if $(findstring windows,$(os)),.exe,) main.go; \
	))

.PHONY: npm-copy-binaries
npm-copy-binaries: build-all-platforms
	$(foreach os,$(OSES),$(foreach arch,$(ARCHS), \
		EXECUTABLE=./$(BINARY_NAME)-$(os)-$(arch)$(if $(findstring windows,$(os)),.exe,); \
		DIRNAME=$(BINARY_NAME)-$(os)-$(arch); \
		mkdir -p ./npm/$$DIRNAME/bin; \
		cp $$EXECUTABLE ./npm/$$DIRNAME/bin/; \
	))

.PHONY: npm-publish
npm-publish: npm-copy-binaries ## Publish the npm packages
	$(foreach os,$(OSES),$(foreach arch,$(ARCHS), \
		DIRNAME="$(BINARY_NAME)-$(os)-$(arch)"; \
		cd npm/$$DIRNAME; \
		echo '//registry.npmjs.org/:_authToken=$(NPM_TOKEN)' >> .npmrc; \
		jq '.version = "$(NPM_VERSION)"' package.json > tmp.json && mv tmp.json package.json; \
		npm publish; \
		cd ../..; \
	))
	cp README.md LICENSE ./npm/gitee/
	echo '//registry.npmjs.org/:_authToken=$(NPM_TOKEN)' >> ./npm/gitee/.npmrc
	jq '.version = "$(NPM_VERSION)"' ./npm/gitee/package.json > tmp.json && mv tmp.json ./npm/gitee/package.json; \
	jq '.optionalDependencies |= with_entries(.value = "$(NPM_VERSION)")' ./npm/gitee/package.json > tmp.json && mv tmp.json ./npm/gitee/package.json; \
	cd npm/gitee && npm publish


# Clean up release directory
clean-release:
	rm -rf release
	@echo "Clean up release directory complete."

# Create a tarball for the given platform
define create_tarball
	@echo "Packaging for $(1)..."
	@mkdir -p release/$(1)
	@cp bin/gitee release/$(1)/gitee$(2)
	@cp LICENSE release/$(1)/
	@cp README.md release/$(1)/
	@cp README_CN.md release/$(1)/
	@tar -czvf release/gitee-$(1).tar.gz -C release/$(1) .
	@rm -rf release/$(1)
endef

release: clean clean-release
	@mkdir -p release
	@echo "Building for Linux..."
	GOOS=linux GOARCH=amd64 $(GO) build $(BUILD_FLAGS) -v main.go
	$(call create_tarball,linux-amd64,)
	@echo "Building for Windows..."
	GOOS=windows GOARCH=amd64 $(GO) build $(BUILD_FLAGS) -v main.go
	$(call create_tarball,windows-amd64,.exe)
	@echo "Building for macOS..."
	GOOS=darwin GOARCH=amd64 $(GO) build $(BUILD_FLAGS) -v main.go
	$(call create_tarball,darwin-amd64,)
	@echo "Building for macOS ARM..."
	GOOS=darwin GOARCH=arm64 $(GO) build $(BUILD_FLAGS) -v main.go
	$(call create_tarball,darwin-arm64,)
	@echo "Building for Linux ARM..."
	GOOS=linux GOARCH=arm $(GO) build $(BUILD_FLAGS) -v main.go
	$(call create_tarball,linux-arm,)
	@echo "Release complete. Artifacts are in the release directory."

# Upload artifacts to a specific release
upload-gitee-release:
	@echo "Uploading artifacts to gitee release..."
	@for file in release/*; do \
		curl -X POST \
			-H "Content-Type: multipart/form-data" \
			-F "access_token=$(GITEE_ACCESS_TOKEN)" \
			-F "owner=$(GITEE_OWNER)" \
			-F "repo=$(GITEE_REPO)" \
			-F "release_id=$(GITEE_RELEASE_ID)" \
			-F "file=@$$file" \
			https://gitee.com/api/v5/repos/$(GITEE_OWNER)/$(GITEE_REPO)/releases/$(GITEE_RELEASE_ID)/attach_files; \
	done
	@echo "Upload complete."
