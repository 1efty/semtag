package lib

import (
	"strings"

	"github.com/coreos/go-semver/semver"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

// GetRepository tries to open current directory as git repository
func GetRepository() *git.Repository {
	repository, err := git.PlainOpen(".")
	CheckIfError(err)
	return repository
}

// GetTagsAsSemver retrieves slice of tags from repository as coreos.Semver objects
func GetTagsAsSemver(repository *git.Repository) []*Version {
	var tagsAsSemver []*Version

	// Get all tags (annotated and light)
	iter, err := repository.Tags()
	CheckIfError(err)

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

		tagsAsSemver = append(tagsAsSemver, &Version{LeadingV: leadingV, Semver: semver.New(versionString)})

		return nil
	})
	CheckIfError(err)

	return tagsAsSemver
}

// CreateTag tags HEAD in a given repository
func CreateTag(repository *git.Repository, tag string) error {
	h, err := repository.Head()
	CheckIfError(err)

	_, err = repository.CreateTag(tag, h.Hash(), &git.CreateTagOptions{
		Message: tag,
	})
	CheckIfError(err)

	return nil
}

// BumpVersion creates new version and bumps according to scope
func BumpVersion(v *Version, scope string, preRelease string, metadata string) (*Version, error) {
	newVersion := &Version{
		LeadingV: v.LeadingV,
		Semver:   v.Semver,
	}
	err := newVersion.Bump(scope)
	CheckIfError(err)

	// set pre-release and metadata
	newVersion.Semver.PreRelease = semver.PreRelease(preRelease)
	newVersion.Semver.Metadata = metadata

	return newVersion, nil
}
