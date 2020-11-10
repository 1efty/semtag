package lib

import (
	"sort"
	"strings"

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
func GetFinalVersion(repository *git.Repository) *Version {
	var finalLeadingV bool
	var finalSemver = semver.New("0.0.0")
	var finalVersion = &Version{LeadingV: false, Semver: finalSemver}

	iter, err := getTags(repository)
	CheckIfError(err)

	// iterate through all tags, and determine which one is final
	err = iter.ForEach(func(ref *plumbing.Reference) error {
		var tempLeadingV bool
		var tempVersionString string

		tempVersionString, err := getTagString(repository, ref)
		CheckIfError(err)

		// handle leading `v`
		if strings.HasPrefix(tempVersionString, "v") {
			tempLeadingV = true
			tempVersionString = strings.TrimPrefix(tempVersionString, "v")
		}

		// create temp semver
		tempSemver := semver.New(tempVersionString)

		// change final variables if a tag was found that is newer, and is not an alpha, beta, or release-candidate
		if !tempSemver.LessThan(*finalSemver) && tempSemver.PreRelease == "" {
			finalSemver = tempSemver
			finalLeadingV = tempLeadingV
			finalVersion = &Version{LeadingV: finalLeadingV, Semver: tempSemver}
		}

		return nil
	})

	return finalVersion
}

// GetTagsAsVersion retrieves a slice of tags from a repository and converts them to Version objects
func GetTagsAsVersion(repository *git.Repository) []*Version {
	var tagsAsSemver Versions

	// Get all tags (annotated and light)
	iter, err := getTags(repository)
	CheckIfError(err)

	err = iter.ForEach(func(ref *plumbing.Reference) error {
		var versionString string
		var leadingV bool = false

		versionString, err = getTagString(repository, ref)
		CheckIfError(err)

		// handle leading `v`
		if strings.HasPrefix(versionString, "v") {
			leadingV = true
			versionString = strings.TrimPrefix(versionString, "v")
		}

		tagsAsSemver = append(tagsAsSemver, &Version{LeadingV: leadingV, Semver: semver.New(versionString)})

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
