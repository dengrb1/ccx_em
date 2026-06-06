//go:build !windows

package appdirs

func windowsPackagedRoamingDir() string {
	return ""
}
