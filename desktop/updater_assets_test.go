package main

import (
	"testing"

	"github.com/wailsapp/wails/v3/pkg/updater"
	updaterGithub "github.com/wailsapp/wails/v3/pkg/updater/providers/github"
)

func TestDesktopUpdateAssetMatcherSelectsMacOSZip(t *testing.T) {
	assets := []updaterGithub.ReleaseAsset{
		{Name: "ccx-darwin-arm64"},
		{Name: "CCX-Desktop-1.2.3-darwin-arm64.dmg"},
		{Name: "CCX-Desktop-1.2.3-darwin-arm64.dmg.sha256"},
		{Name: "CCX-Desktop-1.2.3-darwin-arm64.zip"},
	}

	got := desktopUpdateAssetMatcher(updater.CheckRequest{Platform: "darwin", Arch: "arm64"}, assets)
	if got != 3 {
		t.Fatalf("desktopUpdateAssetMatcher() = %d, want 3", got)
	}
}

func TestDesktopUpdateAssetMatcherRejectsWindowsInstallers(t *testing.T) {
	assets := []updaterGithub.ReleaseAsset{
		{Name: "CCX-Desktop-1.2.3-windows-amd64-setup.exe"},
		{Name: "CCX-Desktop-1.2.3-windows-amd64-store.msix"},
		{Name: "ccx-windows-amd64.exe"},
	}

	got := desktopUpdateAssetMatcher(updater.CheckRequest{Platform: "windows", Arch: "amd64"}, assets)
	if got != -1 {
		t.Fatalf("desktopUpdateAssetMatcher() = %d, want -1", got)
	}
}
