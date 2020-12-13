.PHONY: build/container changelog ci/build ci/release/dryrun release test

export GOVERSION = $(shell cat .go-version)

build/container:
	docker build . -t semtag --build-arg GOVERSION

changelog:
	git-chglog -o CHANGELOG.md

ci/build:
	docker run -it --rm -v $(PWD):/src -w /src goreleaser/goreleaser:latest build --snapshot --rm-dist

ci/release/dryrun:
	docker run -it --rm -v $(PWD):/src -w /src goreleaser/goreleaser:latest release --snapshot --rm-dist --skip-publish

test:
	go test github.com/1efty/semtag/cmd

release:
	semtag final -s minor
