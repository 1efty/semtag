package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/coreos/go-semver/semver"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
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
func getTagsAsSemver(repository *git.Repository) []*Version {
	var tagsAsSemver []*Version

	// Get all tags (annotated and light)
	iter, err := repository.Tags()
	checkIfError(err)

	err = iter.ForEach(func(ref *plumbing.Reference) error {
		var versionString string = "0.0.0"
		var leadingV bool = false
		obj, err := repository.TagObject(ref.Hash())

		// check if annotated tag
		switch err {
		case nil:
			// If annotated, can simply take the Name
			versionString = obj.Name
		case plumbing.ErrObjectNotFound:
			// If not, will need to do some hacking
			versionString = strings.Split(ref.String(), "/")[2]
		}

		// handle leading `v`
		if strings.HasPrefix(versionString, "v") {
			leadingV = true
			versionString = strings.TrimPrefix(versionString, "v")
		}

		tagsAsSemver = append(tagsAsSemver, &Version{leadingV: leadingV, semver: semver.New(versionString)})

		return nil
	})
	checkIfError(err)

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
func bumpVersion(v *Version, scope string, preRelease string, metadata string) (*Version, error) {
	newVersion := &Version{
		leadingV: v.leadingV,
		semver:   v.semver,
	}
	err := newVersion.Bump(scope)
	checkIfError(err)

	// set pre-release and metadata
	newVersion.semver.PreRelease = semver.PreRelease(preRelease)
	newVersion.semver.Metadata = metadata

	return newVersion, nil
}
