package main

import (
	"runtime"
	"strings"

	"github.com/wailsapp/wails/v3/pkg/updater"
	updaterGithub "github.com/wailsapp/wails/v3/pkg/updater/providers/github"
)

func supportsInAppUpdate(distribution string) bool {
	if strings.EqualFold(distribution, "store") {
		return false
	}
	// Windows Store 由 Store 自动更新；NSIS/GitHub 包含 ccx-go.exe sidecar，
	// 不适合使用 Wails 当前的单目标替换流程。
	return runtime.GOOS == "darwin"
}

func newGitHubUpdaterProvider() (updater.Provider, error) {
	return updaterGithub.New(updaterGithub.Config{
		Repository:    "BenedictKing/ccx",
		ChecksumAsset: checksumAssetName(),
		AssetMatcher:  desktopUpdateAssetMatcher,
	})
}

func checksumAssetName() string {
	switch runtime.GOOS {
	case "darwin":
		return "checksums-macos.txt"
	case "windows":
		return "checksums-windows.txt"
	case "linux":
		return "checksums-linux.txt"
	default:
		return ""
	}
}

func desktopUpdateAssetMatcher(req updater.CheckRequest, assets []updaterGithub.ReleaseAsset) int {
	platform := strings.ToLower(req.Platform)
	arch := strings.ToLower(req.Arch)
	for i, asset := range assets {
		name := strings.ToLower(asset.Name)
		if !strings.HasPrefix(name, "ccx-desktop-") {
			continue
		}
		if platform != "" && !strings.Contains(name, platform) {
			continue
		}
		if arch != "" && !assetNameContainsArch(name, arch) {
			continue
		}
		if !isUpdaterPayload(name, platform) {
			continue
		}
		return i
	}
	return -1
}

func isUpdaterPayload(name, platform string) bool {
	return platform == "darwin" && strings.HasSuffix(name, ".zip")
}

func assetNameContainsArch(name, arch string) bool {
	if strings.Contains(name, arch) {
		return true
	}
	if arch == "amd64" {
		return strings.Contains(name, "x86_64") || strings.Contains(name, "x64")
	}
	if arch == "arm64" {
		return strings.Contains(name, "aarch64")
	}
	return false
}
