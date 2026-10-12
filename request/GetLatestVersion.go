package request

import "errors"

// Upstream auto-updates are disabled for this community continuation.
func GetLatestVersion() (string, error) {
	return "", errors.New("check the community GitHub Releases page for updates")
}
