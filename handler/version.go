package handler

import "errors"

// Set with release ldflags. This edition never installs upstream binaries.
var Version = "dev"

func (a *App) GetVersion() string { return Version }

// Preserve old bindings without contacting the upstream update service.
func (a *App) GetLatestVersion() (string, error) { return Version, nil }

func (a *App) UpdateTo(version string) error {
	return errors.New("automatic updates are disabled; download this edition from its GitHub Releases page")
}
