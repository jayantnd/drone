package gogs

import (
	"io"
	"net/http"

	"github.com/drone/drone/model"
)

const (
	hookEvent   = "X-Gogs-Event"
	hookPush    = "push"
	hookCreated = "create"

	refBranch = "branch"
	refTag    = "tag"
)

// parseHook parses a Bitbucket hook from an http.Request request and returns
// Repo and Build detail. If a hook type is unsupported nil values are returned.
func parseHook(r *http.Request) (*model.Repo, *model.Build, error) {
	switch r.Header.Get(hookEvent) {
	case hookPush:
		return parsePushHook(r.Body)
	case hookCreated:
		return parseCreatedHook(r.Body)
	}
	return nil, nil, nil
}

// parsePushHook parses a push hook and returns the Repo and Build details.
// If the commit type is unsupported nil values are returned.
func parsePushHook(payload io.Reader) (*model.Repo, *model.Build, error) {
	var (
		repo  *model.Repo
		build *model.Build
	)

	push, err := parsePush(payload)
	if err != nil {
		return nil, nil, err
	}

	// is this even needed?
	if push.RefType == refBranch {
		return nil, nil, nil
	}

	repo = repoFromPush(push)
	build = buildFromPush(push)
	return repo, build, err
}

// parseCreatedHook parses a push hook and returns the Repo and Build details.
// If the commit type is unsupported nil values are returned.
func parseCreatedHook(payload io.Reader) (*model.Repo, *model.Build, error) {
	var (
		repo  *model.Repo
		build *model.Build
	)

	push, err := parsePush(payload)
	if err != nil {
		return nil, nil, err
	}

	if push.RefType != refTag {
		return nil, nil, nil
	}

	repo = repoFromPush(push)
	build = buildFromTag(push)
	return repo, build, err
}
