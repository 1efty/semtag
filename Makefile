.PHONY: changelog release

container:
	docker build . -t semtag

changelog:
	git-chglog -o CHANGELOG.md

release:
	semtag final -s minor
