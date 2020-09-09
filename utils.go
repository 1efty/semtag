package main

import (
	"fmt"
	"os"

	"github.com/coreos/go-semver/semver"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// checkIfError ... checks err and exists if necessary
func checkIfError(err error) {
	if err == nil {
		return
	}

	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	os.Exit(1)
}

// info ... print info to screen
func info(format string, args ...interface{}) {
	fmt.Printf("\x1b[34;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}

// getRepository ... tries to open current directory as git repository
func getRepository() *git.Repository {
	repository, err := git.PlainOpen(".")
	checkIfError(err)
	return repository
}

// getTagsAsSemver ... get slice of tags from repository
func getTagsAsSemver(repository *git.Repository) []*semver.Version {
	var tagsAsSemver []*semver.Version

	tagObjs, err := repository.TagObjects()
	checkIfError(err)

	err = tagObjs.ForEach(func(t *object.Tag) error {
		tagsAsSemver = append(tagsAsSemver, semver.New(t.Name))
		return nil
	})
	checkIfError(err)

	// sort slice
	semver.Sort(tagsAsSemver)

	return tagsAsSemver
}

// createTag ... create a tag given a string
func createTag(repository *git.Repository, tag string) error {
	h, err := repository.Head()
	checkIfError(err)

	_, err = repository.CreateTag(tag, h.Hash(), &git.CreateTagOptions{
		Message: tag,
	})
	checkIfError(err)

	return nil
}

// bumpVersion ... create new version and bump according to scope
func bumpVersion(v *semver.Version, scope string, preRelease string, metadata string) (*semver.Version, error) {

	newVersion := &semver.Version{
		Major: v.Major,
		Minor: v.Minor,
		Patch: v.Patch,
	}

	switch scope {
	case "patch":
		newVersion.BumpPatch()
	case "minor":
		newVersion.BumpMinor()
	case "major":
		newVersion.BumpMajor()
	default:
		newVersion.BumpPatch()
	}

	// set pre-release and metadata
	newVersion.PreRelease = semver.PreRelease(preRelease)
	newVersion.Metadata = metadata

	return newVersion, nil
}