package main

import (
	"errors"

	"github.com/coreos/go-semver/semver"
)

// Version ... wraps semver and includes flag for leading `v`
type Version struct {
	leadingV bool
	semver   *semver.Version
}

// String ... returns semantic version string including leading `v`
// if necessary
func (v Version) String() string {
	if v.leadingV {
		return "v" + v.semver.String()
	}
	return v.semver.String()
}

// Bump ... bumps the version string according to scope
func (v Version) Bump(scope string) error {
	switch scope {
	case "patch":
		v.semver.BumpPatch()
	case "minor":
		v.semver.BumpMinor()
	case "major":
		v.semver.BumpMajor()
	default:
		return errors.New("scope must be one of: patch, minor, or major")
	}
	return nil
}
