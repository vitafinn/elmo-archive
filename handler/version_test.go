package handler

import "testing"

func TestCommunityBuildDoesNotUseUpstreamUpdater(t *testing.T) {
	app := NewApp()
	version, err := app.GetLatestVersion()
	if err != nil || version != Version {
		t.Fatalf("unexpected version: %q, %v", version, err)
	}
	if app.UpdateTo("anything") == nil {
		t.Fatal("binary replacement must be disabled")
	}
}
