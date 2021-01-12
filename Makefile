.PHONY: container changelog ci/build

container:
	docker build . -t semtag

changelog:
	git-chglog -o CHANGELOG.md

ci/build:
	docker run -it --rm -v $(PWD):/src -w /src goreleaser/goreleaser:latest build --snapshot --rm-dist

ci/release/dryrun:
	docker run -it --rm -v $(PWD):/src -w /src goreleaser/goreleaser:latest release --snapshot --rm-dist --skip-publish
