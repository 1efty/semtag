.PHONY: changelog release

changelog:
	git-chglog -o CHANGELOG.md

release:
	semtag final -s minor
