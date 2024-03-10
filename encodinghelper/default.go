package encodinghelper

import "runtime"

// PlatformDefaultEncoding returns the default encoding for the current platform.
func PlatformDefaultEncoding() string {
	if runtime.GOOS == `windows` {
		return `win1252`
	} else {
		return `utf8`
	}
}
