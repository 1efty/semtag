package lib

import (
	"sort"
	"strings"

	"github.com/1efty/semtag/pkg/version"
	"github.com/coreos/go-semver/semver"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/storer"
)

func getTags(repository *git.Repository) (storer.ReferenceIter, error) {
	return repository.Tags()
}

func getTagString(repository *git.Repository, ref *plumbing.Reference) (string, error) {
	var versionString string
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

	return versionString, nil
}

// GetStatus returns the status of the worktree
func GetStatus(repository *git.Repository) (git.Status, error) {
	worktree, err := repository.Worktree()
	CheckIfError(err)
	status, err := worktree.Status()
	return status, err
}

// GetRepository tries to open current directory as git repository
func GetRepository() *git.Repository {
	repository, err := git.PlainOpen(".")
	CheckIfError(err)
	return repository
}

// GetFinalVersion retrieves the last tag that is neither an alpha, beta, or release-candidate
func GetFinalVersion(repository *git.Repository) *version.Version {
	var finalVersion = version.New("0.0.0")

	iter, err := getTags(repository)
	CheckIfError(err)

	// iterate through all tags, and determine which one is final
	err = iter.ForEach(func(ref *plumbing.Reference) error {
		var tempVersionString string

		tempVersionString, err := getTagString(repository, ref)
		CheckIfError(err)

		// create temp version
		tempVersion := version.New(tempVersionString)

		// change final variables if a tag was found that is newer, and is not an alpha, beta, or release-candidate
		if !tempVersion.Semver.LessThan(*finalVersion.Semver) && tempVersion.Semver.PreRelease == "" {
			finalVersion = tempVersion
		}

		return nil
	})

	return finalVersion
}

// GetTagsAsVersion retrieves a slice of tags from a repository and converts them to Version objects
func GetTagsAsVersion(repository *git.Repository) []*version.Version {
	var tagsAsSemver version.Versions

	// Get all tags (annotated and light)
	iter, err := getTags(repository)
	CheckIfError(err)

	err = iter.ForEach(func(ref *plumbing.Reference) error {
		var versionString string

		versionString, err = getTagString(repository, ref)
		CheckIfError(err)

		tagsAsSemver = append(tagsAsSemver, version.New(versionString))

		return nil
	})
	CheckIfError(err)

	sort.Sort(tagsAsSemver)

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
func BumpVersion(v *version.Version, scope string, preRelease string, metadata string) (*version.Version, error) {
	newVersion := v
	err := newVersion.Bump(scope)
	CheckIfError(err)

	// set pre-release and metadata
	newVersion.Semver.PreRelease = semver.PreRelease(preRelease)
	newVersion.Semver.Metadata = metadata

	return v, nil
}
