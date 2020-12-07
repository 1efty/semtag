package git

import (
	"sort"
	"strings"

	"github.com/1efty/semtag/pkg/utils"
	"github.com/1efty/semtag/pkg/version"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/storer"
)

// Repo represents a git repository
type Repo struct {
	repo *git.Repository

	Tags []*version.Version

	Status git.Status

	CurrentVersion *version.Version
	FinalVersion   *version.Version
	FirstVersion   *version.Version
	LastVersion    *version.Version
}

// New creates a new repo from a string
func New(dir string) *Repo {
	repository := getRepository(dir)
	tags := getTagsAsVersion(repository)

	// determine first, last, current, and final version
	var firstVersion, lastVersion, currentVersion, finalVersion *version.Version

	switch numOfTags := len(tags); numOfTags {
	case 0:
		firstVersion = version.New("0.0.0")
		lastVersion = firstVersion
		currentVersion = firstVersion
		finalVersion = firstVersion
	case 1:
		firstVersion = tags[0]
		lastVersion = firstVersion
		currentVersion = firstVersion
		finalVersion = firstVersion
	default:
		firstVersion = tags[0]
		lastVersion = tags[len(tags)-1]
		currentVersion = lastVersion
		finalVersion = getFinalVersion(repository)
	}

	return &Repo{
		repo:           repository,
		Tags:           getTagsAsVersion(repository),
		Status:         getStatus(repository),
		FirstVersion:   firstVersion,
		LastVersion:    lastVersion,
		CurrentVersion: currentVersion,
		FinalVersion:   finalVersion,
	}
}

// CreateTag tags HEAD in a given repository
func (r Repo) CreateTag(tag string) error {
	h, err := r.repo.Head()
	utils.CheckIfError(err)

	_, err = r.repo.CreateTag(tag, h.Hash(), &git.CreateTagOptions{
		Message: tag,
	})
	utils.CheckIfError(err)

	return nil
}

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

func getTagsAsVersion(repository *git.Repository) []*version.Version {
	var tagsAsSemver version.Versions

	// Get all tags (annotated and light)
	iter, err := getTags(repository)
	utils.CheckIfError(err)

	err = iter.ForEach(func(ref *plumbing.Reference) error {
		var versionString string

		versionString, err = getTagString(repository, ref)
		utils.CheckIfError(err)

		tagsAsSemver = append(tagsAsSemver, version.New(versionString))

		return nil
	})
	utils.CheckIfError(err)

	sort.Sort(tagsAsSemver)

	return tagsAsSemver
}

func getStatus(repository *git.Repository) git.Status {
	worktree, err := repository.Worktree()
	utils.CheckIfError(err)
	status, _ := worktree.Status()
	return status
}

func getRepository(dir string) *git.Repository {
	repository, err := git.PlainOpen(dir)
	utils.CheckIfError(err)
	return repository
}

func getFinalVersion(repository *git.Repository) *version.Version {
	var finalVersion = version.New("0.0.0")

	iter, err := getTags(repository)
	utils.CheckIfError(err)

	// iterate through all tags, and determine which one is final
	err = iter.ForEach(func(ref *plumbing.Reference) error {
		var tempVersionString string

		tempVersionString, err := getTagString(repository, ref)
		utils.CheckIfError(err)

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
